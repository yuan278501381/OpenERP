import React from 'react';
import type { Task } from '../api/task';
import { Clock, FileText, CheckCircle, AlertTriangle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import './TaskCard.css';

interface TaskCardProps {
  task: Task;
}

export const TaskCard: React.FC<TaskCardProps> = ({ task }) => {
  const { t } = useTranslation();
  const isWarning = task.sla_status === 'Warning';
  const isCritical = task.sla_status === 'Critical';

  return (
    <div className="panel task-card fade-in">
      <div className="task-header">
        <div className="task-type">
          <FileText size={16} />
          {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
          <span>{t(task.type as any)}</span>
        </div>
        <span className={`status-badge ${isWarning || isCritical ? 'warning' : 'normal'}`}>
          {isWarning || isCritical ? <AlertTriangle size={12}/> : <CheckCircle size={12}/>}
          {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
          {t(task.sla_status as any)}
        </span>
      </div>
      
      {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
      <h3 className="task-title">{t(task.title as any)}</h3>
      
      <div className="task-footer">
        <div className="task-node">
          <Clock size={14} className="node-icon" />
          {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
          <span>{t('Current Node')}{t(task.node as any)}</span>
        </div>
        <button className="btn-primary">{t('Process Now')}</button>
      </div>
    </div>
  );
};
