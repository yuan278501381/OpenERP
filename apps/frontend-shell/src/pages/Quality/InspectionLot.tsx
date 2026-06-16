import React from 'react';
import { useNavigate } from 'react-router-dom';
import { ClipboardCheck, Search, CheckCircle, XCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import './InspectionLot.css';

interface InspectionLotData {
  id: string;
  lotNo: string;
  material: string;
  batch: string;
  quantity: number;
  status: 'Pending' | 'Passed' | 'Failed';
  date: string;
}

const mockLots: InspectionLotData[] = [
  { id: '1', lotNo: 'QA-2026-1001', material: 'CPU Intel Xeon', batch: 'B-001', quantity: 100, status: 'Pending', date: '2026-06-17' },
  { id: '2', lotNo: 'QA-2026-1002', material: 'Motherboard X99', batch: 'B-002', quantity: 50, status: 'Passed', date: '2026-06-16' },
  { id: '3', lotNo: 'QA-2026-1003', material: 'Cooling Fan 120mm', batch: 'B-003', quantity: 500, status: 'Failed', date: '2026-06-15' },
];

export const InspectionLot: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className="inspection-lot-page">
      <div className="dashboard-header">
        <div className="header-title">
          <h1>{t('Inspection Lots')}</h1>
          <p className="subtitle">{t('Inspection Subtitle')}</p>
        </div>
        <div className="header-actions">
          <button className="btn btn-secondary">
            <Search size={14} />
            {t('Filter')}
          </button>
          <button type="button" onClick={() => navigate('/quality/document')} className="btn btn-primary">
            <ClipboardCheck size={14} />
            {t('Record Results')}
          </button>
        </div>
      </div>

      <div className="panel data-grid-container hover-lift">
        <table className="data-grid">
          <thead>
            <tr>
              <th>{t('Inspection Lots')}</th>
              <th>{t('Material Code')}</th>
              <th>{t('Batch')}</th>
              <th>{t('Quantity')}</th>
              <th>{t('Status')}</th>
              <th>{t('Created Date')}</th>
              <th>{t('Result')}</th>
            </tr>
          </thead>
          <tbody>
            {mockLots.map(lot => (
              <tr key={lot.id}>
                <td className="cell-code">{lot.lotNo}</td>
                <td className="cell-name">{lot.material}</td>
                <td>{lot.batch}</td>
                <td className="cell-number">{lot.quantity}</td>
                <td>
                  <span className={`status-badge ${lot.status.toLowerCase()}`}>
                    {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
                    {t(lot.status as any)}
                  </span>
                </td>
                <td>{lot.date}</td>
                <td className="cell-action">
                  {lot.status === 'Pending' ? (
                    <button type="button" onClick={() => navigate('/quality/document')} className="btn btn-secondary btn-small">
                      {t('Record Results')}
                    </button>
                  ) : lot.status === 'Passed' ? (
                    <span className="result-pass"><CheckCircle size={16}/> {t('Pass')}</span>
                  ) : (
                    <span className="result-fail"><XCircle size={16}/> {t('Fail')}</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
