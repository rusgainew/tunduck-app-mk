/**
 * Service Worker Registration
 * Enables offline support and asset caching
 * 
 * Usage: Call registerServiceWorker() in app initialization
 */

export const registerServiceWorker = async () => {
  if (typeof window === 'undefined') {
    console.warn('Service Worker can only be registered on client side');
    return;
  }

  if (!('serviceWorker' in navigator)) {
    console.warn('Service Workers are not supported in this browser');
    return;
  }

  try {
    const registration = await navigator.serviceWorker.register(
      '/service-worker.js',
      { scope: '/' }
    );

    console.log('Service Worker registered successfully:', registration);

    // Handle Service Worker updates
    registration.addEventListener('updatefound', () => {
      const newWorker = registration.installing;
      if (!newWorker) return;

      newWorker.addEventListener('statechange', () => {
        if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
          // New SW available, prompt user to reload
          console.log('New Service Worker available. Please reload.');
          // Optional: Show toast notification to user
        }
      });
    });

    return registration;
  } catch (error) {
    console.error('Service Worker registration failed:', error);
  }
};

/**
 * Unregister Service Worker
 * Use for cleanup or debugging
 */
export const unregisterServiceWorker = async () => {
  if (typeof window === 'undefined') return;

  try {
    const registrations = await navigator.serviceWorker.getRegistrations();
    for (const registration of registrations) {
      await registration.unregister();
    }
    console.log('Service Worker unregistered');
  } catch (error) {
    console.error('Failed to unregister Service Worker:', error);
  }
};

/**
 * Check if Service Worker is active
 */
export const isServiceWorkerActive = async (): Promise<boolean> => {
  if (typeof window === 'undefined') return false;

  try {
    const registration = await navigator.serviceWorker.ready;
    return !!registration.active;
  } catch {
    return false;
  }
};

/**
 * Skip waiting (force update) of Service Worker
 * Use after new SW is installed
 */
export const skipServiceWorkerWaiting = () => {
  if (typeof window === 'undefined') return;

  navigator.serviceWorker.addEventListener('controllerchange', () => {
    window.location.reload();
  });

  if (navigator.serviceWorker.controller) {
    navigator.serviceWorker.controller.postMessage({ type: 'SKIP_WAITING' });
  }
};
