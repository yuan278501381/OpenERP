import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';

const generateTraceId = () => Math.random().toString(36).substring(2, 15);
const logger = {
  info: (msg: string, traceId: string) => console.log(`[INFO] [${traceId}] ${msg}`),
};

export const DocumentMatrix: React.FC = () => {
  const { t } = useTranslation();
  const traceId = generateTraceId();
  const [lines, setLines] = useState([
    { id: 1, itemCode: 'A0001', description: 'Server PCB', quantity: 10, price: 1500, total: 15000 },
    { id: 2, itemCode: 'A0002', description: 'Cooling Fan', quantity: 50, price: 20, total: 1000 }
  ]);

  const addLine = () => {
    logger.info('Added new line', traceId);
    setLines([...lines, { id: lines.length + 1, itemCode: '', description: '', quantity: 0, price: 0, total: 0 }]);
  };

  const handleInputChange = (id: number, field: string, value: string | number) => {
    setLines(lines.map(line => {
      if (line.id === id) {
        const newLine = { ...line, [field]: value };
        if (field === 'quantity' || field === 'price') {
          newLine.total = Number(newLine.quantity) * Number(newLine.price);
        }
        return newLine;
      }
      return line;
    }));
  };

  const inputStyle = {
    width: '100%',
    border: '1px solid transparent',
    padding: '4px 8px',
    outline: 'none',
    backgroundColor: 'transparent',
    fontSize: '13px',
    borderRadius: '2px',
    transition: 'border-color 0.2s',
  };

  return (
    <div className="panel" style={{ marginTop: '16px' }}>
      <div className="data-grid-container" style={{ maxHeight: '400px' }}>
        <table className="data-grid" style={{ width: '100%', minWidth: '800px' }}>
          <thead>
            <tr>
              <th style={{ width: '40px', textAlign: 'center' }}>#</th>
              <th style={{ width: '150px' }}>{t('Item Code')}</th>
              <th>{t('Item Description')}</th>
              <th style={{ width: '100px' }}>{t('Quantity')}</th>
              <th style={{ width: '120px' }}>{t('Unit Price')}</th>
              <th style={{ width: '120px' }}>{t('Total')}</th>
            </tr>
          </thead>
          <tbody>
            {lines.map((line, index) => (
              <tr key={line.id}>
                <td style={{ textAlign: 'center', color: '#6b7280', fontSize: '12px' }}>{index + 1}</td>
                <td>
                  <input 
                    type="text" 
                    value={line.itemCode} 
                    style={inputStyle}
                    onChange={(e) => handleInputChange(line.id, 'itemCode', e.target.value)}
                    onFocus={(e) => e.target.style.borderColor = '#9ca3af'}
                    onBlur={(e) => e.target.style.borderColor = 'transparent'}
                  />
                </td>
                <td>
                  <input 
                    type="text" 
                    value={line.description} 
                    style={inputStyle}
                    onChange={(e) => handleInputChange(line.id, 'description', e.target.value)}
                    onFocus={(e) => e.target.style.borderColor = '#9ca3af'}
                    onBlur={(e) => e.target.style.borderColor = 'transparent'}
                  />
                </td>
                <td>
                  <input 
                    type="number" 
                    value={line.quantity} 
                    style={inputStyle}
                    onChange={(e) => handleInputChange(line.id, 'quantity', e.target.value)}
                    onFocus={(e) => e.target.style.borderColor = '#9ca3af'}
                    onBlur={(e) => e.target.style.borderColor = 'transparent'}
                  />
                </td>
                <td>
                  <input 
                    type="number" 
                    value={line.price} 
                    style={inputStyle}
                    onChange={(e) => handleInputChange(line.id, 'price', e.target.value)}
                    onFocus={(e) => e.target.style.borderColor = '#9ca3af'}
                    onBlur={(e) => e.target.style.borderColor = 'transparent'}
                  />
                </td>
                <td>
                  <input 
                    type="number" 
                    value={line.total} 
                    style={{ ...inputStyle, backgroundColor: '#f9fafb', color: '#4b5563' }}
                    readOnly 
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div style={{ padding: '12px', borderTop: '1px solid #e5e7eb' }}>
        <button className="btn btn-secondary" onClick={addLine} style={{ fontSize: '12px' }}>+ {t('Add Row')}</button>
      </div>
    </div>
  );
};
