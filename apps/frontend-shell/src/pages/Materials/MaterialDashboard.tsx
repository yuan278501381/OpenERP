import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, Filter, Download, MoreHorizontal } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useMaterials } from '../../hooks/useMaterials';
import './MaterialDashboard.css';

export const MaterialDashboard: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { materials, loading, error } = useMaterials();

  return (
    <div className="material-dashboard">
      <div className="dashboard-header">
        <div className="header-title">
          <h1>{t('Material Master Data')}</h1>
          <p className="subtitle">{t('Material Subtitle')}</p>
        </div>
        <div className="header-actions">
          <button className="btn btn-secondary">
            <Filter size={14} />
            {t('Filter')}
          </button>
          <button className="btn btn-secondary">
            <Download size={14} />
            {t('Export')}
          </button>
          <button type="button" onClick={() => navigate('/master-data/document')} className="btn btn-primary">
            <Plus size={14} />
            {t('New Material')}
          </button>
        </div>
      </div>

      <div className="panel data-grid-container hover-lift">
        {loading ? (
          <div className="loading-state">{t('Loading material data...')}</div>
        ) : error ? (
          <div className="error-state">{error}</div>
        ) : (
          <table className="data-grid">
            <thead>
              <tr>
                <th>{t('Material Code')}</th>
                <th>{t('Name / Description')}</th>
                <th>{t('Category')}</th>
                <th>{t('Stock')}</th>
                <th>{t('Unit')}</th>
                <th>{t('Standard Price')}</th>
                <th>{t('Status')}</th>
                <th>{t('Last Updated')}</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {materials.map((mat) => (
                <tr key={mat.id}>
                  <td className="cell-code">{mat.code}</td>
                  <td className="cell-name">{mat.name}</td>
                  <td>{mat.category}</td>
                  <td className="cell-number">{mat.stock.toLocaleString()}</td>
                  <td>{mat.unit}</td>
                  <td className="cell-number">${mat.price.toFixed(2)}</td>
                  <td>
                    <span className={`status-badge ${mat.status}`}>
                      {mat.status}
                    </span>
                  </td>
                  <td>{mat.lastUpdated}</td>
                  <td className="cell-action">
                    <button className="icon-btn">
                      <MoreHorizontal size={16} />
                    </button>
                  </td>
                </tr>
              ))}
              {materials.length === 0 && (
                <tr>
                  <td colSpan={9} className="empty-state">{t('No materials found.')}</td>
                </tr>
              )}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};
