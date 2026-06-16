import { useState } from 'react';
import { TaskCenter } from './pages/TaskCenter';
import { PurchaseRequest } from './pages/PurchaseRequest';
import { LayoutDashboard, FilePlus } from 'lucide-react';
import './App.css';

function App() {
  const [activeTab, setActiveTab] = useState<'tasks' | 'purchase'>('tasks');

  return (
    <div className="app-layout">
      <nav className="glass-sidebar">
        <div className="logo-area">
          <div className="logo-dot"></div>
          <h1>OpenERP</h1>
        </div>
        
        <div className="nav-menu">
          <button 
            className={`nav-item ${activeTab === 'tasks' ? 'active' : ''}`}
            onClick={() => setActiveTab('tasks')}
          >
            <LayoutDashboard size={20} />
            <span>统一待办大厅</span>
          </button>
          
          <button 
            className={`nav-item ${activeTab === 'purchase' ? 'active' : ''}`}
            onClick={() => setActiveTab('purchase')}
          >
            <FilePlus size={20} />
            <span>采购申请 (发单)</span>
          </button>
        </div>
      </nav>

      <main className="app-content">
        {activeTab === 'tasks' && <TaskCenter />}
        {activeTab === 'purchase' && <PurchaseRequest />}
      </main>
    </div>
  )
}

export default App
