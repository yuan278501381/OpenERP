import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Play, CheckCircle, Factory, Settings, AlertCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import './ProductionBoard.css';

interface ProductionOrder {
  id: string;
  orderNo: string;
  product: string;
  quantity: number;
  progress: number;
  status: 'Planned' | 'In Production' | 'Quality Check' | 'Completed';
  workCenter: string;
}

const mockOrders: ProductionOrder[] = [
  { id: '1', orderNo: 'PO-2026-001', product: 'Server Rack 42U', quantity: 50, progress: 0, status: 'Planned', workCenter: 'WC-Assembly' },
  { id: '2', orderNo: 'PO-2026-002', product: 'Power Supply 1000W', quantity: 200, progress: 0, status: 'Planned', workCenter: 'WC-Electronics' },
  { id: '3', orderNo: 'PO-2026-003', product: 'Cooling Fan 120mm', quantity: 500, progress: 45, status: 'In Production', workCenter: 'WC-Injection' },
  { id: '4', orderNo: 'PO-2026-004', product: 'Motherboard X99', quantity: 100, progress: 80, status: 'In Production', workCenter: 'WC-SMT' },
  { id: '5', orderNo: 'PO-2026-005', product: 'CPU Intel Xeon', quantity: 100, progress: 100, status: 'Quality Check', workCenter: 'WC-Testing' },
  { id: '6', orderNo: 'PO-2026-006', product: 'RAM 32GB DDR4', quantity: 400, progress: 100, status: 'Completed', workCenter: 'WC-Testing' },
];

export const ProductionBoard: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const columns: ('Planned' | 'In Production' | 'Quality Check' | 'Completed')[] = [
    'Planned', 'In Production', 'Quality Check', 'Completed'
  ];

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'Planned': return <Settings size={16} className="text-gray" />;
      case 'In Production': return <Factory size={16} className="text-blue" />;
      case 'Quality Check': return <AlertCircle size={16} className="text-orange" />;
      case 'Completed': return <CheckCircle size={16} className="text-green" />;
      default: return null;
    }
  };

  return (
    <div className="production-board">
      <div className="dashboard-header">
        <div className="header-title">
          <h1>{t('Shop Floor Board')}</h1>
          <p className="subtitle">{t('Production Subtitle')}</p>
        </div>
        <div className="header-actions">
          <button type="button" onClick={() => navigate('/production/document')} className="btn btn-primary">
            <Play size={14} />
            {t('Start Production')}
          </button>
        </div>
      </div>

      <div className="kanban-board">
        {columns.map(col => (
          <div key={col} className="kanban-column">
            <div className="kanban-column-header">
              <h3>{t(col)}</h3>
              <span className="kanban-count">
                {mockOrders.filter(o => o.status === col).length}
              </span>
            </div>
            <div className="kanban-cards">
              {mockOrders.filter(o => o.status === col).map(order => (
                <div key={order.id} className="kanban-card hover-lift">
                  <div className="card-header">
                    <span className="order-no">{order.orderNo}</span>
                    {getStatusIcon(order.status)}
                  </div>
                  <div className="card-body">
                    <h4 className="product-name">{order.product}</h4>
                    <p className="work-center">{t('Work Center')}: {order.workCenter}</p>
                    <p className="quantity">{t('Quantity')}: {order.quantity}</p>
                  </div>
                  <div className="card-footer">
                    <div className="progress-bar-bg">
                      <div className="progress-bar-fill" style={{ width: `${order.progress}%` }}></div>
                    </div>
                    <span className="progress-text">{order.progress}%</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
