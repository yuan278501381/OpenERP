import { useState, useEffect } from 'react';
import axios from 'axios';
import type { Material } from '../types/material';
import { logger } from '../utils/logger';

export const useMaterials = () => {
  const [materials, setMaterials] = useState<Material[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchMaterials = async () => {
      try {
        setLoading(true);
        // Attempt to fetch from real endpoint
        const response = await axios.get('/openerp/v1/materials');
        if (response.data && response.data.data) {
          setMaterials(response.data.data);
          logger.info('Successfully fetched materials from API', { count: response.data.data.length });
        } else {
          setMaterials(response.data);
          logger.info('Successfully fetched materials from API', { count: response.data.length });
        }
      } catch (err: unknown) {
        // Fallback to mock data if endpoint is not available yet
        logger.warn('Failed to fetch from /openerp/v1/materials, using mock data.', { error: err instanceof Error ? err.message : String(err) });
        
        // Mock data to ensure the UI can be tested
        const mockData: Material[] = [
          { id: '1', code: 'MAT-001', name: 'Aluminum Sheet 2mm', category: 'Raw Materials', unit: 'kg', stock: 1250, price: 4.5, status: 'active', lastUpdated: '2026-06-15' },
          { id: '2', code: 'MAT-002', name: 'Steel Beam H-Type', category: 'Raw Materials', unit: 'pcs', stock: 340, price: 120.0, status: 'active', lastUpdated: '2026-06-16' },
          { id: '3', code: 'MAT-003', name: 'Industrial Lubricant XL', category: 'Consumables', unit: 'L', stock: 85, price: 15.2, status: 'inactive', lastUpdated: '2026-06-10' },
          { id: '4', code: 'MAT-004', name: 'Copper Wire 1.5mm', category: 'Raw Materials', unit: 'm', stock: 5000, price: 2.1, status: 'active', lastUpdated: '2026-06-14' },
          { id: '5', code: 'MAT-005', name: 'Packaging Box Standard', category: 'Packaging', unit: 'pcs', stock: 12000, price: 0.5, status: 'active', lastUpdated: '2026-06-16' },
        ];
        setMaterials(mockData);
        // Optional: you can set error state if you want to show error UI
        setError(err instanceof Error ? err.message : String(err));
      } finally {
        setLoading(false);
      }
    };

    fetchMaterials();
  }, []);

  return { materials, loading, error };
};
