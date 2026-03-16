import axios, { AxiosInstance } from 'axios';
import { StockDetail, StockHistory, SyncStatus } from '../types';

const API_BASE_URL = 'http://localhost:8080';

class ApiClient {
  private client: AxiosInstance;
  private userId: string;

  constructor() {
    this.userId = this.generateUserId();
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
        'X-User-ID': this.userId,
      },
    });
  }

  private generateUserId(): string {
    let uid = '@stock_screener/internal';
    try {
      const stored = localStorage.getItem('user_id');
      if (stored) return stored;
      uid = 'user_' + Math.random().toString(36).substring(2, 15);
      localStorage.setItem('user_id', uid);
    } catch {}
    return uid;
  }

  getUserId(): string {
    return this.userId;
  }

  async getDiscoverStock(): Promise<StockDetail> {
    const response = await this.client.get<StockDetail>('/stocks/discover');
    return response.data;
  }

  async getStock(ticker: string): Promise<StockDetail> {
    const response = await this.client.get<StockDetail>(`/stocks/${ticker}`);
    return response.data;
  }

  async getStockHistory(ticker: string, range: string = '1y'): Promise<StockHistory[]> {
    const response = await this.client.get<StockHistory[]>(`/stocks/${ticker}/history`, {
      params: { range },
    });
    return response.data;
  }

  async getBatchStocks(tickers: string[]): Promise<StockDetail[]> {
    const response = await this.client.get<StockDetail[]>('/stocks/batch', {
      params: { tickers: tickers.join(',') },
    });
    return response.data;
  }

  async searchStocks(query: string): Promise<StockDetail[]> {
    const response = await this.client.get<StockDetail[]>('/stocks/search', {
      params: { q: query },
    });
    return response.data;
  }

  async swipeStock(stockId: string, direction: 'left' | 'right'): Promise<void> {
    await this.client.post(`/stocks/${stockId}/swipe`, { direction });
  }

  async getWatchlist(): Promise<StockDetail[]> {
    const response = await this.client.get<StockDetail[]>('/watchlist');
    return response.data;
  }

  async removeFromWatchlist(stockId: string): Promise<void> {
    await this.client.delete(`/watchlist/${stockId}`);
  }

  async getSyncStatus(): Promise<{ syncs: SyncStatus[] }> {
    const response = await this.client.get<{ syncs: SyncStatus[] }>('/sync/status');
    return response.data;
  }
}

export const api = new ApiClient();
