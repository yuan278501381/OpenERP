/* eslint-disable @typescript-eslint/no-explicit-any */
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

export function JournalEntryList() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  
  const entries = [
    { id: 'JE-2026-001', date: '2026-06-17', account: '1001 Cash', debit: '1000.00', credit: '' },
    { id: 'JE-2026-001', date: '2026-06-17', account: '4001 Sales Revenue', debit: '', credit: '1000.00' },
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{t('General Ledger Entries' as any)}</h1>
          <p className="text-gray-500 mt-1">{t('Manage GL accounting entries.' as any)}</p>
        </div>
        <button type="button" onClick={() => navigate('/finance/document')} className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 shadow-sm transition-colors font-medium">
          {t('New Entry' as any)}
        </button>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-100 overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left font-medium text-gray-500">{t('Entry ID' as any)}</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500">{t('Posting Date' as any)}</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500">{t('Account' as any)}</th>
              <th className="px-6 py-3 text-right font-medium text-gray-500">{t('Debit' as any)}</th>
              <th className="px-6 py-3 text-right font-medium text-gray-500">{t('Credit' as any)}</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {entries.map((entry, idx) => (
              <tr key={idx} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-gray-900 font-medium">{entry.id}</td>
                <td className="px-6 py-4 whitespace-nowrap text-gray-500">{entry.date}</td>
                <td className="px-6 py-4 whitespace-nowrap text-gray-900">{entry.account}</td>
                <td className="px-6 py-4 whitespace-nowrap text-gray-900 text-right">{entry.debit}</td>
                <td className="px-6 py-4 whitespace-nowrap text-gray-900 text-right">{entry.credit}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
