"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { organizationSchema, type OrganizationFormData } from "@/lib/validation";
import { ErrorMessage } from "./FormField";
import { useToast } from "@/hooks/useToast";
import { Save, X } from "lucide-react";

interface OrganizationFormProps {
  initialData?: Partial<OrganizationFormData>;
  isEditing?: boolean;
  onCancel?: () => void;
  onSubmit: (data: OrganizationFormData) => Promise<void>;
}

export function OrganizationForm({
  initialData,
  isEditing = false,
  onCancel,
  onSubmit: onSubmitProp,
}: OrganizationFormProps) {
  const { success, error: showError } = useToast();
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<OrganizationFormData>({
    resolver: zodResolver(organizationSchema),
    defaultValues: initialData,
  });

  const onSubmit = async (data: OrganizationFormData) => {
    try {
      setIsLoading(true);
      await onSubmitProp(data);
      success(isEditing ? "Организация обновлена" : "Организация создана");
    } catch (err: any) {
      showError(err.message || "Ошибка при сохранении организации");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Название организации *
        </label>
        <input
          type="text"
          {...register("name")}
          className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
            errors.name
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white"
          }`}
          placeholder="Введите название"
        />
        <ErrorMessage error={errors.name} />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Описание
        </label>
        <textarea
          {...register("description")}
          rows={3}
          className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
            errors.description
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white"
          }`}
          placeholder="Описание организации (опционально)"
        />
        <ErrorMessage error={errors.description} />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Email
        </label>
        <input
          type="email"
          {...register("email")}
          className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
            errors.email
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white"
          }`}
          placeholder="org@example.com"
        />
        <ErrorMessage error={errors.email} />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Телефон
        </label>
        <input
          type="tel"
          {...register("phone")}
          className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
            errors.phone
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white"
          }`}
          placeholder="+996 700 123 456"
        />
        <ErrorMessage error={errors.phone} />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Адрес
        </label>
        <input
          type="text"
          {...register("address")}
          className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
            errors.address
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white"
          }`}
          placeholder="Адрес организации"
        />
        <ErrorMessage error={errors.address} />
      </div>

      <div className="flex gap-2 justify-end pt-4">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors flex items-center gap-2"
          >
            <X size={18} />
            Отмена
          </button>
        )}
        <button
          type="submit"
          disabled={isLoading}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2 disabled:opacity-50"
        >
          <Save size={18} />
          {isLoading ? "Сохранение..." : isEditing ? "Обновить" : "Создать"}
        </button>
      </div>
    </form>
  );
}
