package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/apps/service-gateway/workflow"
	v1 "openerp/apps/service-gateway/api/v1"
	"openerp/packages/pkg-logger"
	"gorm.io/datatypes"

	_ "openerp/apps/service-gateway/docs" // Swagger docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// 提交采购单接口 (支持外部 API 接入)
	r.POST("/openerp/v1/purchase-orders", func(c *gin.Context) {
		log := logger.Ctx(c.Request.Context())
		
		var req struct {
			OrderNo string                 `json:"orderNo"`
			Title   string                 `json:"title"`
			Amount  float64                `json:"amount"`
			ExtData map[string]interface{} `json:"extData"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数解析失败"})
			return
		}

		log.Info("收到采购单提交请求", zap.String("orderNo", req.OrderNo))

		po := models.SysPurchaseOrder{
			OrderNo:     req.OrderNo,
			Title:       req.Title,
			Amount:      req.Amount,
			ApplicantID: "USER-1024",
			DeptID:      "DEPT-FINANCE",
		}

		// 处理无限层级的自定义扩展字段 (序列化为 JSON)
		if req.ExtData != nil {
			extBytes, _ := json.Marshal(req.ExtData)
			po.ExtData = datatypes.JSON(extBytes)
		}

		// 1. 业务落库
		if err := db.DB.Create(&po).Error; err != nil {
			log.Error("采购单持久化失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "采购单保存失败"})
			return
		}

		// 2. 发起事件流转：交由 Zeebe 引擎进行调度
		if err := workflow.TriggerPurchaseApprovalWorkflow(db.DB, po); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "工作流触发失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200, 
			"msg":  "采购申请提交成功，已推入流程引擎",
			"data": po.ID,
		})
	})

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

	bgLog.Info("API Gateway 成功启动，监听在 :8080 端口")
	if err := r.Run(":8080"); err != nil {
		bgLog.Error("服务异常退出", zap.Error(err))
	}
}
