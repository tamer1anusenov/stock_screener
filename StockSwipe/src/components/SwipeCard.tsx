import React from 'react';
import { View, Text, StyleSheet, Dimensions, Image, TouchableOpacity } from 'react-native';
import Animated, {
  useAnimatedStyle,
  useSharedValue,
  withSpring,
  runOnJS,
} from 'react-native-reanimated';
import { Gesture, GestureDetector } from 'react-native-gesture-handler';
import { StockDetail } from '../types';

const { width: SCREEN_WIDTH } = Dimensions.get('window');
const SWIPE_THRESHOLD = SCREEN_WIDTH * 0.25;

interface SwipeCardProps {
  stock: StockDetail;
  onSwipeLeft: () => void;
  onSwipeRight: () => void;
}

export const SwipeCard: React.FC<SwipeCardProps> = ({ stock, onSwipeLeft, onSwipeRight }) => {
  const translateX = useSharedValue(0);
  const translateY = useSharedValue(0);

  const formatMarketCap = (value: number): string => {
    if (!value) return '-';
    if (value >= 1e12) return `$${(value / 1e12).toFixed(1)}T`;
    if (value >= 1e9) return `$${(value / 1e9).toFixed(1)}B`;
    return `$${(value / 1e6).toFixed(1)}M`;
  };

  const formatPrice = (price: number): string => {
    return `$${price.toFixed(2)}`;
  };

  const handleSwipeComplete = (direction: 'left' | 'right') => {
    if (direction === 'left') {
      onSwipeLeft();
    } else {
      onSwipeRight();
    }
  };

  const panGesture = Gesture.Pan()
    .onUpdate((event) => {
      translateX.value = event.translationX;
      translateY.value = event.translationY * 0.5;
    })
    .onEnd((event) => {
      if (Math.abs(event.translationX) > SWIPE_THRESHOLD) {
        const direction = event.translationX > 0 ? 'right' : 'left';
        translateX.value = withSpring(direction === 'right' ? SCREEN_WIDTH * 1.5 : -SCREEN_WIDTH * 1.5);
        translateY.value = withSpring(0);
        runOnJS(handleSwipeComplete)(direction);
      } else {
        translateX.value = withSpring(0);
        translateY.value = withSpring(0);
      }
    });

  const animatedStyle = useAnimatedStyle(() => ({
    transform: [
      { translateX: translateX.value },
      { translateY: translateY.value },
      { rotate: `${translateX.value * 0.02}deg` },
    ],
  }));

  const likeOpacity = useAnimatedStyle(() => ({
    opacity: Math.max(0, Math.min(1, translateX.value / 100)),
  }));

  const nopeOpacity = useAnimatedStyle(() => ({
    opacity: Math.max(0, Math.min(1, -translateX.value / 100)),
  }));

  const isPositive = stock.change_percent >= 0;

  return (
    <GestureDetector gesture={panGesture}>
      <Animated.View style={[styles.card, animatedStyle]}>
        <Animated.View style={[styles.overlay, styles.likeOverlay, likeOpacity]}>
          <Text style={styles.overlayText}>ADD</Text>
        </Animated.View>
        
        <Animated.View style={[styles.overlay, styles.nopeOverlay, nopeOpacity]}>
          <Text style={styles.overlayText}>SKIP</Text>
        </Animated.View>

        <View style={styles.header}>
          {stock.logo_url ? (
            <Image source={{ uri: stock.logo_url }} style={styles.logo} />
          ) : (
            <View style={styles.logoPlaceholder}>
              <Text style={styles.logoText}>{stock.ticker.charAt(0)}</Text>
            </View>
          )}
          <View style={styles.headerInfo}>
            <Text style={styles.ticker}>{stock.ticker}</Text>
            <Text style={styles.companyName} numberOfLines={1}>{stock.company_name}</Text>
          </View>
        </View>

        <View style={styles.priceContainer}>
          <Text style={styles.price}>{formatPrice(stock.price)}</Text>
          <View style={[styles.changeBadge, isPositive ? styles.positiveChange : styles.negativeChange]}>
            <Text style={[styles.changeText, isPositive ? styles.positiveText : styles.negativeText]}>
              {isPositive ? '+' : ''}{stock.change_percent.toFixed(2)}%
            </Text>
          </View>
        </View>

        <View style={styles.badges}>
          <View style={styles.badge}>
            <Text style={styles.badgeLabel}>Sector</Text>
            <Text style={styles.badgeValue}>{stock.sector || 'Unknown'}</Text>
          </View>
          <View style={styles.badge}>
            <Text style={styles.badgeLabel}>Market Cap</Text>
            <Text style={styles.badgeValue}>{formatMarketCap(stock.market_cap)}</Text>
          </View>
        </View>

        <View style={styles.metrics}>
          <View style={styles.metric}>
            <Text style={styles.metricLabel}>P/E Ratio</Text>
            <Text style={styles.metricValue}>{stock.pe_ratio?.toFixed(1) || '-'}</Text>
          </View>
          <View style={styles.metric}>
            <Text style={styles.metricLabel}>EPS</Text>
            <Text style={styles.metricValue}>${stock.eps?.toFixed(2) || '-'}</Text>
          </View>
          <View style={styles.metric}>
            <Text style={styles.metricLabel}>Revenue Growth</Text>
            <Text style={[styles.metricValue, stock.revenue_growth > 0 ? styles.positiveText : styles.negativeText]}>
              {stock.revenue_growth ? `${(stock.revenue_growth * 100).toFixed(1)}%` : '-'}
            </Text>
          </View>
        </View>

        <Text style={styles.swipeHint}>← Swipe to skip | Add to watchlist →</Text>
      </Animated.View>
    </GestureDetector>
  );
};

const styles = StyleSheet.create({
  card: {
    width: SCREEN_WIDTH - 40,
    backgroundColor: '#1a1a2e',
    borderRadius: 20,
    padding: 24,
    position: 'absolute',
    top: 20,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 10 },
    shadowOpacity: 0.3,
    shadowRadius: 20,
    elevation: 10,
  },
  overlay: {
    position: 'absolute',
    top: 40,
    zIndex: 10,
    padding: 10,
    borderRadius: 8,
    borderWidth: 3,
  },
  likeOverlay: {
    right: 40,
    borderColor: '#4CAF50',
    transform: [{ rotate: '15deg' }],
  },
  nopeOverlay: {
    left: 40,
    borderColor: '#F44336',
    transform: [{ rotate: '-15deg' }],
  },
  overlayText: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#fff',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 20,
  },
  logo: {
    width: 60,
    height: 60,
    borderRadius: 12,
    backgroundColor: '#2a2a4e',
  },
  logoPlaceholder: {
    width: 60,
    height: 60,
    borderRadius: 12,
    backgroundColor: '#2a2a4e',
    justifyContent: 'center',
    alignItems: 'center',
  },
  logoText: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#6C63FF',
  },
  headerInfo: {
    marginLeft: 16,
    flex: 1,
  },
  ticker: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#fff',
  },
  companyName: {
    fontSize: 16,
    color: '#888',
    marginTop: 2,
  },
  priceContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 20,
  },
  price: {
    fontSize: 42,
    fontWeight: 'bold',
    color: '#fff',
  },
  changeBadge: {
    marginLeft: 16,
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 20,
  },
  positiveChange: {
    backgroundColor: 'rgba(76, 175, 80, 0.2)',
  },
  negativeChange: {
    backgroundColor: 'rgba(244, 67, 54, 0.2)',
  },
  changeText: {
    fontSize: 16,
    fontWeight: '600',
  },
  positiveText: {
    color: '#4CAF50',
  },
  negativeText: {
    color: '#F44336',
  },
  badges: {
    flexDirection: 'row',
    marginBottom: 20,
  },
  badge: {
    flex: 1,
    backgroundColor: '#2a2a4e',
    borderRadius: 12,
    padding: 12,
    marginRight: 8,
  },
  badgeLabel: {
    fontSize: 12,
    color: '#888',
    marginBottom: 4,
  },
  badgeValue: {
    fontSize: 14,
    fontWeight: '600',
    color: '#fff',
  },
  metrics: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  metric: {
    flex: 1,
    alignItems: 'center',
  },
  metricLabel: {
    fontSize: 11,
    color: '#666',
    marginBottom: 4,
  },
  metricValue: {
    fontSize: 18,
    fontWeight: '600',
    color: '#fff',
  },
  swipeHint: {
    textAlign: 'center',
    color: '#444',
    fontSize: 12,
    marginTop: 20,
  },
});
