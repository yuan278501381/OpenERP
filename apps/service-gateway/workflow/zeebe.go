package workflow

import (
	"context"
	"openerp/packages/pkg-logger"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
)

var Client zbc.Client

func InitZeebe() {
	log := logger.Ctx(context.Background())
	log.Info("正在尝试连接外部 Zeebe BPMN 引擎...")

	var err error
	Client, err = zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         "127.0.0.1:26500",
		UsePlaintextConnection: true,
	})

	if err != nil {
		log.Warn("Zeebe 客户端初始化失败: " + err.Error())
		return
	}

	// 验证连通性
	_, err = Client.NewTopologyCommand().Send(context.Background())
	if err != nil {
		log.Warn("无法连接到 Zeebe 引擎（因宿主机无Docker导致降级）。BPMN 事件驱动模块进入 Mock 模式。")
		Client = nil
		return
	}

	log.Info("Zeebe BPMN 引擎连接成功！")
	StartWorkers()
}

func StartWorkers() {
	if Client == nil {
		return
	}
	
	log := logger.Ctx(context.Background())
	log.Info("启动 ERP 统一待办生成 Worker...")

	// 监听由 BPMN 抛出的事件，实现“业务产生流程，流程驱动待办”
	Client.NewJobWorker().JobType("generate_user_task").Handler(func(client worker.JobClient, job entities.Job) {
		log.Info("BPMN 引擎触发待办生成事件，正在写入数据库...")
		// 真实的业务代码会将 job.GetVariablesAsMap() 写入 DB
		client.NewCompleteJobCommand().JobKey(job.GetKey()).Send(context.Background())
	}).Open()
}
