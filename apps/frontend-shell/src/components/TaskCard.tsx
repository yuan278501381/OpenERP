import React from 'react';
import type { Task } from '../api/task';
import { Clock, FileText, CheckCircle, AlertTriangle } from 'lucide-react';
import './TaskCard.css';

interface TaskCardProps {
  task: Task;
}

export const TaskCard: React.FC<TaskCardProps> = ({ task }) => {
  const isWarning = task.sla_status === 'Warning';
  const isCritical = task.sla_status === 'Critical';

  return (
    <div className="glass-panel task-card fade-in">
      <div className="task-header">
        <div className="task-type">
          <FileText size={16} />
          <span>{task.type}</span>
        </div>
        <span className={`status-badge ${isWarning || isCritical ? 'warning' : 'normal'}`}>
          {isWarning || isCritical ? <AlertTriangle size={12}/> : <CheckCircle size={12}/>}
          {task.sla_status}
        </span>
      </div>
      
      <h3 className="task-title">{task.title}</h3>
      
      <div className="task-footer">
        <div className="task-node">
          <Clock size={14} className="node-icon" />
          <span>当前卡点: {task.node}</span>
        </div>
        <button className="btn-primary">立即处理</button>
      </div>
    </div>
  );
};
