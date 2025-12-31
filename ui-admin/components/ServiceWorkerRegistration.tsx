'use client';

import { useEffect } from 'react';
import { registerServiceWorker } from '@/lib/service-worker';

/**
 * Component to register Service Worker on app initialization
 * Enables offline support and asset caching
 */
export function ServiceWorkerRegistration() {
  useEffect(() => {
    // Only register in browser
    if (typeof window !== 'undefined') {
      registerServiceWorker()
        .then(() => {
          console.log('✓ Service Worker registered successfully');
        })
        .catch((error) => {
          console.warn('⚠ Service Worker registration failed:', error.message);
          // Not critical, app works without SW
        });
    }
  }, []);

  // This component doesn't render anything
  return null;
}
