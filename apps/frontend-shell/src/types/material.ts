export interface Material {
  id: string;
  code: string;
  name: string;
  category: string;
  unit: string;
  stock: number;
  price: number;
  status: 'active' | 'inactive';
  lastUpdated: string;
}
