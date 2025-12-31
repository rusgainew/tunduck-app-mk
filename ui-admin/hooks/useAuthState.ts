"use client";

import { useEffect, useCallback, useState } from "react";
import { useRouter } from "next/navigation";
import { apiClient } from "@/lib/api-client";

interface UseAuthStateReturn {
  isAuthenticated: boolean;
  isLoading: boolean;
  token: string | null;
}

/**
 * Unified hook for auth state management
 * Single source of truth for authentication status
 * Prevents duplicate auth checks across pages
 */
export function useAuthState(): UseAuthStateReturn {
  const router = useRouter();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [token, setToken] = useState<string | null>(null);

  // Check auth on mount
  useEffect(() => {
    const authToken = localStorage.getItem("auth_token");
    
    if (authToken) {
      setToken(authToken);
      setIsAuthenticated(true);
    } else {
      setIsAuthenticated(false);
      router.push("/login");
    }
    
    setIsLoading(false);
  }, [router]);

  const logout = useCallback(() => {
    localStorage.removeItem("auth_token");
    apiClient.clearToken();
    setIsAuthenticated(false);
    setToken(null);
    router.push("/login");
  }, [router]);

  return {
    isAuthenticated,
    isLoading,
    token,
  };
}
