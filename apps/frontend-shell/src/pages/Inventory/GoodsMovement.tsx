import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { AlertCircle, ArrowDownToLine, CheckCircle2, History, Package, RefreshCw } from 'lucide-react';
import {
  fetchGoodsMovements,
  submitGoodsMovement,
  type GoodsMovementRecord,
  type JournalEntryRecord,
  type MovementType,
} from '../../api/inventory';

export const GoodsMovement: React.FC = () => {
  const { t } = useTranslation();
  const [movements, setMovements] = useState<GoodsMovementRecord[]>([]);
  const [form, setForm] = useState({
    mvmtType: receiptMovementType,
    itemCode: 'MAT-001',
    qty: '1',
    unitCost: '4.50',
    costCenter: 'WH01',
  });
  const [loading, setLoading] = useState(false);
  const [historyLoading, setHistoryLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastJournalEntry, setLastJournalEntry] = useState<JournalEntryRecord | null>(null);

  const loadMovements = useCallback(async () => {
    setHistoryLoading(true);
    try {
      const data = await fetchGoodsMovements();
      setMovements(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : t('Movement Load Failed'));
    } finally {
      setHistoryLoading(false);
    }
  }, [t]);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      loadMovements();
    }, 0);
    return () => window.clearTimeout(timer);
  }, [loadMovements]);

  const canSubmit = useMemo(() => {
    const qty = Number(form.qty);
    const unitCost = Number(form.unitCost);
    return Boolean(form.itemCode.trim()) && qty > 0 && (form.mvmtType === 'Issue' || unitCost > 0);
  }, [form]);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!canSubmit) {
      return;
    }

    setLoading(true);
    setError(null);
    setLastJournalEntry(null);
    try {
      const result = await submitGoodsMovement({
        mvmtType: form.mvmtType,
        itemCode: form.itemCode.trim(),
        qty: Number(form.qty),
        unitCost: form.mvmtType === 'Receipt' ? Number(form.unitCost) : undefined,
        costCenter: form.costCenter.trim() || undefined,
      });
      setLastJournalEntry(result.journalEntry);
      await loadMovements();
    } catch (err: unknown) {
      setError(extractErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-4" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--spacing-lg)' }}>
      <div className="panel p-4" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <h1 style={{ fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Package size={20} style={{ color: 'var(--status-normal)' }} />
            {t('Goods Movement')}
          </h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: 'var(--font-sm)', marginTop: '4px' }}>
            {t('Warehouse receipts and issues')}
          </p>
        </div>
      </div>

      {lastJournalEntry && (
        <div className="panel p-4" style={{ borderColor: '#86efac', background: '#f0fdf4', display: 'flex', alignItems: 'center', gap: '8px' }}>
          <CheckCircle2 size={18} style={{ color: 'var(--status-success)' }} />
          <span style={{ fontSize: 'var(--font-sm)', color: '#166534' }}>
            {t('Movement Posted')} {lastJournalEntry.EntryNo} · {t('Amount')}: {lastJournalEntry.TotalAmount.toFixed(2)}
          </span>
        </div>
      )}

      {error && (
        <div className="panel p-4" style={{ borderColor: '#fecaca', background: '#fef2f2', display: 'flex', alignItems: 'center', gap: '8px' }}>
          <AlertCircle size={18} style={{ color: 'var(--status-critical)' }} />
          <span style={{ fontSize: 'var(--font-sm)', color: '#991b1b' }}>{error}</span>
        </div>
      )}

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: 'var(--spacing-lg)' }}>
        <form onSubmit={handleSubmit} className="panel p-4 hover-lift" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--spacing-md)' }}>
          <h2 style={{ fontSize: '14px', fontWeight: 600, marginBottom: '8px' }}>{t('New Movement')}</h2>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
            <label style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>{t('Movement Type')}</label>
            <select
              className="form-input"
              value={form.mvmtType}
              onChange={(event) => setForm(current => ({ ...current, mvmtType: event.target.value as MovementType }))}
              style={fieldStyle}
            >
              <option value={receiptMovementType}>{t('Receipt')}</option>
              <option value={issueMovementType}>{t('Issue')}</option>
            </select>
          </div>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
            <label style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>{t('Material Code')}</label>
            <input
              type="text"
              placeholder={t('Material Placeholder')}
              value={form.itemCode}
              onChange={(event) => setForm(current => ({ ...current, itemCode: event.target.value }))}
              style={fieldStyle}
            />
          </div>

          <div style={{ display: 'flex', gap: '8px' }}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '4px', flex: 1 }}>
              <label style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>{t('Quantity')}</label>
              <input
                type="number"
                min="0.0001"
                step="0.0001"
                value={form.qty}
                onChange={(event) => setForm(current => ({ ...current, qty: event.target.value }))}
                style={fieldStyle}
              />
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '4px', flex: 1 }}>
              <label style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>{t('Unit Cost')}</label>
              <input
                type="number"
                min="0"
                step="0.0001"
                value={form.unitCost}
                disabled={form.mvmtType === 'Issue'}
                onChange={(event) => setForm(current => ({ ...current, unitCost: event.target.value }))}
                style={{ ...fieldStyle, backgroundColor: form.mvmtType === 'Issue' ? '#f3f4f6' : '#fff' }}
              />
            </div>
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
            <label style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>{t('Cost Center')}</label>
            <input
              type="text"
              value={form.costCenter}
              onChange={(event) => setForm(current => ({ ...current, costCenter: event.target.value }))}
              style={fieldStyle}
            />
          </div>

          <button type="submit" className="btn btn-primary" disabled={!canSubmit || loading} style={{ marginTop: '8px', padding: '10px', justifyContent: 'center', opacity: !canSubmit || loading ? 0.7 : 1 }}>
            {loading ? <RefreshCw size={16} /> : <ArrowDownToLine size={16} />}
            {loading ? t('Processing') : t('Submit Movement')}
          </button>
        </form>

        <div className="panel p-0 hover-lift" style={{ display: 'flex', flexDirection: 'column' }}>
          <div style={{ padding: '16px', borderBottom: '1px solid var(--border-color)', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <History size={16} style={{ color: 'var(--text-secondary)' }} />
            <h2 style={{ fontSize: '14px', fontWeight: 600 }}>{t('Recent Movements')}</h2>
            <button type="button" onClick={loadMovements} className="btn btn-secondary" style={{ marginLeft: 'auto', padding: '6px 10px' }}>
              <RefreshCw size={14} />
              {t('Refresh')}
            </button>
          </div>
          <div className="data-grid-container" style={{ border: 'none', borderRadius: '0 0 6px 6px' }}>
            <table className="data-grid">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>{t('Movement Type')}</th>
                  <th>{t('Material Code')}</th>
                  <th style={{ textAlign: 'right' }}>{t('Quantity')}</th>
                  <th>{t('Created Date')}</th>
                </tr>
              </thead>
              <tbody>
                {historyLoading ? (
                  <tr>
                    <td colSpan={5} style={{ color: 'var(--text-secondary)' }}>{t('Loading movements...')}</td>
                  </tr>
                ) : movements.length === 0 ? (
                  <tr>
                    <td colSpan={5} style={{ color: 'var(--text-secondary)' }}>{t('No movements found.')}</td>
                  </tr>
                ) : movements.map(m => (
                  <tr key={m.ID}>
                    <td style={{ fontWeight: 500 }}>{m.MvmtNo}</td>
                    <td>
                      <span className={`status-badge ${m.MvmtType.toLowerCase()}`}>
                        {t(m.MvmtType)}
                      </span>
                    </td>
                    <td>{m.ItemCode}</td>
                    <td style={{ textAlign: 'right', fontWeight: 500, color: m.MvmtType === 'Receipt' ? 'var(--status-success)' : 'var(--status-critical)' }}>
                      {m.MvmtType === 'Receipt' ? `+${m.Qty}` : `-${m.Qty}`}
                    </td>
                    <td style={{ color: 'var(--text-secondary)' }}>{formatDate(m.CreatedAt)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};

const fieldStyle: React.CSSProperties = {
  padding: '8px',
  border: '1px solid var(--border-color)',
  borderRadius: '4px',
  fontSize: 'var(--font-sm)',
  outline: 'none',
};

const receiptMovementType: MovementType = 'Receipt';
const issueMovementType: MovementType = 'Issue';

const formatDate = (value: string) => {
  if (!value) {
    return '';
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString();
};

const extractErrorMessage = (err: unknown) => {
  if (typeof err === 'object' && err !== null && 'response' in err) {
    const response = (err as { response?: { data?: { error?: string; msg?: string } } }).response;
    return response?.data?.error || response?.data?.msg || 'Submit failed';
  }
  return err instanceof Error ? err.message : 'Submit failed';
};
