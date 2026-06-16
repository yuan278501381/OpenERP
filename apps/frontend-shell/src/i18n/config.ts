import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

const resources = {
  en: {
    translation: {
      "Dashboard": "Dashboard",
      "Master Data": "Master Data",
      "Sales & Distribution": "Sales & Distribution",
      "Purchasing": "Purchasing",
      "Inventory Management": "Inventory Management",
      "Production (PP)": "Production (PP)",
      "Quality Management": "Quality Management",
      "Plant Maintenance": "Plant Maintenance",
      "Finance (FI/CO)": "Finance (FI/CO)",
      "Human Resources": "Human Resources",
      "System & Org": "System & Org",
      "Task Center": "Task Center",
      "Task Subtitle": "Tasks find people · Smart Collaboration · Global Visibility",
      "Syncing": "Syncing system tasks...",
      "Process Now": "Process Now",
      "Current Node": "Current Node: ",
      "Switch Language": "切换至中文",
      "New Purchase Request": "New Purchase Request",
      "Submit description": "Upon submission, ERP event-driven workflow will be triggered automatically.",
      "Success Banner": "✅ Submitted successfully! Workflow initiated. Check the Task Center for approval tasks.",
      "PR Title": "Request Title",
      "PR Title Placeholder": "e.g.: Q3 2026 Server Hardware Procurement",
      "Total Amount": "Estimated Total Amount (¥)",
      "Amount Placeholder": "0.00 (Over 100k triggers special approval)",
      "Dynamic Fields Section": "Dynamic ExtData JSONB Fields",
      "Linked Material Group": "Linked Material Group (MRP Reservation)",
      "Material Placeholder": "e.g.: PCB / High-freq Chipsets",
      "Project Code": "Project Code",
      "Project Placeholder": "e.g.: Project-A-Secret",
      "Submit PR": "Submit Purchase Request",
      "Processing": "Processing & Dispatching...",
      "Submit Failed": "Submit failed, check your network",
      "Material Master Data": "Material Master Data",
      "Material Subtitle": "Manage all raw materials, consumables, and finished goods.",
      "Filter": "Filter",
      "Export": "Export",
      "New Material": "New Material",
      "Loading material data...": "Loading material data...",
      "No materials found.": "No materials found.",
      "Material Code": "Material Code",
      "Name / Description": "Name / Description",
      "Category": "Category",
      "Stock": "Stock",
      "Unit": "Unit",
      "Standard Price": "Standard Price",
      "Status": "Status",
      "Last Updated": "Last Updated",
      "Search Placeholder": "Search ERP (Ctrl+K)",
      "Admin User": "Admin User",
      "System Admin": "System Admin",
      "Coming Soon": " Module Coming Soon",
      "Approval": "Approval",
      "Task": "Task",
      "Notice": "Notice",
      "Warning": "Warning",
      "Critical": "Critical",
      "Normal": "Normal",
      "采购单审批 - (数据库真实数据)": "Purchase Order Approval - (Live DB Data)",
      "系统权限开通申请 - (数据库真实数据)": "System Access Request - (Live DB Data)",
      "月度盘点确认 - (数据库真实数据)": "Monthly Inventory Confirmation - (Live DB Data)",
      "财务总监复核": "CFO Review",
      "IT 运维节点": "IT Ops Node",
      "仓储主管确认": "Warehouse Manager Confirmation"
    }
  },
  zh: {
    translation: {
      "Dashboard": "控制台",
      "Master Data": "主数据中心",
      "Sales & Distribution": "销售与分销",
      "Purchasing": "采购管理",
      "Inventory Management": "库存管理",
      "Production (PP)": "生产制造",
      "Quality Management": "质量管理",
      "Plant Maintenance": "设备维护",
      "Finance (FI/CO)": "财务核算",
      "Human Resources": "人力资源",
      "System & Org": "系统设置",
      "Task Center": "统一任务中枢",
      "Task Subtitle": "事找人 · 智能协同 · 全局可视",
      "Syncing": "正在同步全系统待办...",
      "Process Now": "立即处理",
      "Current Node": "当前卡点: ",
      "Switch Language": "Switch to English",
      "New Purchase Request": "新建采购申请单",
      "Submit description": "提交后将自动触发 ERP 事件驱动工作流，生成审批待办",
      "Success Banner": "✅ 提交成功！工作流已自动流转，请前往【统一待办】大厅查看生成的审批任务。",
      "PR Title": "申请单标题",
      "PR Title Placeholder": "例如：2026年Q3服务器硬件采购",
      "Total Amount": "预估总金额 (¥)",
      "Amount Placeholder": "0.00 (超过10万将触发特批流程)",
      "Dynamic Fields Section": "动态扩展表单字段 (ExtData JSONB)",
      "Linked Material Group": "关联物料族 (MRP预留机制)",
      "Material Placeholder": "例如：PCB主板 / 高频芯片组",
      "Project Code": "归属项目组代码",
      "Project Placeholder": "例如：Project-A-Secret",
      "Submit PR": "提交采购单 (触发流程)",
      "Processing": "业务处理与事件分发中...",
      "Submit Failed": "提交失败，请检查网络",
      "Material Master Data": "物料主数据",
      "Material Subtitle": "管理所有原材料、消耗品与产成品。",
      "Filter": "筛选",
      "Export": "导出",
      "New Material": "新增物料",
      "Loading material data...": "正在加载物料数据...",
      "No materials found.": "未找到物料数据。",
      "Material Code": "物料编码",
      "Name / Description": "物料名称 / 描述",
      "Category": "分类",
      "Stock": "库存",
      "Unit": "单位",
      "Standard Price": "标准价",
      "Status": "状态",
      "Last Updated": "最后更新",
      "Search Placeholder": "全局搜索 ERP (Ctrl+K)",
      "Admin User": "管理员用户",
      "System Admin": "系统管理员",
      "Coming Soon": " 模块正在建设中",
      "Approval": "审批单",
      "Task": "业务任务",
      "Notice": "系统通知",
      "Warning": "警告 (将逾期)",
      "Critical": "严重 (已逾期)",
      "Normal": "正常流转",
      "采购单审批 - (数据库真实数据)": "采购单审批 - (数据库真实数据)",
      "系统权限开通申请 - (数据库真实数据)": "系统权限开通申请 - (数据库真实数据)",
      "月度盘点确认 - (数据库真实数据)": "月度盘点确认 - (数据库真实数据)",
      "财务总监复核": "财务总监复核",
      "IT 运维节点": "IT 运维节点",
      "仓储主管确认": "仓储主管确认"
    }
  }
} as const;

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: 'translation';
    resources: typeof resources['en'];
  }
}

i18n
  .use(initReactI18next)
  .init({
    resources,
    lng: "zh",
    fallbackLng: "en",
    interpolation: {
      escapeValue: false
    }
  });

export default i18n;
