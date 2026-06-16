import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { ShoppingCart, Plus, Search, Eye, Edit2 } from 'lucide-react';

export const SalesOrderList: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  
  // Mock Data
  const [orders] = useState([
    { id: 'SO-2026-001', customer: 'Tech Corp', amount: 150000, status: 'Completed', date: '2026-06-15' },
    { id: 'SO-2026-002', customer: 'Global Industries', amount: 85000, status: 'Pending', date: '2026-06-16' },
    { id: 'SO-2026-003', customer: 'Local Retailer', amount: 12000, status: 'Processing', date: '2026-06-16' },
  ]);

  return (
    <div className="p-4" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--spacing-lg)' }}>
      <div className="panel p-4" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <h1 style={{ fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
            <ShoppingCart size={20} style={{ color: 'var(--brand-accent)' }} />
            {t('Sales Orders')}
          </h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: 'var(--font-sm)', marginTop: '4px' }}>
            {t('Manage customer sales orders')}
          </p>
        </div>
        <div style={{ display: 'flex', gap: '8px' }}>
          <div style={{ position: 'relative' }}>
            <Search size={16} style={{ position: 'absolute', left: '8px', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-secondary)' }} />
            <input 
              type="text" 
              placeholder={t('Search Placeholder')} 
              style={{ padding: '6px 12px 6px 28px', border: '1px solid var(--border-color)', borderRadius: '4px', fontSize: 'var(--font-sm)', outline: 'none' }} 
            />
          </div>
          <button className="btn btn-primary hover-lift" onClick={() => navigate('/sales/document')}>
            <Plus size={16} />
            {t('New Sales Order')}
          </button>
        </div>
      </div>

      <div className="panel data-grid-container hover-lift">
        <table className="data-grid">
          <thead>
            <tr>
              <th>{t('Order ID')}</th>
              <th>{t('Customer')}</th>
              <th>{t('Amount')}</th>
              <th>{t('Status')}</th>
              <th>{t('Created Date')}</th>
              <th style={{ textAlign: 'right' }}>{t('Actions')}</th>
            </tr>
          </thead>
          <tbody>
            {orders.map(order => (
              <tr key={order.id}>
                <td style={{ fontWeight: 500, color: 'var(--brand-accent)' }}>{order.id}</td>
                <td>{order.customer}</td>
                <td>¥{order.amount.toLocaleString()}</td>
                <td>
                  <span style={{ 
                    padding: '2px 8px', 
                    borderRadius: '12px', 
                    fontSize: '11px',
                    backgroundColor: order.status === 'Completed' ? '#d1fae5' : order.status === 'Pending' ? '#fef3c7' : '#dbeafe',
                    color: order.status === 'Completed' ? '#065f46' : order.status === 'Pending' ? '#92400e' : '#1e40af'
                  }}>
                    {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
                    {t(order.status as any)}
                  </span>
                </td>
                <td>{order.date}</td>
                <td style={{ textAlign: 'right' }}>
                  <button className="btn btn-secondary" style={{ padding: '4px', marginRight: '4px' }} title={t('View')} onClick={() => navigate('/sales/document')}>
                    <Eye size={14} />
                  </button>
                  <button className="btn btn-secondary" style={{ padding: '4px' }} title={t('Edit')} onClick={() => navigate('/sales/document')}>
                    <Edit2 size={14} />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
