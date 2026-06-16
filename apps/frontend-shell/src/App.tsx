import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { MainLayout } from './components/Layout/MainLayout';
import { TaskCenter } from './pages/TaskCenter';
import { PurchaseRequest } from './pages/Purchase/PurchaseRequest';
import { MaterialDashboard } from './pages/Materials/MaterialDashboard';
import { SalesOrderList } from './pages/Sales/SalesOrderList';
import { GoodsMovement } from './pages/Inventory/GoodsMovement';
import { ProductionBoard } from './pages/Production/ProductionBoard';
import { InspectionLot } from './pages/Quality/InspectionLot';
import { JournalEntryList } from './pages/Finance/JournalEntryList';
import { EmployeeMaster } from './pages/HR/EmployeeMaster';
import { SalesOrderDocument } from './pages/Sales/SalesOrderDocument';
import { PurchaseDocument } from './pages/Purchase/PurchaseDocument';
import { InventoryDocument } from './pages/Inventory/InventoryDocument';
import { ProductionDocument } from './pages/Production/ProductionDocument';
import { QualityDocument } from './pages/Quality/QualityDocument';
import { FinanceDocument } from './pages/Finance/FinanceDocument';
import { MaterialDocument } from './pages/Materials/MaterialDocument';
import { useTranslation } from 'react-i18next';
import './App.css';

function App() {
  const { t } = useTranslation();
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MainLayout />}>
          <Route index element={<TaskCenter />} />
          <Route path="master-data" element={<MaterialDashboard />} />
          <Route path="master-data/document" element={<MaterialDocument />} />
          <Route path="sales" element={<SalesOrderList />} />
          <Route path="sales/document" element={<SalesOrderDocument />} />
          <Route path="purchase" element={<PurchaseRequest />} />
          <Route path="purchase/document" element={<PurchaseDocument />} />
          <Route path="inventory" element={<GoodsMovement />} />
          <Route path="inventory/document" element={<InventoryDocument />} />
          <Route path="production" element={<ProductionBoard />} />
          <Route path="production/document" element={<ProductionDocument />} />
          <Route path="quality" element={<InspectionLot />} />
          <Route path="quality/document" element={<QualityDocument />} />
          <Route path="maintenance" element={<div className="panel p-4">{t('Plant Maintenance')}{t('Coming Soon')}</div>} />
          <Route path="finance" element={<JournalEntryList />} />
          <Route path="finance/document" element={<FinanceDocument />} />
          <Route path="hr" element={<EmployeeMaster />} />
          <Route path="system" element={<div className="panel p-4">{t('System & Org')}{t('Coming Soon')}</div>} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;

