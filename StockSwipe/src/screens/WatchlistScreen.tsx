import React, { useEffect } from 'react';
import { View, Text, StyleSheet, FlatList, ActivityIndicator } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { StockListItem } from '../components/StockListItem';
import { useStore } from '../store';
import { StockDetail } from '../types';

export const WatchlistScreen: React.FC = () => {
  const navigation = useNavigation<any>();
  const { watchlist, isLoading, fetchWatchlist, removeFromWatchlist } = useStore();

  useEffect(() => {
    fetchWatchlist();
  }, []);

  const handleStockPress = (stock: StockDetail) => {
    navigation.navigate('StockDetail', { ticker: stock.ticker });
  };

  const renderItem = ({ item }: { item: StockDetail }) => (
    <StockListItem 
      stock={item} 
      onPress={() => handleStockPress(item)}
      onRemove={() => removeFromWatchlist(item.id)}
    />
  );

  if (isLoading && watchlist.length === 0) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#6C63FF" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Watchlist</Text>
      <Text style={styles.subtitle}>{watchlist.length} stocks</Text>
      
      {watchlist.length === 0 ? (
        <View style={styles.empty}>
          <Text style={styles.emptyEmoji}>📝</Text>
          <Text style={styles.emptyText}>No stocks in your watchlist yet</Text>
          <Text style={styles.emptySubtext}>Swipe right on stocks to add them here</Text>
        </View>
      ) : (
        <FlatList
          data={watchlist}
          renderItem={renderItem}
          keyExtractor={(item) => item.id}
          contentContainerStyle={styles.list}
          showsVerticalScrollIndicator={false}
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0f0f1e',
    paddingTop: 60,
    paddingHorizontal: 20,
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#fff',
  },
  subtitle: {
    fontSize: 16,
    color: '#888',
    marginTop: 4,
    marginBottom: 20,
  },
  list: {
    paddingBottom: 20,
  },
  empty: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  emptyEmoji: {
    fontSize: 64,
    marginBottom: 16,
  },
  emptyText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: '600',
  },
  emptySubtext: {
    color: '#888',
    marginTop: 8,
  },
});
