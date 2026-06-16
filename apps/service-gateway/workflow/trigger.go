package workflow

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"
)

// TriggerPurchaseApprovalWorkflow 触发采购审批工作流
// 这正是“业务产生流程，流程驱动待办”的核心枢纽。
func TriggerPurchaseApprovalWorkflow(db *gorm.DB, po models.SysPurchaseOrder) error {
	log := logger.Ctx(context.Background())

	if Client != nil {
		// 环境探测：如果 Docker Zeebe 真实存在，则向引擎发射 BPMN 实例启动请求
		request, err := Client.NewCreateInstanceCommand().
			BPMNProcessId("ERP_Purchase_Approval").
			LatestVersion().
			VariablesFromMap(map[string]interface{}{
				"orderNo": po.OrderNo,
				"amount":  po.Amount,
			})
		if err == nil {
			_, err = request.Send(context.Background())
			if err == nil {
				log.Info("成功触发真实 Zeebe 采购审批工作流！流程引擎已接管业务流转。")
				return nil
			}
		}
		log.Warn("真实 Zeebe 触发失败，准备降级为 Mock 事件处理: " + err.Error())
	}

	// ==========================================
	// 优雅降级 (Mock) 模式
	// 无论底层引擎是否就绪，业务层表现绝不能挂断！
	// 在无 Docker 状态下，网关自动接管工作流的 “代偿动作”：直接派发任务到统一待办大厅。
	// ==========================================
	log.Info("触发事件总线：采购单已创建，正在模拟 BPMN 计算并指派审批任务...")

	mockTask := models.SysTask{
		TaskID:    fmt.Sprintf("T-PO-%d", po.ID), // 关联采购单ID
		Title:     fmt.Sprintf("[急件审批] 采购单 %s", po.OrderNo),
		Type:      "Approval",
		Node:      "财务总监复核",
		SlaStatus: "Normal",
		TenantID:  "TENANT-001",
		UserID:    "USER-1024", // 模拟引擎将任务指派给财务总监的员工账号
	}

	// 模拟流控引擎的网关条件判断机制 (Amount > 100000 走特批)
	if po.Amount > 100000 {
		mockTask.SlaStatus = "Critical"
		mockTask.Node = "CEO 资金特批"
	}

	// 写入统一待办任务表 (模拟任务已抵达用户桌面)
	if err := db.Create(&mockTask).Error; err != nil {
		log.Error("派发统一待办失败: " + err.Error())
		return err
	}

	log.Info("🎉 统一待办派发成功！任务已推送到负责人的 TaskCenter。")
	return nil
}
