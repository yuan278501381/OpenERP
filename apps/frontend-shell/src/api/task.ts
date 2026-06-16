import axios from 'axios';

export interface Task {
  task_id: string;
  title: string;
  type: string;
  node: string;
  sla_status: 'Normal' | 'Warning' | 'Critical';
}

export const fetchTasks = async (): Promise<Task[]> => {
  // 借助 vite proxy，请求直接打向 /openerp/v1/tasks
  const response = await axios.get('/openerp/v1/tasks', {
    headers: {
      'X-Tenant-ID': 'TENANT-001',
      'X-User-ID': 'USER-1024',
    }
  });
  return response.data.data;
};
