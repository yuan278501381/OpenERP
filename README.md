# 🚀 OpenERP - 下一代世界级开源企业核心中枢

![Version](https://img.shields.io/badge/version-1.0.0--alpha-blue.svg)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)
![Architecture](https://img.shields.io/badge/Architecture-World--Class-success)

> **Vision (愿景)**  
> 打造一款能与 SAP S/4HANA、SAP B1 正面交锋的世界级、极高扩展性、现代化的开源 ERP 系统。
> 我们坚信，中小企业同样值得拥有极致性能、极致 UI 与高度自由扩展的“业财一体化”信息化底座。通过融入最前沿的 UDF (动态字段) 与 UDT (动态表单) 技术，OpenERP 让千万级复杂度的业务流转变得轻盈如燕。

---

## 🌟 核心设计理念

1. **绝对的业财一体化**：每一笔供应链与生产的实物流转，都将自动且精准地转化为财务总账凭证，彻底消灭业务与财务的数据孤岛。
2. **极简的高密度交互**：摒弃传统 ERP 繁琐臃肿的操作界面，采用苹果风与专业级金融数据看板融合的浅色高密度 UI，让业务员聚焦于“数据本身”。
3. **“无边界”扩展架构 (UDF / UDT / UDO)**：
   - **JSONB 无极扩表 (UDF)**：底层拥抱 PostgreSQL/MySQL JSONB 技术，所有核心单据均内置 `ExtData`，客户新增自定义字段（User Defined Fields）再也无需修改数据库 Schema 和 API 源码。
   - **元数据驱动渲染 (UDT)**：前端表单摒弃硬编码，依靠 `SysFormField` 动态表单接口获取元数据渲染界面，实现自定义表（User Defined Tables）的无限热插拔。
   - **用户定义对象 (UDO)**：将自定义的主表与明细表打包，自动生成标准 CRUD 接口、权限校验与业务对象级缓存。无需编写一行代码，即可让业务人员“像搭积木一样”捏出全新的业务单据（如：车辆管理单、宿舍分配单），并使其具备原生单据的完全体体验。

---

## ⚡ 当前已上线功能 (Phase 1 & Wave 1)

目前，我们已经铺设了最核心的基础设施并打通了首条核心业务流：

### 🛠️ 坚若磐石的技术底座
- **Go 语言核心网关**：基于洋葱架构（Clean Architecture），高内聚低耦合。
- **BPMN 流程引擎接入**：内置 Zeebe 客户端，支持复杂审批流与“事找人”的任务驱动机制。
- **40+ 核心业务模型**：涵盖了多组织架构、财务期间、多币种、物料组等核心骨架，并向下兼容至复杂的制造成本体系。
- **零死角的国际化引擎 (i18n)**：前端不仅实现了中英双语无缝秒切，还集成了 `eslint-plugin-i18next` 的 AST 强校验与 TypeScript 强类型泛型锁，在编译期彻底杜绝漏译与硬编码，符合顶级跨国架构标准。

### 📦 供应链中枢 (MM/SD)
- **销售接单与 ATP 校验**：销售单工作台（Sales Orders）具备实时的 ATP (Available To Promise) 模拟能力，接单即秒级校验底层可用库存。
- **动态库存与加权移动平均**：物料移动大屏（Goods Movement）支持收发货。在核心层，入库动作已实现了基于数据库事务与悲观锁的“移动加权平均成本 (Moving Average Cost)”实时重算。
- **统一任务大厅**：融合了审批单、系统通知、业务预警的卡片式任务中枢。

---

## 🚀 未来开发路线图 (Roadmap)

我们的征途才刚刚开始，以下模块正在我们的作战沙盘中稳步推进：

### 🌊 Wave 2: 制造与品控引擎 (PP/QM/PM)
- **MRP 算料爆炸引擎**：根据销售订单、BOM（物料清单）与现有库存，通过有向无环图自动推演生成生产工单与采购申请。
- **车间排产大屏**：实时展示各工作中心 (Work Center) 与机器设备的产能负载情况。
- **全链路质量追溯**：从采购收货到生产工单完工，嵌入检验批次 (Inspection Lot) 与质检结果判定机制。

### 🌊 Wave 3: 深度业财一体与人力 (FI/CO & HCM)
- **事件驱动的自动过账**：当仓库完成收发货时，引擎将根据 `SysAccountDetermination` (科目决定表) 自动生成带借贷方向的日记账凭证 (Journal Entry)，并实现月末费用的自动分摊。
- **动态预算控制 (Budget Control)**：实时的部门级与项目级费控卡点，超过预算的报销与采购将自动触发 BPMN 特批流。
- **员工与组织树**：构建带有时间切片的复杂汇报线与组织架构体系。

### 🌊 Wave 4: 数字化与智能化看板 (BI & AI)
- 引入前端极速图表库，直接挖掘底层的 JSONB 数据源。
- 构建涵盖应收账款账龄 (AR Aging)、库存周转率、生产准交率的高层驾驶舱。

---

## ⚙️ 如何在本地启动体验

### 后端服务 (Service Gateway)
\`\`\`bash
cd apps/service-gateway
go mod tidy
swag init --parseDependency --parseInternal
go run main.go
\`\`\`
*默认监听在 `http://localhost:8080`，启动时会自动同步数据库并初始化演示测试数据。*

### 前端看板 (Frontend Shell)
\`\`\`bash
cd apps/frontend-shell
npm install
npm run dev
\`\`\`
*默认监听在 `http://localhost:3000`，极致流畅的交互等您检阅。*

---

> *"我们不是在写代码，我们是在用代码重构一家企业的数字生命体系。" — OpenERP 架构组*
