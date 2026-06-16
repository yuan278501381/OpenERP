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
        <div className="logo-dot"></div>
        <h2>OpenERP</h2>
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
