import { create } from 'zustand';
import { StockDetail } from '../types';
import { api } from '../api/client';

interface AppState {
  userId: string;
  currentStock: StockDetail | null;
  watchlist: StockDetail[];
  isLoading: boolean;
  error: string | null;

  fetchDiscoverStock: () => Promise<void>;
  swipeStock: (direction: 'left' | 'right') => Promise<void>;
  fetchWatchlist: () => Promise<void>;
  removeFromWatchlist: (stockId: string) => Promise<void>;
  clearError: () => void;
}

export const useStore = create<AppState>((set, get) => ({
  userId: api.getUserId(),
  currentStock: null,
  watchlist: [],
  isLoading: false,
  error: null,

  fetchDiscoverStock: async () => {
    set({ isLoading: true, error: null });
    try {
      const stock = await api.getDiscoverStock();
      set({ currentStock: stock, isLoading: false });
    } catch (err: any) {
      if (err.response?.status === 204) {
        set({ currentStock: null, isLoading: false, error: 'No more stocks to discover!' });
      } else {
        set({ isLoading: false, error: 'Failed to load stock' });
      }
    }
  },

  swipeStock: async (direction: 'left' | 'right') => {
    const { currentStock } = get();
    if (!currentStock) return;

    set({ isLoading: true, error: null });
    try {
      await api.swipeStock(currentStock.id, direction);
      
      if (direction === 'right') {
        const { watchlist } = get();
        set({ watchlist: [...watchlist, currentStock] });
      }
      
      set({ currentStock: null, isLoading: false });
      await get().fetchDiscoverStock();
    } catch (err) {
      set({ isLoading: false, error: 'Failed to swipe stock' });
    }
  },

  fetchWatchlist: async () => {
    set({ isLoading: true, error: null });
    try {
      const watchlist = await api.getWatchlist();
      set({ watchlist, isLoading: false });
    } catch (err) {
      set({ isLoading: false, error: 'Failed to load watchlist' });
    }
  },

  removeFromWatchlist: async (stockId: string) => {
    try {
      await api.removeFromWatchlist(stockId);
      const { watchlist } = get();
      set({ watchlist: watchlist.filter(s => s.id !== stockId) });
    } catch (err) {
      set({ error: 'Failed to remove from watchlist' });
    }
  },

  clearError: () => set({ error: null }),
}));
