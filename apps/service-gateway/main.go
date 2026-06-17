package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "openerp/apps/service-gateway/api/v1"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/apps/service-gateway/workflow"
	"openerp/packages/pkg-logger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "openerp/apps/service-gateway/docs" // Swagger docs
)

func traceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = "trace-uuid-placeholder"
		}

		ctx := context.WithValue(c.Request.Context(), logger.TraceIDKey, traceID)
		ctx = context.WithValue(ctx, logger.TenantIDKey, c.GetHeader("X-Tenant-ID"))
		ctx = context.WithValue(ctx, logger.UserIDKey, c.GetHeader("X-User-ID"))

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func seedFinanceDefaults(log *zap.Logger) {
	accounts := []models.SysGLAccount{
		{AccountCode: "1405", AccountName: "Inventory", AccountType: "Asset"},
		{AccountCode: "2202", AccountName: "GR/IR Clearing", AccountType: "Liability"},
		{AccountCode: "6401", AccountName: "Cost of Goods Sold", AccountType: "Expense"},
	}
	for _, account := range accounts {
		if err := db.DB.Where("account_code = ?", account.AccountCode).FirstOrCreate(&account).Error; err != nil {
			log.Error("默认总账科目初始化失败", zap.String("accountCode", account.AccountCode), zap.Error(err))
		}
	}

	rules := []models.SysAccountDetermination{
		{
			TransactionType: v1.PostingEventGoodsReceipt,
			ItemCategory:    v1.DefaultItemCategory,
			DebitAccount:    "1405",
			CreditAccount:   "2202",
		},
		{
			TransactionType: v1.PostingEventGoodsIssue,
			ItemCategory:    v1.DefaultItemCategory,
			DebitAccount:    "6401",
			CreditAccount:   "1405",
		},
	}
	for _, rule := range rules {
		err := db.DB.Where(
			"transaction_type = ? AND item_category = ?",
			rule.TransactionType,
			rule.ItemCategory,
		).FirstOrCreate(&rule).Error
		if err != nil {
			log.Error("默认科目决定规则初始化失败", zap.String("transactionType", rule.TransactionType), zap.String("itemCategory", rule.ItemCategory), zap.Error(err))
		}
	}
}

func seedMaterialDefaults(log *zap.Logger) {
	materials := []models.SysMaterial{
		{
			MaterialCode:      "MAT-001",
			MaterialName:      "Aluminum Sheet 2mm",
			Category:          "RAW",
			Unit:              "kg",
			IsActive:          true,
			Stock:             1250,
			MovingAverageCost: 4.5,
		},
		{
			MaterialCode:      "MAT-005",
			MaterialName:      "Packaging Box Standard",
			Category:          "PACK",
			Unit:              "pcs",
			IsActive:          true,
			Stock:             12000,
			MovingAverageCost: 0.5,
		},
	}
	for _, material := range materials {
		if err := db.DB.Where("material_code = ?", material.MaterialCode).FirstOrCreate(&material).Error; err != nil {
			log.Error("默认物料初始化失败", zap.String("materialCode", material.MaterialCode), zap.Error(err))
		}
	}
}

func serverAddress() string {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}

// @title OpenERP 企业级信息化中枢 API
// @version 1.0
// @description 这是一个高度可扩展的开源 ERP 底座，支持动态字段(JSONB)、复杂审批流以及物料需求计划(MRP)。
// @host localhost:8080
// @BasePath /
func main() {
	logger.Init("service-gateway", "info")
	bgLog := logger.Ctx(context.Background())
	bgLog.Info("统一任务中枢 (API Gateway) 正在启动...")

	// 1. 初始化数据库并自动迁移表结构
	db.InitDB()
	db.DB.AutoMigrate(&models.SysTask{})
	db.DB.AutoMigrate(&models.SysPurchaseOrder{})
	db.DB.AutoMigrate(&models.SysPurchaseReceipt{})
	db.DB.AutoMigrate(&models.SysPurchaseReceiptLine{})
	db.DB.AutoMigrate(&models.SysMaterial{})
	db.DB.AutoMigrate(&models.SysBOM{})
	db.DB.AutoMigrate(&models.SysFormField{})

	// Stage 1 Models
	db.DB.AutoMigrate(
		&models.SysOrganization{},
		&models.SysFinancialPeriod{},
		&models.SysPaymentTerm{},
		&models.SysUoMGroup{},
		&models.SysAccountDetermination{},
		&models.SysNotificationTemplate{},
	)

	// Stage 2 Models
	db.DB.AutoMigrate(
		&models.SysGLAccount{},
		&models.SysJournalEntry{},
		&models.SysJournalEntryLine{},
		&models.SysCostCenter{},
		&models.SysBudget{},
		&models.SysARInvoice{},
		&models.SysAPInvoice{},
		&models.SysPayment{},
		&models.SysReconciliation{},
		&models.SysExpenseClaim{},
		&models.SysFixedAsset{},
		&models.SysBillOfExchange{},
		&models.SysExchangeRate{},
		&models.SysTaxGroup{},
		&models.SysTaxInvoice{},
	)

	// Stage 3 Models
	db.DB.AutoMigrate(
		&models.SysBatch{},
		&models.SysSerialNumber{},
		&models.SysBusinessPartner{},
		&models.SysPricingCondition{},
		&models.SysSalesQuotation{},
		&models.SysSalesQuotationLine{},
		&models.SysSalesOrder{},
		&models.SysSalesOrderLine{},
		&models.SysDelivery{},
		&models.SysDeliveryLine{},
		&models.SysSalesReturnRequest{},
		&models.SysSalesReturnRequestLine{},
		&models.SysSalesReturn{},
		&models.SysSalesReturnLine{},
		&models.SysPurchaseRequisition{},
		&models.SysPurchaseRequisitionLine{},
		&models.SysPurchaseOrderLine{},
		&models.SysLandedCost{},
		&models.SysPurchaseReturnRequest{},
		&models.SysPurchaseReturnRequestLine{},
		&models.SysPurchaseReturn{},
		&models.SysPurchaseReturnLine{},
		&models.SysGoodsMovement{},
		&models.SysInventoryTransferRequest{},
		&models.SysInventoryTransfer{},
		&models.SysInventoryCounting{},
		&models.SysWarehouseBin{},
		&models.SysInventoryValuation{},
	)

	// Stage 4 Models
	db.DB.AutoMigrate(
		&models.SysWorkCenter{},
		&models.SysResource{},
		&models.SysRouting{},
		&models.SysProductionOrder{},
		&models.SysInspectionLot{},
		&models.SysEquipment{},
		&models.SysEmployee{},
	)
	seedFinanceDefaults(bgLog)
	seedMaterialDefaults(bgLog)

	// 2. 初始化 BPMN 引擎客户端 (Zeebe)
	workflow.InitZeebe()

	// 2. 如果数据库为空，注入测试数据（模拟由工作流引擎产生的待办）
	var count int64
	db.DB.Model(&models.SysTask{}).Count(&count)
	if count == 0 {
		bgLog.Info("数据库为空，自动注入初始测试待办任务...")
		db.DB.Create(&models.SysTask{
			TaskID:    "T-DB-1001",
			Title:     "采购单审批 - (数据库真实数据)",
			Type:      "Approval",
			Node:      "财务总监复核",
			SlaStatus: "Warning",
			TenantID:  "TENANT-001",
			UserID:    "USER-1024",
			ExtraData: "{}",
		})
		db.DB.Create(&models.SysTask{
			TaskID:    "T-DB-1002",
			Title:     "系统权限开通申请 - (数据库真实数据)",
			Type:      "Task",
			Node:      "IT 运维节点",
			SlaStatus: "Critical",
			TenantID:  "TENANT-001",
			UserID:    "USER-1024",
			ExtraData: "{}",
		})
		db.DB.Create(&models.SysTask{
			TaskID:    "T-DB-1003",
			Title:     "月度盘点确认 - (数据库真实数据)",
			Type:      "Notice",
			Node:      "仓储主管确认",
			SlaStatus: "Normal",
			TenantID:  "TENANT-001",
			UserID:    "USER-1024",
			ExtraData: "{}",
		})
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(traceMiddleware())

	// Swagger 在线文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 从真实的数据库读取待办列表
	r.GET("/openerp/v1/tasks", func(c *gin.Context) {
		log := logger.Ctx(c.Request.Context())

		var tasks []models.SysTask

		// 获取租户上下文，实现严格的数据隔离
		tenantID, _ := c.Request.Context().Value(logger.TenantIDKey).(string)
		if tenantID == "" {
			tenantID = "TENANT-001"
		}

		log.Info("正在从数据库拉取待办", zap.String("tenant", tenantID))

		// GORM 查询
		result := db.DB.Where("tenant_id = ?", tenantID).Find(&tasks)
		if result.Error != nil {
			log.Error("查询失败", zap.Error(result.Error))
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "数据库查询失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": tasks,
		})
	})

	// Purchasing API
	r.POST("/openerp/v1/purchase-orders", v1.CreatePurchaseOrder)
	r.GET("/openerp/v1/purchase-orders", v1.GetPurchaseOrders)
	r.POST("/openerp/v1/purchase-receipts", v1.CreatePurchaseReceipt)
	r.GET("/openerp/v1/purchase-receipts", v1.GetPurchaseReceipts)

	// 动态表单元数据接口
	r.GET("/openerp/v1/forms/:type/meta", v1.GetFormMeta)
	r.POST("/openerp/v1/forms/:type/meta", v1.AddFormMeta)

	// 物料主数据接口
	r.GET("/openerp/v1/materials", v1.GetMaterials)

	// Organizations API
	r.POST("/openerp/v1/organizations", v1.CreateOrganization)
	r.GET("/openerp/v1/organizations", v1.GetOrganizations)
	r.PUT("/openerp/v1/organizations/:id", v1.UpdateOrganization)
	r.DELETE("/openerp/v1/organizations/:id", v1.DeleteOrganization)

	// GLAccounts API
	r.POST("/openerp/v1/gl-accounts", v1.CreateGLAccount)
	r.GET("/openerp/v1/gl-accounts", v1.GetGLAccounts)
	r.PUT("/openerp/v1/gl-accounts/:id", v1.UpdateGLAccount)
	r.DELETE("/openerp/v1/gl-accounts/:id", v1.DeleteGLAccount)

	// ExpenseClaims API
	r.POST("/openerp/v1/expense-claims", v1.CreateExpenseClaim)
	r.GET("/openerp/v1/expense-claims", v1.GetExpenseClaims)
	r.PUT("/openerp/v1/expense-claims/:id", v1.UpdateExpenseClaim)
	r.DELETE("/openerp/v1/expense-claims/:id", v1.DeleteExpenseClaim)

	// Business Partners API
	r.POST("/openerp/v1/business-partners", v1.CreateBusinessPartner)
	r.GET("/openerp/v1/business-partners", v1.GetBusinessPartners)
	r.PUT("/openerp/v1/business-partners/:id", v1.UpdateBusinessPartner)
	r.DELETE("/openerp/v1/business-partners/:id", v1.DeleteBusinessPartner)

	// Sales Orders API
	r.POST("/openerp/v1/sales-orders", v1.CreateSalesOrder)
	r.GET("/openerp/v1/sales-orders", v1.GetSalesOrders)
	r.PUT("/openerp/v1/sales-orders/:id", v1.UpdateSalesOrder)
	r.DELETE("/openerp/v1/sales-orders/:id", v1.DeleteSalesOrder)

	// Inventory & Goods Movement API
	r.POST("/openerp/v1/goods-movement", v1.CreateGoodsMovement)
	r.GET("/openerp/v1/goods-movement", v1.GetGoodsMovements)

	// Production Orders API
	r.POST("/openerp/v1/production-orders", v1.CreateProductionOrder)
	r.GET("/openerp/v1/production-orders", v1.GetProductionOrders)
	r.PUT("/openerp/v1/production-orders/:id", v1.UpdateProductionOrder)
	r.DELETE("/openerp/v1/production-orders/:id", v1.DeleteProductionOrder)
	r.POST("/openerp/v1/production-orders/:id/report-completion", v1.ReportCompletion)

	// MRP API
	r.POST("/openerp/v1/mrp/run", v1.RunMRP)

	// Finance API (自动过账引擎 & 日记账)
	r.POST("/openerp/v1/finance/auto-post", v1.AutoPostJournalEntry)
	r.GET("/openerp/v1/finance/journal-entries", v1.GetJournalEntries)

	addr := serverAddress()
	bgLog.Info("API Gateway 成功启动", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		bgLog.Error("服务异常退出", zap.Error(err))
	}
}
