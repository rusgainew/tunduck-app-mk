"use client";

import { useRouter } from "next/navigation";
import { LoginForm } from "@/components/LoginForm";
import { LogIn } from "lucide-react";

export default function LoginPage() {
  const router = useRouter();

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4">
      <div className="max-w-md w-full bg-white rounded-lg shadow-xl p-8">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-blue-600 rounded-full mb-4">
            <LogIn className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900">Tunduck Admin</h1>
          <p className="text-gray-600 mt-2">Войдите в систему управления</p>
        </div>

        <LoginForm onSuccess={() => router.push("/dashboard")} />

        <div className="mt-6 text-center text-sm text-gray-600 space-y-3">
          <p>
            Нет аккаунта?{" "}
            <a href="/register" className="text-blue-600 hover:underline">
              Зарегистрироваться
            </a>
          </p>
          <p className="pt-2 border-t border-gray-200">
            Администратор?{" "}
            <a href="/register-admin" className="text-purple-600 hover:underline font-medium">
              Регистрация администратора
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
