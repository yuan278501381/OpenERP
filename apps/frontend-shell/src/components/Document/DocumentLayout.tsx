/* eslint-disable i18next/no-literal-string */
import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

const generateTraceId = () => Math.random().toString(36).substring(2, 15);
const logger = {
  info: (msg: string, traceId: string) => console.log(`[INFO] [${traceId}] ${msg}`),
  warn: (msg: string, traceId: string) => console.warn(`[WARN] [${traceId}] ${msg}`),
  error: (msg: string, traceId: string) => console.error(`[ERROR] [${traceId}] ${msg}`),
};

interface DocumentLayoutProps {
  title: string;
  children: React.ReactNode;
}

export const DocumentLayout: React.FC<DocumentLayoutProps> = ({ title, children }) => {
  const { t } = useTranslation();
  const [activeSection, setActiveSection] = useState<string>('header');
  const traceId = generateTraceId();

  useEffect(() => {
    logger.info(`DocumentLayout mounted for ${title}`, traceId);
    const handleScroll = () => {
      const sections = ['header', 'lines', 'logistics', 'accounting'];
      for (const section of sections) {
        const el = document.getElementById(section);
        if (el) {
          const rect = el.getBoundingClientRect();
          if (rect.top >= 0 && rect.top < 300) {
            setActiveSection(section);
            break;
          }
        }
      }
    };
    
    const container = document.getElementById('document-scroll-container');
    if (container) {
      container.addEventListener('scroll', handleScroll);
      return () => container.removeEventListener('scroll', handleScroll);
    }
  }, [title, traceId]);

  const scrollTo = (id: string) => {
    logger.info(`Scrolling to section: ${id}`, traceId);
    const el = document.getElementById(id);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth' });
      setActiveSection(id);
    }
  };

  return (
    <div id="document-scroll-container" className="document-layout" style={{ height: '100%', overflowY: 'auto', backgroundColor: 'var(--bg-primary)' }}>
      {/* Action Bar */}
      <div className="action-bar" style={{ position: 'sticky', top: 0, zIndex: 20, display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '12px 24px', backgroundColor: 'rgba(255, 255, 255, 0.75)', backdropFilter: 'blur(16px)', WebkitBackdropFilter: 'blur(16px)', borderBottom: '1px solid var(--border-color)' }}>
        <h2 style={{ margin: 0, fontSize: '18px', fontWeight: 600 }}>{title}</h2>
        <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
          <select className="btn btn-secondary" onChange={(e) => logger.info(`Copy From: ${e.target.value}`, traceId)} defaultValue="">
            <option value="" disabled>{t('Copy From')}</option>
            <option value="quotation">{t('Quotation')}</option>
          </select>
          <select className="btn btn-secondary" onChange={(e) => logger.info(`Copy To: ${e.target.value}`, traceId)} defaultValue="">
            <option value="" disabled>{t('Copy To')}</option>
            <option value="delivery">{t('Delivery')}</option>
            <option value="invoice">{t('Invoice')}</option>
          </select>
          <button className="btn btn-primary active-press" onClick={() => logger.info('Save clicked', traceId)}>{t('Save')}</button>
        </div>
      </div>
      
      {/* Horizontal Anchor Nav */}
      <nav style={{ position: 'sticky', top: '57px', zIndex: 10, backgroundColor: 'rgba(255, 255, 255, 0.75)', backdropFilter: 'blur(16px)', WebkitBackdropFilter: 'blur(16px)', borderBottom: '1px solid var(--border-color)', padding: '0 24px' }}>
        <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'flex', gap: '32px' }}>
          {['header', 'lines', 'logistics', 'accounting'].map((id) => (
            <li key={id} style={{ margin: 0 }}>
              <a
                href={`#${id}`}
                onClick={(e) => { e.preventDefault(); scrollTo(id); }}
                style={{
                  display: 'inline-block',
                  padding: '16px 4px',
                  color: activeSection === id ? 'var(--text-primary)' : 'var(--text-secondary)',
                  fontWeight: activeSection === id ? 600 : 500,
                  fontSize: '14px',
                  textDecoration: 'none',
                  transition: 'all 0.2s ease',
                  borderBottom: activeSection === id ? '2px solid var(--text-primary)' : '2px solid transparent',
                  marginBottom: '-1px'
                }}
              >
                {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
                {t(id.charAt(0).toUpperCase() + id.slice(1) as any)}
              </a>
            </li>
          ))}
        </ul>
      </nav>
      
      {/* Content */}
      <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto', width: '100%' }}>
        {children}
      </div>
    </div>
  );
};
