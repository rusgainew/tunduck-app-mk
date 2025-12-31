"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useEffect, useState, type ReactNode } from "react";
import { registerServiceWorker } from "@/lib/service-worker";

/**
 * Service Worker initialization component
 * Registers SW for offline support and asset caching
 */
function ServiceWorkerInitializer() {
  useEffect(() => {
    if (typeof window !== "undefined") {
      registerServiceWorker()
        .then(() => {
          console.log("✓ Service Worker registered successfully");
        })
        .catch((error) => {
          console.warn("⚠ Service Worker registration failed:", error?.message);
        });
    }
  }, []);

  return null;
}

export function Providers({ children }: { children: ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 5 * 60 * 1000, // 5 minutes - data stays fresh longer
            gcTime: 10 * 60 * 1000, // 10 minutes - keep in cache longer
            refetchOnWindowFocus: false, // Don't refetch on window focus
            retry: 1, // Retry once on failure
            retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000), // Exponential backoff
          },
          mutations: {
            retry: 1,
            retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000),
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
      <ServiceWorkerInitializer />
      {children}
    </QueryClientProvider>
  );
}
