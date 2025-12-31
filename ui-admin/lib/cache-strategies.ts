/**
 * Stale While Revalidate Pattern for TanStack Query
 * 
 * Pattern: Show stale data immediately while fetching fresh data in background
 * Benefit: Faster perceived performance, always shows data (even if stale)
 */

import { UseQueryOptions } from '@tanstack/react-query';

/**
 * Creates a query option with Stale While Revalidate pattern
 * 
 * @param staleTime - Time before data is considered stale (default: 5 minutes)
 * @param gcTime - Time to keep data in cache (default: 10 minutes)
 * 
 * @example
 * const { data } = useQuery({
 *   queryKey: ['documents'],
 *   queryFn: () => documentsApi.getAll(),
 *   ...getStaleWhileRevalidateOptions(),
 * });
 */
export const getStaleWhileRevalidateOptions = (
  staleTime: number = 5 * 60 * 1000, // 5 minutes
  gcTime: number = 10 * 60 * 1000    // 10 minutes
): Partial<UseQueryOptions> => ({
  staleTime,
  gcTime,
  refetchOnWindowFocus: false,
  refetchOnReconnect: true,  // Refetch when reconnecting
  refetchOnMount: false,     // Don't refetch on mount if fresh
});

/**
 * Aggressive cache pattern - keeps data longest
 * Use for rarely-changing data (orgs, users list)
 */
export const getAggressiveCacheOptions = (): Partial<UseQueryOptions> => ({
  staleTime: 30 * 60 * 1000,  // 30 minutes
  gcTime: 60 * 60 * 1000,     // 1 hour
  refetchOnWindowFocus: false,
  refetchOnReconnect: false,
  refetchOnMount: false,
});

/**
 * Normal cache pattern - balanced freshness and performance
 * Use for changing data (documents)
 */
export const getNormalCacheOptions = (): Partial<UseQueryOptions> => ({
  staleTime: 5 * 60 * 1000,    // 5 minutes
  gcTime: 10 * 60 * 1000,      // 10 minutes
  refetchOnWindowFocus: false,
  refetchOnReconnect: true,
  refetchOnMount: false,
});

/**
 * Fresh data pattern - minimal caching
 * Use for critical real-time data (user profile, settings)
 */
export const getFreshDataOptions = (): Partial<UseQueryOptions> => ({
  staleTime: 1 * 60 * 1000,    // 1 minute
  gcTime: 5 * 60 * 1000,       // 5 minutes
  refetchOnWindowFocus: true,
  refetchOnReconnect: true,
  refetchOnMount: true,
});

/**
 * Cache-first pattern - never refetch unless explicitly called
 * Use for static data (settings, config)
 */
export const getCacheFirstOptions = (): Partial<UseQueryOptions> => ({
  staleTime: Infinity,         // Never stale
  gcTime: 24 * 60 * 60 * 1000, // 24 hours
  refetchOnWindowFocus: false,
  refetchOnReconnect: false,
  refetchOnMount: false,
});

/**
 * Custom cache pattern
 */
export const getCustomCacheOptions = (config: {
  staleTimeMinutes?: number;
  cacheTimeMinutes?: number;
  refetchOnFocus?: boolean;
  refetchOnReconnect?: boolean;
}): Partial<UseQueryOptions> => ({
  staleTime: (config.staleTimeMinutes ?? 5) * 60 * 1000,
  gcTime: (config.cacheTimeMinutes ?? 10) * 60 * 1000,
  refetchOnWindowFocus: config.refetchOnFocus ?? false,
  refetchOnReconnect: config.refetchOnReconnect ?? true,
  refetchOnMount: false,
});
