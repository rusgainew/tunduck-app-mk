"use client";

import Link from "next/link";
import { useState, useCallback, memo } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/store";
import { authApi } from "@/lib/api";
import { LogOut, User, Shield, ChevronDown, Users, Building2 } from "lucide-react";

/**
 * Компонент header'а с информацией пользователя
 * Мемоизирован для предотвращения лишних переrender'ов
 */
function HeaderComponent() {
  const router = useRouter();
  const { user, clearAuth } = useAuthStore();
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const handleLogout = useCallback(async () => {
    try {
      setIsLoading(true);
      await authApi.logout();
      clearAuth();
      router.push("/login");
    } catch (error) {
      console.error("Logout error:", error);
      clearAuth();
      router.push("/login");
    }
  }, [clearAuth, router]);

  const toggleDropdown = useCallback(() => {
    setIsDropdownOpen(prev => !prev);
  }, []);

  if (!user) return null;

  const isAdmin = user.role === "admin";
  const roleColor = isAdmin ? "bg-purple-100 text-purple-700" : "bg-blue-100 text-blue-700";
  const roleIcon = isAdmin ? <Shield className="w-4 h-4" /> : <User className="w-4 h-4" />;

  return (
    <header className="bg-white shadow border-b border-gray-200">
      <div className="px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold text-gray-900">Tunduck Admin</h1>
          {isAdmin && (
            <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-purple-100 text-purple-700 text-sm font-semibold">
              <Shield className="w-4 h-4" />
              АДМИНИСТРАТОР
            </span>
          )}
        </div>

        {/* User Menu */}
        <div className="relative">
          <button
            onClick={toggleDropdown}
            className="flex items-center gap-3 px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
            aria-label="User menu"
            aria-expanded={isDropdownOpen}
          >
            <div className="text-right">
              <p className="text-sm font-medium text-gray-900">{user.fullName}</p>
              <p className="text-xs text-gray-600">@{user.username}</p>
              <p className={`text-xs px-2 py-1 rounded inline-flex items-center gap-1 mt-1 ${roleColor}`}>
                {roleIcon}
                {user.role === "admin" ? "Администратор" : user.role === "user" ? "Пользователь" : "Просмотр"}
              </p>
            </div>
            <ChevronDown className="w-4 h-4 text-gray-600" />
          </button>

          {/* Dropdown Menu */}
          {isDropdownOpen && (
            <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg z-50 border border-gray-200"
              role="menu"
              aria-orientation="vertical"
            >
              <div className="p-4 border-b border-gray-200">
                <p className="text-sm text-gray-600">Username:</p>
                <p className="text-sm font-medium text-gray-900 break-all">@{user.username}</p>
                <p className="text-sm text-gray-600 mt-2">Email:</p>
                <p className="text-sm font-medium text-gray-900 break-all">{user.email}</p>
                {isAdmin && (
                  <>
                    <p className="text-sm text-gray-600 mt-2">ID:</p>
                    <p className="text-xs text-gray-500 break-all font-mono">{user.id}</p>
                  </>
                )}
              </div>

              <nav className="py-2">
                {isAdmin && (
                  <>
                    <Link
                      href="/dashboard/users"
                      className="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                      <Users className="w-4 h-4" />
                      Управление пользователями
                    </Link>
                    <Link
                      href="/dashboard/organizations"
                      className="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                      <Building2 className="w-4 h-4" />
                      Управление организациями
                    </Link>
                  </>
                )}
                <Link
                  href={`/dashboard/users/${user.id}`}
                  className="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  <User className="w-4 h-4" />
                  Мой профиль
                </Link>
              </nav>

              <div className="border-t border-gray-200 p-2">
                <button
                  onClick={handleLogout}
                  disabled={isLoading}
                  className="w-full flex items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 rounded transition-colors disabled:opacity-50"
                  role="menuitem"
                >
                  <LogOut className="w-4 h-4" />
                  {isLoading ? "Выход..." : "Выход"}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}

/**
 * Экспортируем мемоизированный компонент Header
 */
export const Header = memo(HeaderComponent);
