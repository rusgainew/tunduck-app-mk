"use client";

import { loginSchema, type LoginFormData } from "@/lib/validation";
import { authApi } from "@/lib/api";
import { LogIn } from "lucide-react";
import { AuthForm } from "./AuthForm";

interface LoginFormProps {
  onSuccess?: () => void;
}

/**
 * Компонент формы входа в систему
 * Использует универсальный AuthForm компонент
 */
export function LoginForm({ onSuccess }: LoginFormProps) {
  return (
    <AuthForm<LoginFormData>
      schema={loginSchema}
      fields={[
        {
          name: "username",
          label: "Username",
          type: "text",
          placeholder: "Введите ваш username",
          required: true,
        },
        {
          name: "password",
          label: "Пароль",
          type: "password",
          placeholder: "Введите пароль",
          required: true,
        },
      ]}
      onSubmit={async (data) => {
        await authApi.login(data);
      }}
      buttonLabel="Войти"
      buttonIcon={LogIn}
      successMessage="Вы успешно вошли в систему"
      redirectTo="/dashboard"
      onSuccess={onSuccess}
    />
  );
}
