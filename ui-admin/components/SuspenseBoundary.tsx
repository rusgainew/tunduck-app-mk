/**
 * Lazy Loading Suspense Boundaries
 * Provides consistent loading UI for code-split components
 */

import React from 'react';
import { LoadingState } from '@/components/states/LoadingState';

/**
 * Suspense fallback for page components
 * Shows while lazy-loaded pages are loading
 */
export const PageLoadingFallback = () => (
  <div className="min-h-screen flex items-center justify-center bg-gray-50">
    <LoadingState 
      message="Загрузка страницы..." 
      fullHeight={false}
    />
  </div>
);

/**
 * Suspense boundary wrapper for pages
 * Usage:
 * <SuspenseBoundary>
 *   <LazyComponent />
 * </SuspenseBoundary>
 */
export const SuspenseBoundary: React.FC<{ children: React.ReactNode }> = ({ 
  children 
}) => (
  <React.Suspense fallback={<PageLoadingFallback />}>
    {children}
  </React.Suspense>
);

/**
 * Suspense boundary with custom fallback
 */
export const CustomSuspenseBoundary: React.FC<{
  children: React.ReactNode;
  fallback?: React.ReactNode;
}> = ({ children, fallback = <PageLoadingFallback /> }) => (
  <React.Suspense fallback={fallback}>
    {children}
  </React.Suspense>
);
