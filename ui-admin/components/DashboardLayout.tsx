"use client";

import { useEffect, memo, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/store";
import { Header } from "@/components/Header";
import { Loader2 } from "lucide-react";

const DashboardLayoutComponent = ({ children }: { children: React.ReactNode }) => {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();

  const handleAuthRedirect = useCallback(() => {
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  useEffect(() => {
    handleAuthRedirect();
  }, [handleAuthRedirect]);

  if (!isAuthenticated) {
    return (
      <div 
        className="flex items-center justify-center h-screen"
        role="status"
        aria-label="Loading"
      >
        <Loader2 className="w-12 h-12 text-blue-500 animate-spin" aria-hidden="true" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main className="p-6 max-w-7xl mx-auto">{children}</main>
    </div>
  );
};

export const DashboardLayout = memo(DashboardLayoutComponent);
