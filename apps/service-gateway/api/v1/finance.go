package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"
)

const (
	PostingEventGoodsReceipt = "GoodsReceipt"
	PostingEventGoodsIssue   = "GoodsIssue"
	DefaultItemCategory      = "DEFAULT"
)

// AutoPostRequest 定义自动过账的请求参数
type AutoPostRequest struct {
	EventType    string  `json:"eventType" binding:"required"`
	ItemCategory string  `json:"itemCategory" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	CostCenter   string  `json:"costCenter"`
	Reference    string  `json:"reference"`
	PayerBpID    string  `json:"payerBpId"`
	BillToBpID   string  `json:"billToBpId"`
}

// @Summary 核心自动过账引擎 (Core Posting Engine)
// @Description 模拟接收事件(如收货)并自动生成复式记账(借贷必相等)
// @Tags Finance
// @Accept json
// @Produce json
// @Param request body AutoPostRequest true "过账请求参数"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/finance/auto-post [post]
func AutoPostJournalEntry(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var req AutoPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("自动过账参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	log.Info("开始执行自动过账", zap.String("eventType", req.EventType), zap.String("itemCategory", req.ItemCategory))

	var entry *models.SysJournalEntry

	// 使用事务确保借贷及日记账头同时成功或失败
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		entry, err = postJournalEntry(tx, req)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "自动过账失败", "error": err.Error()})
		return
	}

	log.Info("自动过账完成", zap.Uint("entryID", entry.ID))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "过账成功", "data": entry.ID})
}

func postJournalEntry(tx *gorm.DB, req AutoPostRequest) (*models.SysJournalEntry, error) {
	req.ItemCategory = normalizeItemCategory(req.ItemCategory)
	if strings.TrimSpace(req.EventType) == "" {
		return nil, fmt.Errorf("交易类型不能为空")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("过账金额必须大于 0")
	}

	rule, err := findAccountDetermination(tx, req.EventType, req.ItemCategory)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	entryNo := fmt.Sprintf("JE-%d", now.UnixNano())
	entry := models.SysJournalEntry{
		EntryNo:     entryNo,
		PostingDate: now,
		DocDate:     now,
		TotalAmount: req.Amount,
		ExtData:     buildJournalEntryExtData(req),
	}

	if err := tx.Create(&entry).Error; err != nil {
		return nil, fmt.Errorf("创建日记账头失败: %w", err)
	}

	lines := []models.SysJournalEntryLine{
		{
			EntryID:      entry.ID,
			AccountCode:  rule.DebitAccount,
			PayerBpID:    req.PayerBpID,
			BillToBpID:   req.BillToBpID,
			DebitAmount:  req.Amount,
			CreditAmount: 0,
			CostCenter:   req.CostCenter,
		},
		{
			EntryID:      entry.ID,
			AccountCode:  rule.CreditAccount,
			PayerBpID:    req.PayerBpID,
			BillToBpID:   req.BillToBpID,
			DebitAmount:  0,
			CreditAmount: req.Amount,
			CostCenter:   req.CostCenter,
		},
	}

	if err := tx.Create(&lines).Error; err != nil {
		return nil, fmt.Errorf("创建日记账行失败: %w", err)
	}
	entry.Lines = lines
	return &entry, nil
}

func findAccountDetermination(tx *gorm.DB, eventType, itemCategory string) (models.SysAccountDetermination, error) {
	var rule models.SysAccountDetermination
	err := tx.Where("transaction_type = ? AND item_category = ?", eventType, itemCategory).First(&rule).Error
	if err == nil {
		return rule, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) || itemCategory == DefaultItemCategory {
		return rule, fmt.Errorf("科目决定规则未配置: %w", err)
	}

	err = tx.Where("transaction_type = ? AND item_category = ?", eventType, DefaultItemCategory).First(&rule).Error
	if err != nil {
		return rule, fmt.Errorf("科目决定规则未配置: %w", err)
	}
	return rule, nil
}

func normalizeItemCategory(category string) string {
	category = strings.TrimSpace(category)
	if category == "" {
		return DefaultItemCategory
	}
	return category
}

func buildJournalEntryExtData(req AutoPostRequest) datatypes.JSON {
	payload := map[string]string{
		"eventType":    req.EventType,
		"itemCategory": req.ItemCategory,
		"reference":    req.Reference,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return datatypes.JSON(data)
}

// @Summary 获取日记账分录
// @Description 获取自动或手动创建的日记账分录及其行项目
// @Tags Finance
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/finance/journal-entries [get]
func GetJournalEntries(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var entries []models.SysJournalEntry
	// 预加载 Lines 关联数据，以便一起返回
	if err := db.DB.Preload("Lines").Find(&entries).Error; err != nil {
		log.Error("查询日记账分录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": entries})
}
