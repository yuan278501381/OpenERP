import React from 'react';
import { NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { LayoutDashboard, Package, ShoppingCart, Truck, Landmark, Globe } from 'lucide-react';
import './Layout.css';

export const Sidebar: React.FC = () => {
  const { t, i18n } = useTranslation();

  const toggleLanguage = () => {
    i18n.changeLanguage(i18n.language === 'zh' ? 'en' : 'zh');
  };

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <svg width="28" height="28" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M16 2L2 9.5V22.5L16 30L30 22.5V9.5L16 2Z" fill="#002FA7" />
          <path d="M16 8L7 13V21L16 26L25 21V13L16 8Z" fill="#4D88FF" />
          <path d="M16 12L12 15V19L16 22L20 19V15L16 12Z" fill="#FFFFFF" />
        </svg>
        {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
        <h2 style={{ margin: 0 }}>{t('OpenERP' as any)}</h2>
      </div>
      <nav className="sidebar-nav">
        <NavLink to="/" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`} end>
          <LayoutDashboard size={18} />
          <span>{t('Dashboard')}</span>
        </NavLink>
        <NavLink to="/master-data" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <Package size={18} />
          <span>{t('Master Data')}</span>
        </NavLink>
        <NavLink to="/sales" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <ShoppingCart size={18} />
          <span>{t('Sales & Distribution')}</span>
        </NavLink>
        <NavLink to="/purchase" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <Truck size={18} />
          <span>{t('Purchasing')}</span>
        </NavLink>
        <NavLink to="/inventory" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <Package size={18} />
          <span>{t('Inventory Management')}</span>
        </NavLink>
        <NavLink to="/production" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <Package size={18} />
          <span>{t('Production (PP)')}</span>
        </NavLink>
        <NavLink to="/quality" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <Package size={18} />
          <span>{t('Quality Management')}</span>
        </NavLink>
        <NavLink to="/maintenance" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <Package size={18} />
          <span>{t('Plant Maintenance')}</span>
        </NavLink>
        <NavLink to="/finance" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <Landmark size={18} />
          <span>{t('Finance (FI/CO)')}</span>
        </NavLink>
        <NavLink to="/hr" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <Package size={18} />
          <span>{t('Human Resources')}</span>
        </NavLink>
        <NavLink to="/system" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
          <Package size={18} />
          <span>{t('System & Org')}</span>
        </NavLink>
      </nav>
      <div style={{ padding: '16px', marginTop: 'auto', borderTop: '1px solid var(--border-color)' }}>
        <button className="btn btn-secondary" style={{ width: '100%', justifyContent: 'flex-start' }} onClick={toggleLanguage}>
          <Globe size={16} />
          <span>{t('Switch Language')}</span>
        </button>
      </div>
    </aside>
  );
};
