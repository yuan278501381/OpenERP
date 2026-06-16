import React, { useEffect, useState } from 'react';
import { fetchTasks, type Task } from '../api/task';
import { useTranslation } from 'react-i18next';
import { TaskCard } from '../components/TaskCard';
import './TaskCenter.css';

export const TaskCenter: React.FC = () => {
  const { t } = useTranslation();
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
          <h1>{t('Task Center')}</h1>
          <p>{t('Task Subtitle')}</p>
        </div>
      </header>
      
      <main className="task-grid">
        {loading ? (
          <div className="loading-state">
            <div className="spinner"></div>
            <p>{t('Syncing')}</p>
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
