import React, { useState } from 'react';
import { submitPurchaseOrder } from '../api/purchase';
import { useTranslation } from 'react-i18next';
import { Send, FilePlus, DollarSign, Package } from 'lucide-react';
import './PurchaseRequest.css';

export const PurchaseRequest: React.FC = () => {
  const { t } = useTranslation();
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
      alert(t('Submit Failed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="purchase-container fade-in">
      <div className="glass-form-card">
        <div className="form-header">
          <div className="icon-wrapper"><FilePlus size={28} /></div>
          <h2>{t('New Purchase Request')}</h2>
          <p>{t('Submit description')}</p>
        </div>

        {success && (
          <div className="success-banner">
            {t('Success Banner')}
          </div>
        )}

        <form onSubmit={handleSubmit} className="purchase-form">
          <div className="form-group">
            <label>{t('PR Title')}</label>
            <input 
              required
              type="text" 
              placeholder={t('PR Title Placeholder')}
              value={formData.title}
              onChange={e => setFormData({...formData, title: e.target.value})}
            />
          </div>

          <div className="form-group">
            <label>{t('Total Amount')}</label>
            <div className="input-with-icon">
              <DollarSign size={18} className="input-icon" />
              <input 
                required
                type="number" 
                min="0"
                step="0.01"
                placeholder={t('Amount Placeholder')}
                value={formData.amount}
                onChange={e => setFormData({...formData, amount: e.target.value})}
              />
            </div>
          </div>

          {/* 这里的自定义字段将统一封装进入 JSONB 结构中，体现极致扩展性 */}
          <div className="custom-fields-section">
            <div className="section-title">{t('Dynamic Fields Section')}</div>
            
            <div className="form-group">
              <label>{t('Linked Material Group')}</label>
              <div className="input-with-icon">
                <Package size={18} className="input-icon" />
                <input 
                  type="text" 
                  placeholder={t('Material Placeholder')}
                  value={formData.customMaterial}
                  onChange={e => setFormData({...formData, customMaterial: e.target.value})}
                />
              </div>
            </div>

            <div className="form-group">
              <label>{t('Project Code')}</label>
              <input 
                type="text" 
                placeholder={t('Project Placeholder')}
                value={formData.customProject}
                onChange={e => setFormData({...formData, customProject: e.target.value})}
              />
            </div>
          </div>

          <button type="submit" className="submit-btn" disabled={loading}>
            <Send size={18} /> {loading ? t('Processing') : t('Submit PR')}
          </button>
        </form>
      </div>
    </div>
  );
};
