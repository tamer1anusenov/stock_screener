import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image } from 'react-native';
import { StockDetail } from '../types';

interface StockListItemProps {
  stock: StockDetail;
  onPress: () => void;
  onRemove?: () => void;
}

export const StockListItem: React.FC<StockListItemProps> = ({ stock, onPress, onRemove }) => {
  const formatMarketCap = (value: number): string => {
    if (!value) return '-';
    if (value >= 1e12) return `$${(value / 1e12).toFixed(1)}T`;
    if (value >= 1e9) return `$${(value / 1e9).toFixed(1)}B`;
    return `$${(value / 1e6).toFixed(1)}M`;
  };

  const isPositive = stock.change_percent >= 0;

  return (
    <TouchableOpacity style={styles.container} onPress={onPress}>
      <View style={styles.left}>
        {stock.logo_url ? (
          <Image source={{ uri: stock.logo_url }} style={styles.logo} />
        ) : (
          <View style={styles.logoPlaceholder}>
            <Text style={styles.logoText}>{stock.ticker.charAt(0)}</Text>
          </View>
        )}
        <View style={styles.info}>
          <Text style={styles.ticker}>{stock.ticker}</Text>
          <Text style={styles.companyName} numberOfLines={1}>{stock.company_name}</Text>
        </View>
      </View>
      <View style={styles.right}>
        <Text style={styles.price}>${stock.price.toFixed(2)}</Text>
        <Text style={[styles.change, isPositive ? styles.positive : styles.negative]}>
          {isPositive ? '+' : ''}{stock.change_percent.toFixed(2)}%
        </Text>
        {onRemove && (
          <TouchableOpacity onPress={onRemove} style={styles.removeBtn}>
            <Text style={styles.removeText}>✕</Text>
          </TouchableOpacity>
        )}
      </View>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    backgroundColor: '#1a1a2e',
    padding: 16,
    borderRadius: 12,
    marginBottom: 8,
  },
  left: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  logo: {
    width: 44,
    height: 44,
    borderRadius: 10,
    backgroundColor: '#2a2a4e',
  },
  logoPlaceholder: {
    width: 44,
    height: 44,
    borderRadius: 10,
    backgroundColor: '#2a2a4e',
    justifyContent: 'center',
    alignItems: 'center',
  },
  logoText: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#6C63FF',
  },
  info: {
    marginLeft: 12,
    flex: 1,
  },
  ticker: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#fff',
  },
  companyName: {
    fontSize: 13,
    color: '#888',
    marginTop: 2,
  },
  right: {
    alignItems: 'flex-end',
  },
  price: {
    fontSize: 18,
    fontWeight: '600',
    color: '#fff',
  },
  change: {
    fontSize: 14,
    marginTop: 2,
  },
  positive: {
    color: '#4CAF50',
  },
  negative: {
    color: '#F44336',
  },
  removeBtn: {
    marginTop: 4,
    padding: 4,
  },
  removeText: {
    color: '#666',
    fontSize: 16,
  },
});
