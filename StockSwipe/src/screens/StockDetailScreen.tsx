import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, ScrollView, Image, ActivityIndicator, TouchableOpacity } from 'react-native';
import { useRoute, useNavigation } from '@react-navigation/native';
import { api } from '../api/client';
import { StockDetail, StockHistory } from '../types';

export const StockDetailScreen: React.FC = () => {
  const route = useRoute<any>();
  const navigation = useNavigation();
  const { ticker } = route.params;

  const [stock, setStock] = useState<StockDetail | null>(null);
  const [history, setHistory] = useState<StockHistory[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [stockData, historyData] = await Promise.all([
          api.getStock(ticker),
          api.getStockHistory(ticker, '1y'),
        ]);
        setStock(stockData);
        setHistory(historyData);
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [ticker]);

  const formatMarketCap = (value: number): string => {
    if (!value) return '-';
    if (value >= 1e12) return `$${(value / 1e12).toFixed(1)}T`;
    if (value >= 1e9) return `$${(value / 1e9).toFixed(1)}B`;
    return `$${(value / 1e6).toFixed(1)}M`;
  };

  if (loading) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#6C63FF" />
      </View>
    );
  }

  if (!stock) {
    return (
      <View style={styles.container}>
        <Text style={styles.errorText}>Stock not found</Text>
      </View>
    );
  }

  const isPositive = stock.change_percent >= 0;
  const latestPrice = history.length > 0 ? history[history.length - 1].close : stock.price;
  const earliestPrice = history.length > 0 ? history[0].close : stock.price;
  const yearlyChange = earliestPrice > 0 ? ((latestPrice - earliestPrice) / earliestPrice) * 100 : 0;

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      <TouchableOpacity style={styles.backButton} onPress={() => navigation.goBack()}>
        <Text style={styles.backText}>← Back</Text>
      </TouchableOpacity>

      <View style={styles.header}>
        {stock.logo_url ? (
          <Image source={{ uri: stock.logo_url }} style={styles.logo} />
        ) : (
          <View style={styles.logoPlaceholder}>
            <Text style={styles.logoText}>{stock.ticker.charAt(0)}</Text>
          </View>
        )}
        <Text style={styles.ticker}>{stock.ticker}</Text>
        <Text style={styles.companyName}>{stock.company_name}</Text>
      </View>

      <View style={styles.priceSection}>
        <Text style={styles.price}>${stock.price.toFixed(2)}</Text>
        <View style={[styles.changeBadge, isPositive ? styles.positiveChange : styles.negativeChange]}>
          <Text style={[styles.changeText, isPositive ? styles.positiveText : styles.negativeText]}>
            {isPositive ? '+' : ''}{stock.change_percent.toFixed(2)}%
          </Text>
        </View>
        <Text style={[styles.yearlyChange, yearlyChange >= 0 ? styles.positiveText : styles.negativeText]}>
          {yearlyChange >= 0 ? '+' : ''}{yearlyChange.toFixed(1)}% (1Y)
        </Text>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Overview</Text>
        <View style={styles.grid}>
          <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>Sector</Text>
            <Text style={styles.gridValue}>{stock.sector || 'Unknown'}</Text>
          </View>
          <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>Industry</Text>
            <Text style={styles.gridValue}>{stock.industry || 'Unknown'}</Text>
          </View>
          <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>Market Cap</Text>
            <Text style={styles.gridValue}>{formatMarketCap(stock.market_cap)}</Text>
          </View>
          <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>P/E Ratio</Text>
            <Text style={styles.gridValue}>{stock.pe_ratio?.toFixed(1) || '-'}</Text>
          </View>
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Financials</Text>
        <View style={styles.grid}>
          <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>EPS</Text>
            <Text style={styles.gridValue}>${stock.eps?.toFixed(2) || '-'}</Text>
          </View>
          <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>Revenue Growth</Text>
            <Text style={[styles.gridValue, stock.revenue_growth > 0 ? styles.positiveText : styles.negativeText]}>
              {stock.revenue_growth ? `${(stock.revenue_growth * 100).toFixed(1)}%` : '-'}
            </Text>
          </View>
        </View>
      </View>

      {stock.description && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>About</Text>
          <Text style={styles.description}>{stock.description}</Text>
        </View>
      )}

      <View style={{ height: 40 }} />
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0f0f1e',
    paddingHorizontal: 20,
  },
  backButton: {
    marginTop: 60,
    marginBottom: 20,
  },
  backText: {
    color: '#6C63FF',
    fontSize: 16,
  },
  header: {
    alignItems: 'center',
    marginBottom: 20,
  },
  logo: {
    width: 80,
    height: 80,
    borderRadius: 16,
    backgroundColor: '#2a2a4e',
  },
  logoPlaceholder: {
    width: 80,
    height: 80,
    borderRadius: 16,
    backgroundColor: '#2a2a4e',
    justifyContent: 'center',
    alignItems: 'center',
  },
  logoText: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#6C63FF',
  },
  ticker: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#fff',
    marginTop: 12,
  },
  companyName: {
    fontSize: 16,
    color: '#888',
    marginTop: 4,
  },
  priceSection: {
    alignItems: 'center',
    marginBottom: 30,
  },
  price: {
    fontSize: 48,
    fontWeight: 'bold',
    color: '#fff',
  },
  changeBadge: {
    marginTop: 8,
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 20,
  },
  positiveChange: {
    backgroundColor: 'rgba(76, 175, 80, 0.2)',
  },
  negativeChange: {
    backgroundColor: 'rgba(244, 67, 54, 0.2)',
  },
  changeText: {
    fontSize: 18,
    fontWeight: '600',
  },
  positiveText: {
    color: '#4CAF50',
  },
  negativeText: {
    color: '#F44336',
  },
  yearlyChange: {
    marginTop: 8,
    fontSize: 14,
  },
  section: {
    marginBottom: 24,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: '600',
    color: '#fff',
    marginBottom: 12,
  },
  grid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
  },
  gridItem: {
    width: '50%',
    marginBottom: 16,
  },
  gridLabel: {
    fontSize: 12,
    color: '#888',
    marginBottom: 4,
  },
  gridValue: {
    fontSize: 16,
    color: '#fff',
    fontWeight: '500',
  },
  description: {
    fontSize: 14,
    color: '#aaa',
    lineHeight: 22,
  },
  errorText: {
    color: '#fff',
    fontSize: 18,
    textAlign: 'center',
    marginTop: 100,
  },
});
