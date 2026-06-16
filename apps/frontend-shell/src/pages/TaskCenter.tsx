import React, { useEffect, useState } from 'react';
import { fetchTasks, type Task } from '../api/task';
import { TaskCard } from '../components/TaskCard';
import './TaskCenter.css';

export const TaskCenter: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchTasks().then(data => {
      setTasks(data);
      setLoading(false);
    }).catch(err => {
      console.error(err);
      setLoading(false);
    });
  }, []);

  return (
    <div className="task-center-container">
      <header className="page-header fade-in">
        <div className="header-content">
          <h1>统一任务中枢 (Task Center)</h1>
          <p>事找人 · 智能协同 · 全局可视</p>
        </div>
      </header>
      
      <main className="task-grid">
        {loading ? (
          <div className="loading-state">
            <div className="spinner"></div>
            <p>正在同步全系统待办...</p>
          </div>
        ) : (
          tasks.map(task => (
            <TaskCard key={task.task_id} task={task} />
          ))
        )}
      </main>
    </div>
  );
};
