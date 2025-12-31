"use client";

import { useState } from "react";
import { useForm, FieldValues, UseFormRegister, FormState } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRouter } from "next/navigation";
import { useToast } from "@/hooks/useToast";
import { ErrorMessage } from "./FormField";
import { LucideIcon } from "lucide-react";

/**
 * Конфигурация поля для формы аутентификации
 */
interface AuthFormFieldConfig {
  name: string;
  label: string;
  type?: "text" | "email" | "password";
  placeholder?: string;
  required?: boolean;
  hint?: string;
}

/**
 * Пропсы для универсальной формы аутентификации
 */
interface AuthFormProps<T extends FieldValues> {
  schema: any; // ZodSchema with proper type inference
  fields: AuthFormFieldConfig[];
  onSubmit: (data: T) => Promise<void>;
  buttonLabel: string;
  buttonIcon?: LucideIcon;
  successMessage?: string;
  redirectTo?: string;
  onSuccess?: () => void;
}

/**
 * Универсальный компонент для форм аутентификации (логин, регистрация и т.д.)
 * Избавляет от дублирования кода между LoginForm и RegisterForm
 */
export function AuthForm<T extends FieldValues>({
  schema,
  fields,
  onSubmit,
  buttonLabel,
  buttonIcon: ButtonIcon,
  successMessage = "Операция выполнена успешно",
  redirectTo = "/dashboard",
  onSuccess,
}: AuthFormProps<T>) {
  const router = useRouter();
  const { success, error: showError } = useToast();
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<T>({
    resolver: zodResolver(schema),
  });

  const handleFormSubmit = async (data: T) => {
    try {
      setIsLoading(true);
      await onSubmit(data);
      success(successMessage);
      onSuccess?.();
      if (redirectTo) {
        router.push(redirectTo);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Произошла ошибка";
      showError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      {fields.map((field) => (
        <div key={field.name}>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            {field.label}
            {field.required && <span className="text-red-600 ml-1">*</span>}
          </label>
          <input
            type={field.type || "text"}
            {...register(field.name as any)}
            className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
              errors[field.name]
                ? "border-red-500 bg-red-50"
                : "border-gray-300 bg-white"
            }`}
            placeholder={field.placeholder}
            aria-label={field.label}
            aria-describedby={errors[field.name] ? `${field.name}-error` : undefined}
            disabled={isLoading}
          />
          {field.hint && <p className="text-xs text-gray-500 mt-1">{field.hint}</p>}
          <ErrorMessage 
            error={errors[field.name] as any}
            className="mt-1"
          />
        </div>
      ))}

      <button
        type="submit"
        disabled={isLoading}
        className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        aria-busy={isLoading}
      >
        {ButtonIcon && <ButtonIcon size={18} aria-hidden="true" />}
        {isLoading ? "Обработка..." : buttonLabel}
      </button>
    </form>
  );
}
