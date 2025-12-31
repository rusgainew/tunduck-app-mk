"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { registerSchema, type RegisterFormData } from "@/lib/validation";
import { ErrorMessage } from "./FormField";
import { useRouter } from "next/navigation";
import { apiClient } from "@/lib/api-client";
import { useToast } from "@/hooks/useToast";
import { UserPlus, Eye, EyeOff } from "lucide-react";

interface RegisterFormProps {
  onSuccess?: () => void;
}

/**
 * Компонент формы регистрации
 * Содержит сложную валидацию пароля и toggle для показа пароля
 */
export function RegisterForm({ onSuccess }: RegisterFormProps) {
  const router = useRouter();
  const { success, error: showError } = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });

  const password = watch("password");

  const onSubmit = async (data: RegisterFormData) => {
    try {
      setIsLoading(true);
      const { confirmPassword, ...submitData } = data;
      
      const response = await apiClient.post<{ token: string }>("/auth/register", submitData);
      
      if (response.data?.token) {
        apiClient.setToken(response.data.token);
        success("Вы успешно зарегистрировались");
        onSuccess?.();
        router.push("/dashboard");
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Ошибка регистрации";
      showError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const renderFormField = (
    fieldName: keyof RegisterFormData,
    label: string,
    type: string = "text",
    placeholder: string = ""
  ) => (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-1">
        {label}
        <span className="text-red-600 ml-1">*</span>
      </label>
      <input
        type={type}
        {...register(fieldName)}
        className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
          errors[fieldName]
            ? "border-red-500 bg-red-50"
            : "border-gray-300 bg-white"
        }`}
        placeholder={placeholder}
        disabled={isLoading}
        aria-label={label}
      />
      <ErrorMessage error={errors[fieldName]} />
    </div>
  );

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {renderFormField("email", "Email", "email", "your@email.com")}
      {renderFormField("username", "Имя пользователя", "text", "username")}
      {renderFormField("fullName", "Полное имя", "text", "Иван Иванов")}

      {/* Password field with visibility toggle and requirements */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Пароль
          <span className="text-red-600 ml-1">*</span>
        </label>
        <div className="relative">
          <input
            type={showPassword ? "text" : "password"}
            {...register("password")}
            className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition pr-12 ${
              errors.password
                ? "border-red-500 bg-red-50"
                : "border-gray-300 bg-white"
            }`}
            placeholder="Минимум 8 символов"
            disabled={isLoading}
            aria-label="Пароль"
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 transition-colors"
            aria-label={showPassword ? "Скрыть пароль" : "Показать пароль"}
            tabIndex={-1}
          >
            {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
          </button>
        </div>
        <ErrorMessage error={errors.password} />
        {password && (
          <div className="mt-2 text-xs space-y-1">
            <p className={password.length >= 8 ? "text-green-600" : "text-gray-500"}>
              ✓ Минимум 8 символов
            </p>
            <p className={/[A-Z]/.test(password) ? "text-green-600" : "text-gray-500"}>
              ✓ Заглавная буква
            </p>
            <p className={/[a-z]/.test(password) ? "text-green-600" : "text-gray-500"}>
              ✓ Строчная буква
            </p>
            <p className={/\d/.test(password) ? "text-green-600" : "text-gray-500"}>
              ✓ Цифра
            </p>
          </div>
        )}
      </div>

      {/* Confirm password field with visibility toggle */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Подтверждение пароля
          <span className="text-red-600 ml-1">*</span>
        </label>
        <div className="relative">
          <input
            type={showConfirmPassword ? "text" : "password"}
            {...register("confirmPassword")}
            className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition pr-12 ${
              errors.confirmPassword
                ? "border-red-500 bg-red-50"
                : "border-gray-300 bg-white"
            }`}
            placeholder="Повторите пароль"
            disabled={isLoading}
            aria-label="Подтверждение пароля"
          />
          <button
            type="button"
            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 transition-colors"
            aria-label={showConfirmPassword ? "Скрыть пароль" : "Показать пароль"}
            tabIndex={-1}
          >
            {showConfirmPassword ? <EyeOff size={18} /> : <Eye size={18} />}
          </button>
        </div>
        <ErrorMessage error={errors.confirmPassword} />
      </div>

      <button
        type="submit"
        disabled={isLoading}
        className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        aria-busy={isLoading}
      >
        <UserPlus size={18} aria-hidden="true" />
        {isLoading ? "Регистрация..." : "Зарегистрироваться"}
      </button>
    </form>
  );
}
