export interface Stock {
  id: string;
  ticker: string;
  company_name: string;
  sector: string;
  industry: string;
  market_cap: number;
  pe_ratio: number;
  eps: number;
  revenue_growth: number;
  ranking_score: number;
  created_at: string;
  updated_at: string;
  logo_url?: string;
  description?: string;
}

export interface StockDetail extends Stock {
  price: number;
  change_percent: number;
  price_updated_at: string;
}

export interface StockHistory {
  stock_id: string;
  date: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
}

export interface Swipe {
  id: string;
  user_id: string;
  stock_id: string;
  direction: 'left' | 'right';
  created_at: string;
}

export interface SyncStatus {
  sync_type: string;
  last_sync: string | null;
  status: 'idle' | 'running' | 'completed' | 'failed';
  records_synced: number;
  error_message: string | null;
}

export interface ApiError {
  error: string;
}
