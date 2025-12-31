/**
 * Advanced Query Hooks with Cache Strategy Integration
 * Uses cache-strategies.ts for optimal caching patterns
 */

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import {
  getStaleWhileRevalidateOptions,
  getAggressiveCacheOptions,
  getNormalCacheOptions,
  getFreshDataOptions,
  getCacheFirstOptions,
} from '@/lib/cache-strategies';

/**
 * Hook for data that changes occasionally (stale-while-revalidate pattern)
 * Useful for: User profiles, Organization lists, Document metadata
 * Updates: Every 5 minutes, but serves stale data while fetching
 */
export function useOptimizedQuery<TData>(
  options: UseQueryOptions<TData>,
  strategy: 'stale-while-revalidate' | 'aggressive' | 'normal' | 'fresh' | 'cache-first' = 'stale-while-revalidate'
) {
  const strategyOptions = {
    'stale-while-revalidate': getStaleWhileRevalidateOptions(),
    'aggressive': getAggressiveCacheOptions(),
    'normal': getNormalCacheOptions(),
    'fresh': getFreshDataOptions(),
    'cache-first': getCacheFirstOptions(),
  }[strategy];

  return useQuery({
    ...options,
    ...strategyOptions,
  } as UseQueryOptions<TData>);
}

/**
 * Hook for rarely-changing data (organizations, roles)
 * Cache for 30 minutes, very low refetch frequency
 */
export function useAggressiveCachedQuery<TData>(
  options: UseQueryOptions<TData>
) {
  return useQuery({
    ...options,
    ...getAggressiveCacheOptions(),
  } as UseQueryOptions<TData>);
}

/**
 * Hook for frequently-needed fresh data (current user, real-time counts)
 * Cache for 1 minute, prioritizes freshness
 */
export function useFreshDataQuery<TData>(
  options: UseQueryOptions<TData>
) {
  return useQuery({
    ...options,
    ...getFreshDataOptions(),
  } as UseQueryOptions<TData>);
}

/**
 * Hook for truly static data (never changes after initial load)
 * Cache indefinitely, only updates on manual invalidation
 */
export function useCacheFirstQuery<TData>(
  options: UseQueryOptions<TData>
) {
  return useQuery({
    ...options,
    ...getCacheFirstOptions(),
  } as UseQueryOptions<TData>);
}
