"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { authApi } from "@/lib/api";
import { useAuthStore } from "@/lib/store";
import { ShieldPlus } from "lucide-react";

const registerAdminSchema = z
  .object({
    username: z
      .string()
      .min(3, "Имя пользователя должно быть не менее 3 символов"),
    email: z.string().email("Неверный формат email"),
    fullName: z.string().min(2, "Полное имя должно быть не менее 2 символов"),
    phone: z
      .string()
      .min(10, "Номер телефона должен быть не менее 10 символов"),
    password: z.string().min(6, "Пароль должен быть не менее 6 символов"),
    confirmPassword: z.string(),
    adminSecret: z
      .string()
      .min(1, "Секретный код администратора обязателен"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Пароли не совпадают",
    path: ["confirmPassword"],
  });

type RegisterAdminForm = z.infer<typeof registerAdminSchema>;

export default function RegisterAdminPage() {
  const router = useRouter();
  const setAuth = useAuthStore((state) => state.setAuth);
  const [error, setError] = useState<string>("");
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterAdminForm>({
    resolver: zodResolver(registerAdminSchema),
  });

  const onSubmit = async (data: RegisterAdminForm) => {
    try {
      setIsLoading(true);
      setError("");

      // Отправляем админа через специальный endpoint регистрации администратора
      const response = await authApi.registerAdmin({
        username: data.username,
        email: data.email,
        fullName: data.fullName,
        phone: data.phone,
        password: data.password,
        confirmPassword: data.confirmPassword,
        adminSecret: data.adminSecret,
      });

      setAuth(response.user, response.token);
      router.push("/dashboard");
    } catch (err: unknown) {
      console.error("Register admin error:", err);
      const errorMsg =
        (err as { response?: { data?: { message?: string; error?: string } } })?.response?.data?.message ||
        (err as { response?: { data?: { message?: string; error?: string } } })?.response?.data?.error ||
        "Ошибка регистрации администратора. Попробуйте снова.";
      setError(errorMsg);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-50 to-indigo-100 px-4 py-8">
      <div className="max-w-md w-full bg-white rounded-lg shadow-xl p-8">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-purple-600 rounded-full mb-4">
            <ShieldPlus className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900">
            Регистрация Администратора
          </h1>
          <p className="text-gray-600 mt-2">Создайте аккаунт администратора</p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-800 text-sm">{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Имя пользователя
            </label>
            <input
              {...register("username")}
              type="text"
              id="username"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              placeholder="admin_username"
            />
            {errors.username && (
              <p className="mt-1 text-sm text-red-600">
                {errors.username.message}
              </p>
            )}
          </div>

          <div>
            <label
              htmlFor="email"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Email
            </label>
            <input
              {...register("email")}
              type="email"
              id="email"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              placeholder="admin@example.com"
            />
            {errors.email && (
              <p className="mt-1 text-sm text-red-600">
                {errors.email.message}
              </p>
            )}
          </div>

          <div>
            <label
              htmlFor="fullName"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Полное имя
            </label>
            <input
              {...register("fullName")}
              type="text"
              id="fullName"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              placeholder="Иван Администратор"
            />
            {errors.fullName && (
              <p className="mt-1 text-sm text-red-600">
                {errors.fullName.message}
              </p>
            )}
          </div>

          <div>
            <label
              htmlFor="phone"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Телефон
            </label>
            <input
              {...register("phone")}
              type="tel"
              id="phone"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              placeholder="+996123456789"
            />
            {errors.phone && (
              <p className="mt-1 text-sm text-red-600">
                {errors.phone.message}
              </p>
            )}
          </div>

          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Пароль
            </label>
            <input
              {...register("password")}
              type="password"
              id="password"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              placeholder="••••••••"
            />
            {errors.password && (
              <p className="mt-1 text-sm text-red-600">
                {errors.password.message}
              </p>
            )}
          </div>

          <div>
            <label
              htmlFor="confirmPassword"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Подтвердите пароль
            </label>
            <input
              {...register("confirmPassword")}
              type="password"
              id="confirmPassword"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              placeholder="••••••••"
            />
            {errors.confirmPassword && (
              <p className="mt-1 text-sm text-red-600">
                {errors.confirmPassword.message}
              </p>
            )}
          </div>

          <div>
            <label
              htmlFor="adminSecret"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Секретный код администратора
            </label>
            <input
              {...register("adminSecret")}
              type="password"
              id="adminSecret"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              placeholder="••••••••"
            />
            {errors.adminSecret && (
              <p className="mt-1 text-sm text-red-600">
                {errors.adminSecret.message}
              </p>
            )}
            <p className="mt-1 text-xs text-gray-500">
              Введите секретный код для подтверждения административных прав
            </p>
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-purple-600 text-white py-3 rounded-lg hover:bg-purple-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed font-medium mt-6"
          >
            {isLoading ? "Регистрация..." : "Создать администратора"}
          </button>
        </form>

        <div className="mt-6 text-center text-sm text-gray-600">
          <p>
            Уже есть аккаунт?{" "}
            <a href="/login" className="text-purple-600 hover:underline">
              Войти
            </a>
          </p>
          <p className="mt-2">
            Обычная регистрация?{" "}
            <a href="/register" className="text-blue-600 hover:underline">
              Регистрация пользователя
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
