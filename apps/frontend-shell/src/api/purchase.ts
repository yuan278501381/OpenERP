import axios from 'axios';

export interface PurchaseOrderRequest {
  orderNo: string;
  title: string;
  amount: number;
  extData?: Record<string, unknown>;
}

export const submitPurchaseOrder = async (data: PurchaseOrderRequest) => {
  const response = await axios.post('/openerp/v1/purchase-orders', data);
  return response.data;
};
