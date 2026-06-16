/* eslint-disable i18next/no-literal-string */
import React, { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { DocumentLayout } from '../../components/Document/DocumentLayout';
import { DocumentMatrix } from '../../components/Document/DocumentMatrix';

const generateTraceId = () => Math.random().toString(36).substring(2, 15);
const logger = {
  info: (msg: string, traceId: string) => console.log(`[INFO] [${traceId}] ${msg}`),
};

export const SalesOrderDocument: React.FC = () => {
  const { t } = useTranslation();
  const traceId = generateTraceId();

  useEffect(() => {
    logger.info('SalesOrderDocument initialized', traceId);
  }, [traceId]);

  const sectionStyle = {
    padding: '24px',
    backgroundColor: '#fff',
    borderRadius: '6px',
    marginBottom: '24px',
    boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
    border: '1px solid #e5e7eb'
  };

  const sectionHeaderStyle = {
    marginTop: 0, 
    marginBottom: '20px', 
    borderBottom: '1px solid #f3f4f6', 
    paddingBottom: '12px',
    fontSize: '16px',
    fontWeight: 600,
    color: '#111827'
  };

  const gridStyle = {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
    gap: '24px'
  };

  const fieldStyle = {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '6px'
  };

  const labelStyle = {
    fontSize: '12px',
    color: '#6b7280',
    fontWeight: 500
  };

  const inputStyle = {
    padding: '8px 12px',
    border: '1px solid #d1d5db',
    borderRadius: '4px',
    fontSize: '13px',
    outline: 'none',
    transition: 'border-color 0.2s',
    backgroundColor: '#fff'
  };

  return (
    <DocumentLayout title={t('Sales Orders')}>
      <div id="header" style={sectionStyle}>
        <h3 style={sectionHeaderStyle}>{t('Header')}</h3>
        <div style={gridStyle}>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Sold-to Party')}</label>
            <input type="text" defaultValue="C0001 - Acme Corp" style={inputStyle} />
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Posting Date')}</label>
            <input type="date" defaultValue="2026-06-17" style={inputStyle} />
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>
              {t('Related Party')} 
              <span style={{ color: '#ef4444', fontSize: '11px', marginLeft: '6px', fontWeight: 600 }}>{t('Intercompany')}</span>
            </label>
            <select style={inputStyle} defaultValue="Y">
              <option value="N">{t('No')}</option>
              <option value="Y">{t('Yes - Affiliate')}</option>
            </select>
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('End Customer')}</label>
            <input type="text" placeholder={t('Distributor indirect sales flow...')} style={inputStyle} />
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Status')}</label>
            <input type="text" defaultValue={t('Pending')} style={{...inputStyle, backgroundColor: '#f3f4f6', color: '#6b7280'}} readOnly />
          </div>
        </div>
      </div>

      <div id="lines" style={sectionStyle}>
        <h3 style={sectionHeaderStyle}>{t('Lines')}</h3>
        <DocumentMatrix />
      </div>

      <div id="logistics" style={sectionStyle}>
        <h3 style={sectionHeaderStyle}>{t('Logistics')}</h3>
        <div style={gridStyle}>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Ship-to Party')}</label>
            <textarea style={{...inputStyle, minHeight: '80px', resize: 'vertical'}} defaultValue="123 Industrial Way&#10;Tech City, TC 90210" />
          </div>
          <div style={{ ...fieldStyle, justifyContent: 'flex-start' }}>
            <div style={fieldStyle}>
              <label style={labelStyle}>{t('Shipping Type')}</label>
              <select style={inputStyle} defaultValue="motor">
                <option value="motor">{t('Motor Express')}</option>
                <option value="air">{t('Air Freight')}</option>
              </select>
            </div>
            <div style={{ ...fieldStyle, marginTop: '16px' }}>
              <label style={labelStyle}>{t('Tracking Number')}</label>
              <input type="text" placeholder="1Z9999999999999999" style={inputStyle} />
            </div>
          </div>
        </div>
      </div>

      <div id="accounting" style={sectionStyle}>
        <h3 style={sectionHeaderStyle}>{t('Accounting')}</h3>
        <div style={gridStyle}>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Bill-to Party')}</label>
            <input type="text" defaultValue="C0001 - Acme Corp" style={inputStyle} />
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Payer')}</label>
            <input type="text" defaultValue="C0001 - Acme Corp" style={inputStyle} />
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Payment Terms')}</label>
            <select style={inputStyle} defaultValue="net30">
              <option value="net30">{t('Net 30')}</option>
              <option value="cod">{t('Cash on Delivery')}</option>
            </select>
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Currency')}</label>
            <select style={inputStyle} defaultValue="CNY">
              <option value="CNY">CNY - 人民币</option>
              <option value="USD">USD - US Dollar</option>
            </select>
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Journal Remark')}</label>
            <input type="text" defaultValue="Sales Order - C0001" style={inputStyle} />
          </div>
        </div>
      </div>
    </DocumentLayout>
  );
};
