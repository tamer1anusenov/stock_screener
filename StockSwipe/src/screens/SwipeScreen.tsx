import React, { useEffect } from 'react';
import { View, Text, StyleSheet, ActivityIndicator, TouchableOpacity } from 'react-native';
import { SwipeCard } from '../components/SwipeCard';
import { useStore } from '../store';

export const SwipeScreen: React.FC = () => {
  const { currentStock, isLoading, error, fetchDiscoverStock, swipeStock, clearError } = useStore();

  useEffect(() => {
    fetchDiscoverStock();
  }, []);

  const handleSwipeLeft = () => {
    swipeStock('left');
  };

  const handleSwipeRight = () => {
    swipeStock('right');
  };

  if (isLoading && !currentStock) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#6C63FF" />
        <Text style={styles.loadingText}>Finding stocks for you...</Text>
      </View>
    );
  }

  if (error && !currentStock) {
    return (
      <View style={styles.container}>
        <Text style={styles.errorEmoji}>🎉</Text>
        <Text style={styles.errorText}>{error}</Text>
        <TouchableOpacity style={styles.button} onPress={() => { clearError(); fetchDiscoverStock(); }}>
          <Text style={styles.buttonText}>Refresh</Text>
        </TouchableOpacity>
      </View>
    );
  }

  if (!currentStock) {
    return (
      <View style={styles.container}>
        <Text style={styles.errorEmoji}>🏁</Text>
        <Text style={styles.errorText}>You've seen all available stocks!</Text>
        <Text style={styles.subText}>Check back later for more</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Discover</Text>
      <Text style={styles.subtitle}>Swipe to find your next investment</Text>
      
      <View style={styles.cardContainer}>
        <SwipeCard 
          stock={currentStock} 
          onSwipeLeft={handleSwipeLeft}
          onSwipeRight={handleSwipeRight}
        />
      </View>

      <View style={styles.buttons}>
        <TouchableOpacity 
          style={[styles.actionButton, styles.skipButton]} 
          onPress={handleSwipeLeft}
          disabled={isLoading}
        >
          <Text style={styles.buttonEmoji}>✕</Text>
          <Text style={styles.buttonLabel}>Skip</Text>
        </TouchableOpacity>
        
        <TouchableOpacity 
          style={[styles.actionButton, styles.addButton]} 
          onPress={handleSwipeRight}
          disabled={isLoading}
        >
          <Text style={styles.buttonEmoji}>♥</Text>
          <Text style={styles.buttonLabel}>Add</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0f0f1e',
    alignItems: 'center',
    paddingTop: 60,
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
  cardContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  buttons: {
    flexDirection: 'row',
    justifyContent: 'center',
    paddingBottom: 40,
    gap: 40,
  },
  actionButton: {
    width: 70,
    height: 70,
    borderRadius: 35,
    justifyContent: 'center',
    alignItems: 'center',
  },
  skipButton: {
    backgroundColor: 'rgba(244, 67, 54, 0.2)',
    borderWidth: 2,
    borderColor: '#F44336',
  },
  addButton: {
    backgroundColor: 'rgba(76, 175, 80, 0.2)',
    borderWidth: 2,
    borderColor: '#4CAF50',
  },
  buttonEmoji: {
    fontSize: 28,
  },
  buttonLabel: {
    fontSize: 12,
    color: '#fff',
    marginTop: 2,
  },
  loadingText: {
    color: '#888',
    marginTop: 16,
    fontSize: 16,
  },
  errorEmoji: {
    fontSize: 64,
    marginBottom: 16,
  },
  errorText: {
    color: '#fff',
    fontSize: 20,
    fontWeight: '600',
    textAlign: 'center',
  },
  subText: {
    color: '#888',
    marginTop: 8,
  },
  button: {
    marginTop: 20,
    backgroundColor: '#6C63FF',
    paddingHorizontal: 32,
    paddingVertical: 12,
    borderRadius: 25,
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
