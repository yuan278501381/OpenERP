import React, { useState } from 'react';
import { submitPurchaseOrder } from '../api/purchase';
import { Send, FilePlus, DollarSign, Package } from 'lucide-react';
import './PurchaseRequest.css';

export const PurchaseRequest: React.FC = () => {
  const [formData, setFormData] = useState({
    title: '',
    amount: '',
    customMaterial: '',
    customProject: ''
  });
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await submitPurchaseOrder({
        orderNo: `PO-${Date.now()}`, // 实际开发应由后端生成单号
        title: formData.title,
        amount: parseFloat(formData.amount),
        extData: {
          material: formData.customMaterial,
          project: formData.customProject,
        }
      });
      setSuccess(true);
      setTimeout(() => setSuccess(false), 4000);
      setFormData({ title: '', amount: '', customMaterial: '', customProject: '' });
    } catch (error) {
      alert('提交失败，请检查网络');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="purchase-container fade-in">
      <div className="glass-form-card">
        <div className="form-header">
          <div className="icon-wrapper"><FilePlus size={28} /></div>
          <h2>新建采购申请单</h2>
          <p>提交后将自动触发 ERP 事件驱动工作流，生成审批待办</p>
        </div>

        {success && (
          <div className="success-banner">
            ✅ 提交成功！工作流已自动流转，请前往【统一待办】大厅查看生成的审批任务。
          </div>
        )}

        <form onSubmit={handleSubmit} className="purchase-form">
          <div className="form-group">
            <label>申请单标题</label>
            <input 
              required
              type="text" 
              placeholder="例如：2026年Q3服务器硬件采购"
              value={formData.title}
              onChange={e => setFormData({...formData, title: e.target.value})}
            />
          </div>

          <div className="form-group">
            <label>预估总金额 (¥)</label>
            <div className="input-with-icon">
              <DollarSign size={18} className="input-icon" />
              <input 
                required
                type="number" 
                min="0"
                step="0.01"
                placeholder="0.00 (超过10万将触发特批流程)"
                value={formData.amount}
                onChange={e => setFormData({...formData, amount: e.target.value})}
              />
            </div>
          </div>

          {/* 这里的自定义字段将统一封装进入 JSONB 结构中，体现极致扩展性 */}
          <div className="custom-fields-section">
            <div className="section-title">动态扩展表单字段 (ExtData JSONB)</div>
            
            <div className="form-group">
              <label>关联物料族 (MRP预留机制)</label>
              <div className="input-with-icon">
                <Package size={18} className="input-icon" />
                <input 
                  type="text" 
                  placeholder="例如：PCB主板 / 高频芯片组"
                  value={formData.customMaterial}
                  onChange={e => setFormData({...formData, customMaterial: e.target.value})}
                />
              </div>
            </div>

            <div className="form-group">
              <label>归属项目组代码</label>
              <input 
                type="text" 
                placeholder="例如：Project-A-Secret"
                value={formData.customProject}
                onChange={e => setFormData({...formData, customProject: e.target.value})}
              />
            </div>
          </div>

          <button type="submit" className="submit-btn" disabled={loading}>
            <Send size={18} /> {loading ? '业务处理与事件分发中...' : '提交采购单 (触发流程)'}
          </button>
        </form>
      </div>
    </div>
  );
};
