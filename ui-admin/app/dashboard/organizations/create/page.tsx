"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { organizationsApi } from "@/lib/api";
import { ArrowLeft, Save } from "lucide-react";
import { DashboardLayout } from "@/components/DashboardLayout";
import Link from "next/link";

interface FormData {
  name: string;
  email: string;
  phone: string;
  address: string;
  description: string;
  isActive: boolean;
}

export default function CreateOrganizationPage() {
  const router = useRouter();
  const [formData, setFormData] = useState<FormData>({
    name: "",
    email: "",
    phone: "",
    address: "",
    description: "",
    isActive: true,
  });
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string>("");

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value, type } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]:
        type === "checkbox" ? (e.target as HTMLInputElement).checked : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      setIsSaving(true);
      setError("");

      const payload = {
        name: formData.name,
        description: formData.description,
        email: formData.email,
        phone: formData.phone,
        address: formData.address,
        isActive: formData.isActive,
      } as Omit<typeof formData, "name"> & {name: string};
      await organizationsApi.create(payload);
      router.push("/dashboard/organizations");
    } catch (err) {
      const errorMessage = 
        (err as { response?: { data?: { message?: string; error?: string } } })
          ?.response?.data?.message ||
        (err as { response?: { data?: { message?: string; error?: string } } })
          ?.response?.data?.error ||
        "Ошибка при создании организации";
      setError(errorMessage);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <DashboardLayout>
      <div className="max-w-2xl space-y-6">
        <div className="flex items-center gap-4">
          <Link
            href="/dashboard/organizations"
            className="p-2 hover:bg-gray-100 rounded transition-colors"
          >
            <ArrowLeft className="w-5 h-5 text-gray-600" />
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">
            Создание организации
          </h1>
        </div>

        {error && (
          <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="bg-white rounded-lg shadow-sm p-6 space-y-6">
            {/* Название */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Название организации *
              </label>
              <input
                name="name"
                type="text"
                placeholder="Введите название"
                value={formData.name}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              />
            </div>

            {/* Email */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Email *
              </label>
              <input
                name="email"
                type="email"
                placeholder="example@company.com"
                value={formData.email}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              />
            </div>

            {/* Телефон */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Телефон *
              </label>
              <input
                name="phone"
                type="tel"
                placeholder="+1 (555) 000-0000"
                value={formData.phone}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              />
            </div>

            {/* Адрес */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Адрес
              </label>
              <input
                name="address"
                type="text"
                placeholder="Введите адрес организации"
                value={formData.address}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              />
            </div>

            {/* Описание */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Описание
              </label>
              <textarea
                name="description"
                placeholder="Введите описание организации"
                rows={4}
                value={formData.description}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              />
            </div>

            {/* Статус */}
            <div className="flex items-center gap-3">
              <input
                name="isActive"
                type="checkbox"
                id="isActive"
                checked={formData.isActive}
                onChange={handleChange}
                className="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
              />
              <label htmlFor="isActive" className="text-sm font-medium text-gray-700">
                Организация активна
              </label>
            </div>
          </div>

          {/* Кнопки действия */}
          <div className="flex gap-4">
            <button
              type="submit"
              disabled={isSaving}
              className="flex items-center gap-2 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
            >
              <Save className="w-4 h-4" />
              {isSaving ? "Сохранение..." : "Создать"}
            </button>
            <Link
              href="/dashboard/organizations"
              className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium"
            >
              Отмена
            </Link>
          </div>
        </form>
      </div>
    </DashboardLayout>
  );
}
