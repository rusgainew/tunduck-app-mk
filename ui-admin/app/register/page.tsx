"use client";

import { useRouter } from "next/navigation";
import { RegisterForm } from "@/components/RegisterForm";
import { UserPlus } from "lucide-react";

export default function RegisterPage() {
  const router = useRouter();

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-emerald-100 px-4">
      <div className="max-w-md w-full bg-white rounded-lg shadow-xl p-8">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-green-600 rounded-full mb-4">
            <UserPlus className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900">Tunduck Admin</h1>
          <p className="text-gray-600 mt-2">Создайте новый аккаунт</p>
        </div>

        <RegisterForm onSuccess={() => router.push("/dashboard")} />

        <div className="mt-6 text-center text-sm text-gray-600">
          <p>
            Уже есть аккаунт?{" "}
            <a href="/login" className="text-blue-600 hover:underline">
              Войти
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
