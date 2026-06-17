import React, { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { DocumentLayout } from '../../components/Document/DocumentLayout';
import { DocumentMatrix } from '../../components/Document/DocumentMatrix';

const generateTraceId = () => Math.random().toString(36).substring(2, 15);
const logger = {
  info: (msg: string, traceId: string) => console.log(`[INFO] [${traceId}] ${msg}`),
};

export const QualityDocument: React.FC = () => {
  const { t } = useTranslation();
  const traceId = generateTraceId();

  useEffect(() => {
    logger.info('QualityDocument initialized', traceId);
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
    <DocumentLayout title={t('Inspection Lot')}>
      <div id="header" style={sectionStyle}>
        <h3 style={sectionHeaderStyle}>{t('Header')}</h3>
        <div style={gridStyle}>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Vendor/Plant')}</label>
            <input type="text" defaultValue="P001 - Main Plant" style={inputStyle} />
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Inspection Date')}</label>
            <input type="date" defaultValue="2026-06-17" style={inputStyle} />
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Material')}</label>
            <input type="text" defaultValue="MAT-001" style={inputStyle} />
          </div>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Status')}</label>
            <input type="text" defaultValue={t('In Inspection')} style={{...inputStyle, backgroundColor: '#f3f4f6', color: '#6b7280'}} readOnly />
          </div>
        </div>
      </div>

      <div id="lines" style={sectionStyle}>
        <h3 style={sectionHeaderStyle}>{t('Inspection Characteristics')}</h3>
        <DocumentMatrix />
      </div>

      <div id="logistics" style={sectionStyle}>
        <h3 style={sectionHeaderStyle}>{t('Logistics')}</h3>
        <div style={gridStyle}>
          <div style={fieldStyle}>
            <label style={labelStyle}>{t('Lab Location')}</label>
            <input type="text" defaultValue="LAB-01" style={inputStyle} />
          </div>
        </div>
      </div>
    </DocumentLayout>
  );
};

