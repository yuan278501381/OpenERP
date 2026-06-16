import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { MainLayout } from './components/Layout/MainLayout';
import { TaskCenter } from './pages/TaskCenter';
import { PurchaseRequest } from './pages/PurchaseRequest';
import { MaterialDashboard } from './pages/Materials/MaterialDashboard';
import { useTranslation } from 'react-i18next';
import './App.css'; // Keep if any global app-specific styles left, otherwise can be empty

function App() {
  const { t } = useTranslation();
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MainLayout />}>
          <Route index element={<TaskCenter />} />
          <Route path="master-data" element={<MaterialDashboard />} />
          <Route path="sales" element={<div className="panel p-4">{t('Sales & Distribution')}{t('Coming Soon')}</div>} />
          <Route path="purchase" element={<PurchaseRequest />} />
          <Route path="inventory" element={<div className="panel p-4">{t('Inventory Management')}{t('Coming Soon')}</div>} />
          <Route path="production" element={<div className="panel p-4">{t('Production (PP)')}{t('Coming Soon')}</div>} />
          <Route path="quality" element={<div className="panel p-4">{t('Quality Management')}{t('Coming Soon')}</div>} />
          <Route path="maintenance" element={<div className="panel p-4">{t('Plant Maintenance')}{t('Coming Soon')}</div>} />
          <Route path="finance" element={<div className="panel p-4">{t('Finance (FI/CO)')}{t('Coming Soon')}</div>} />
          <Route path="hr" element={<div className="panel p-4">{t('Human Resources')}{t('Coming Soon')}</div>} />
          <Route path="system" element={<div className="panel p-4">{t('System & Org')}{t('Coming Soon')}</div>} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
