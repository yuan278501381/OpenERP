/* eslint-disable @typescript-eslint/no-explicit-any */
import { useTranslation } from 'react-i18next';

export function EmployeeMaster() {
  const { t } = useTranslation();
  
  const employees = [
    { id: 'E-001', name: 'John Doe', department: 'IT', position: 'Senior Developer' },
    { id: 'E-002', name: 'Jane Smith', department: 'Finance', position: 'CFO' },
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{t('Employee Master' as any)}</h1>
          <p className="text-gray-500 mt-1">{t('Manage organization tree and employee profiles.' as any)}</p>
        </div>
        <button className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 shadow-sm transition-colors font-medium">
          {t('New Employee' as any)}
        </button>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-100 overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left font-medium text-gray-500">{t('Employee ID' as any)}</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500">{t('Name' as any)}</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500">{t('Department' as any)}</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500">{t('Position' as any)}</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {employees.map((emp) => (
              <tr key={emp.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-gray-900 font-medium">{emp.id}</td>
                <td className="px-6 py-4 whitespace-nowrap text-gray-500">{emp.name}</td>
                <td className="px-6 py-4 whitespace-nowrap text-gray-900">{emp.department}</td>
                <td className="px-6 py-4 whitespace-nowrap text-gray-900">{emp.position}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
