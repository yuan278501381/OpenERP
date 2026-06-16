import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Package, ArrowDownToLine, History } from 'lucide-react';

export const GoodsMovement: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  
  const [movements] = useState([
    { id: 'MVT-001', type: 'Receipt', material: 'MAT-001', qty: 500, date: '2026-06-16 10:30' },
    { id: 'MVT-002', type: 'Issue', material: 'MAT-005', qty: -20, date: '2026-06-16 14:15' },
  ]);

  return (
    <div className="p-4" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--spacing-lg)' }}>
      <div className="panel p-4" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <h1 style={{ fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Package size={20} style={{ color: 'var(--status-normal)' }} />
            {t('Goods Movement')}
          </h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: 'var(--font-sm)', marginTop: '4px' }}>
            {t('Warehouse receipts and issues')}
          </p>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: 'var(--spacing-lg)' }}>
        {/* Form Panel */}
        <div className="panel p-4 hover-lift" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--spacing-md)' }}>
          <h2 style={{ fontSize: '14px', fontWeight: 600, marginBottom: '8px' }}>{t('New Movement')}</h2>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
            <label style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>{t('Movement Type')}</label>
            <select className="form-input">
              <option>{t('Receipt')}</option>
              <option>{t('Issue')}</option>
            </select>
          </div>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
            <label style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>{t('Material Code')}</label>
            <input type="text" placeholder={t('Material Placeholder')} style={{ padding: '8px', border: '1px solid var(--border-color)', borderRadius: '4px', fontSize: 'var(--font-sm)', outline: 'none' }} />
          </div>

          <div style={{ display: 'flex', gap: '8px' }}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '4px', flex: 1 }}>
              <label style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>{t('Quantity')}</label>
              <input type="number" defaultValue={1} style={{ padding: '8px', border: '1px solid var(--border-color)', borderRadius: '4px', fontSize: 'var(--font-sm)', outline: 'none' }} />
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '4px', flex: 1 }}>
              <label style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>{t('Warehouse')}</label>
              <input type="text" defaultValue="WH01" style={{ padding: '8px', border: '1px solid var(--border-color)', borderRadius: '4px', fontSize: 'var(--font-sm)', outline: 'none' }} />
            </div>
          </div>

          <button type="button" onClick={() => navigate('/inventory/document')} className="btn btn-primary" style={{ marginTop: '8px', padding: '10px', justifyContent: 'center' }}>
            <ArrowDownToLine size={16} />
            {t('Submit Movement')}
          </button>
        </div>

        {/* History Panel */}
        <div className="panel p-0 hover-lift" style={{ display: 'flex', flexDirection: 'column' }}>
          <div style={{ padding: '16px', borderBottom: '1px solid var(--border-color)', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <History size={16} style={{ color: 'var(--text-secondary)' }} />
            <h2 style={{ fontSize: '14px', fontWeight: 600 }}>{t('Recent Movements')}</h2>
          </div>
          <div className="data-grid-container" style={{ border: 'none', borderRadius: '0 0 6px 6px' }}>
            <table className="data-grid">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>{t('Movement Type')}</th>
                  <th>{t('Material Code')}</th>
                  <th style={{ textAlign: 'right' }}>{t('Quantity')}</th>
                  <th>{t('Created Date')}</th>
                </tr>
              </thead>
              <tbody>
                {movements.map(m => (
                  <tr key={m.id}>
                    <td style={{ fontWeight: 500 }}>{m.id}</td>
                    <td>
                      <span className={`status-badge ${m.type.toLowerCase()}`}>
                        {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
                        {t(m.type as any)}
                      </span>
                    </td>
                    <td>{m.material}</td>
                    <td style={{ textAlign: 'right', fontWeight: 500, color: m.qty > 0 ? 'var(--status-normal)' : 'var(--status-critical)' }}>
                      {m.qty > 0 ? `+${m.qty}` : m.qty}
                    </td>
                    <td style={{ color: 'var(--text-secondary)' }}>{m.date}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};
