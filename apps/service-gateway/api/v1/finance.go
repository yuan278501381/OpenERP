package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"
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

	var newEntryID uint

	// 使用事务确保借贷及日记账头同时成功或失败
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查找科目决定规则 (Account Determination)
		var rule models.SysAccountDetermination
		if err := tx.Where("transaction_type = ? AND item_category = ?", req.EventType, req.ItemCategory).First(&rule).Error; err != nil {
			log.Error("未找到对应的科目决定规则", zap.Error(err))
			return fmt.Errorf("科目决定规则未配置: %w", err)
		}

		// 2. 构建日记账头
		entryNo := fmt.Sprintf("JE-%d", time.Now().UnixNano())
		je := models.SysJournalEntry{
			EntryNo:     entryNo,
			PostingDate: time.Now(),
			DocDate:     time.Now(),
			TotalAmount: req.Amount,
		}

		if err := tx.Create(&je).Error; err != nil {
			log.Error("创建日记账头失败", zap.Error(err))
			return err
		}
		newEntryID = je.ID

		// 3. 构建日记账行 (借贷必相等，复式记账原则)
		debitLine := models.SysJournalEntryLine{
			EntryID:      je.ID,
			AccountCode:  rule.DebitAccount,
			PayerBpID:    req.PayerBpID,
			BillToBpID:   req.BillToBpID,
			DebitAmount:  req.Amount,
			CreditAmount: 0,
			CostCenter:   req.CostCenter,
		}

		creditLine := models.SysJournalEntryLine{
			EntryID:      je.ID,
			AccountCode:  rule.CreditAccount,
			PayerBpID:    req.PayerBpID,
			BillToBpID:   req.BillToBpID,
			DebitAmount:  0,
			CreditAmount: req.Amount,
			CostCenter:   req.CostCenter,
		}

		if err := tx.Create(&debitLine).Error; err != nil {
			log.Error("创建日记账借方行失败", zap.Error(err))
			return err
		}
		if err := tx.Create(&creditLine).Error; err != nil {
			log.Error("创建日记账贷方行失败", zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "自动过账失败", "error": err.Error()})
		return
	}

	log.Info("自动过账完成", zap.Uint("entryID", newEntryID))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "过账成功", "data": newEntryID})
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
