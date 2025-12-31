"use client";

import { useState, useEffect, useCallback } from "react";
import { useRouter, useParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { usersApi } from "@/lib/api";
import { User } from "@/lib/types";
import { ArrowLeft, Save, Trash2 } from "lucide-react";
import { DashboardLayout } from "@/components/DashboardLayout";
import { LoadingState } from "@/components/states/LoadingState";
import { ErrorState } from "@/components/states/ErrorState";

const editSchema = z.object({
  email: z.string().email("Неверный формат email"),
  fullName: z.string().min(2, "Полное имя должно быть не менее 2 символов"),
  phone: z.string().min(10, "Номер телефона должен быть не менее 10 символов"),
  role: z.enum(["user", "viewer", "admin"], { message: "Выберите роль" }),
  isActive: z.boolean(),
});

type EditForm = z.infer<typeof editSchema>;

export default function EditUserPage() {
  const router = useRouter();
  const params = useParams();
  const userId = params.id as string;

  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState<string>("");
  const [successMessage, setSuccessMessage] = useState<string>("");

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<EditForm>({
    resolver: zodResolver(editSchema),
  });

  // Загружаем данные пользователя
  useEffect(() => {
    const fetchUser = async () => {
      try {
        setIsLoading(true);
        const data = await usersApi.getById(userId);
        setUser(data);
        reset({
          email: data.email,
          fullName: data.fullName,
          phone: data.phone,
          role: data.role,
          isActive: data.isActive,
        });
      } catch {
        setError("Ошибка загрузки пользователя");
      } finally {
        setIsLoading(false);
      }
    };

    fetchUser();
  }, [userId, reset]);

  const onSubmit = useCallback(
    async (data: EditForm) => {
      try {
        setIsSaving(true);
        setError("");
        setSuccessMessage("");

        await usersApi.update(userId, data);
        setSuccessMessage("Пользователь успешно обновлён");

        setTimeout(() => {
          router.push("/dashboard/users");
        }, 1500);
      } catch (err: unknown) {
        const errorMsg = err instanceof Error ? err.message : "Ошибка при обновлении пользователя";
        setError(
          (err as { response?: { data?: { message?: string; error?: string } } })?.response?.data?.message || errorMsg
        );
      } finally {
        setIsSaving(false);
      }
    },
    [userId, router]
  );

  const handleDelete = useCallback(async () => {
    if (!confirm("Вы уверены? Это действие нельзя отменить.")) {
      return;
    }

    try {
      setIsDeleting(true);
      setError("");
      await usersApi.delete(userId);
      router.push("/dashboard/users");
    } catch (err: unknown) {
      const errorMsg = err instanceof Error ? err.message : "Ошибка при удалении пользователя";
      setError((err as { response?: { data?: { message?: string; error?: string } } })?.response?.data?.message || errorMsg);
      setIsDeleting(false);
    }
  }, [userId, router]);

  if (isLoading) {
    return (
      <DashboardLayout>
        <LoadingState message="Загрузка пользователя..." fullHeight />
      </DashboardLayout>
    );
  }

  if (!user) {
    return (
      <DashboardLayout>
        <ErrorState 
          title="Пользователь не найден" 
          message="Запрошенный пользователь не существует или был удален." 
          onRetry={() => router.back()}
        />
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <button
              onClick={() => router.back()}
              className="inline-flex items-center justify-center w-10 h-10 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <ArrowLeft className="w-5 h-5 text-gray-600" />
            </button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Редактировать пользователя</h1>
              <p className="text-gray-600 mt-1">{user.username}</p>
            </div>
          </div>
        </div>

      {/* Error */}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800">{error}</p>
        </div>
      )}

      {/* Success */}
      {successMessage && (
        <div className="bg-green-50 border border-green-200 rounded-lg p-4">
          <p className="text-green-800">{successMessage}</p>
        </div>
      )}

      {/* Form */}
      <div className="bg-white rounded-lg shadow p-8">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* Email */}
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              {...register("email")}
              type="email"
              id="email"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="email@example.com"
            />
            {errors.email && (
              <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
            )}
          </div>

          {/* Full Name */}
          <div>
            <label htmlFor="fullName" className="block text-sm font-medium text-gray-700 mb-2">
              Полное имя
            </label>
            <input
              {...register("fullName")}
              type="text"
              id="fullName"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Иван Иванов"
            />
            {errors.fullName && (
              <p className="mt-1 text-sm text-red-600">{errors.fullName.message}</p>
            )}
          </div>

          {/* Phone */}
          <div>
            <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-2">
              Телефон
            </label>
            <input
              {...register("phone")}
              type="tel"
              id="phone"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="+996123456789"
            />
            {errors.phone && (
              <p className="mt-1 text-sm text-red-600">{errors.phone.message}</p>
            )}
          </div>

          {/* Role */}
          <div>
            <label htmlFor="role" className="block text-sm font-medium text-gray-700 mb-2">
              Роль
            </label>
            <select
              {...register("role")}
              id="role"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="user">Пользователь (user)</option>
              <option value="viewer">Просмотр (viewer)</option>
              <option value="admin">Администратор (admin)</option>
            </select>
            {errors.role && (
              <p className="mt-1 text-sm text-red-600">{errors.role.message}</p>
            )}
          </div>

          {/* IsActive */}
          <div className="flex items-center gap-3">
            <input
              {...register("isActive")}
              type="checkbox"
              id="isActive"
              className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <label htmlFor="isActive" className="text-sm font-medium text-gray-700">
              Активный пользователь
            </label>
          </div>

          {/* Buttons */}
          <div className="flex gap-4 pt-6 border-t border-gray-200">
            <button
              type="submit"
              disabled={isSaving}
              className="flex items-center gap-2 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Save className="w-4 h-4" />
              {isSaving ? "Сохранение..." : "Сохранить"}
            </button>

            <button
              type="button"
              onClick={handleDelete}
              disabled={isDeleting}
              className="flex items-center gap-2 px-6 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Trash2 className="w-4 h-4" />
              {isDeleting ? "Удаление..." : "Удалить"}
            </button>

            <button
              type="button"
              onClick={() => router.back()}
              className="flex items-center gap-2 px-6 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
            >
              Отменить
            </button>
          </div>
        </form>
      </div>

      {/* User Info */}
      <div className="bg-gray-50 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Информация</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-600">Username:</span>
            <p className="font-medium text-gray-900">{user.username}</p>
          </div>
          <div>
            <span className="text-gray-600">ID:</span>
            <p className="font-medium text-gray-900 break-all">{user.id}</p>
          </div>
          <div>
            <span className="text-gray-600">Создан:</span>
            <p className="font-medium text-gray-900">
              {new Date(user.createdAt).toLocaleString("ru-RU")}
            </p>
          </div>
          <div>
            <span className="text-gray-600">Обновлён:</span>
            <p className="font-medium text-gray-900">
              {new Date(user.updatedAt).toLocaleString("ru-RU")}
            </p>
          </div>
        </div>
      </div>
      </div>
    </DashboardLayout>
  );
}
