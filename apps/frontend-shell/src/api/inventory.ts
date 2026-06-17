import axios from 'axios';

export type MovementType = 'Receipt' | 'Issue';

export interface GoodsMovementRecord {
  ID: number;
  CreatedAt: string;
  MvmtNo: string;
  MvmtType: MovementType;
  ItemCode: string;
  Qty: number;
}

export interface JournalEntryLine {
  AccountCode: string;
  DebitAmount: number;
  CreditAmount: number;
}

export interface JournalEntryRecord {
  ID: number;
  EntryNo: string;
  TotalAmount: number;
  Lines: JournalEntryLine[];
}

export interface GoodsMovementRequest {
  mvmtNo?: string;
  mvmtType: MovementType;
  itemCode: string;
  qty: number;
  unitCost?: number;
  costCenter?: string;
}

export interface GoodsMovementPostResult {
  movement: GoodsMovementRecord;
  journalEntry: JournalEntryRecord;
}

export const fetchGoodsMovements = async (): Promise<GoodsMovementRecord[]> => {
  const response = await axios.get('/openerp/v1/goods-movement');
  return response.data.data;
};

export const submitGoodsMovement = async (
  data: GoodsMovementRequest,
): Promise<GoodsMovementPostResult> => {
  const response = await axios.post('/openerp/v1/goods-movement', data);
  return response.data.data;
};
