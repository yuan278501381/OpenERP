import React from 'react';
import { Bell, Search, User } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import './Layout.css';

export const Header: React.FC = () => {
  const { t } = useTranslation();
  return (
    <header className="header">
      <div className="header-search">
        <Search size={16} className="search-icon" />
        <input type="text" placeholder={t('Search Placeholder')} className="search-input" />
      </div>
      <div className="header-actions">
        <button className="header-btn">
          <Bell size={18} />
        </button>
        <div className="user-profile">
          <div className="avatar">
            <User size={18} />
          </div>
          <div className="user-info">
            <span className="user-name">{t('Admin User')}</span>
            <span className="user-role">{t('System Admin')}</span>
          </div>
        </div>
      </div>
    </header>
  );
};
