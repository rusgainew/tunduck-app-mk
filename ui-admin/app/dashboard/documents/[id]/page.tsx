"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { documentsApi, organizationsApi } from "@/lib/api";
import { apiClient } from "@/lib/api-client";
import { useEffect, useState, useCallback } from "react";
import { ArrowLeft, Copy, Check, Edit2, Trash2, FileText, DollarSign, Calendar, Package, TrendingUp, Clock } from "lucide-react";
import { format } from "date-fns";
import { ru } from "date-fns/locale";
import { useToast } from "@/hooks/useToast";
import { DashboardLayout } from "@/components/DashboardLayout";
import { DocumentForm } from "@/components/DocumentForm";
import { LoadingState } from "@/components/states/LoadingState";
import { ErrorState } from "@/components/states/ErrorState";

export default function DocumentDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const { success } = useToast();
  const id = params.id as string;
  const [copiedId, setCopiedId] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [userSelectedOrgId, setUserSelectedOrgId] = useState<string | null>(null);

  // Fetch organizations
  const { data: orgsData } = useQuery({
    queryKey: ["organizations"],
    queryFn: () => organizationsApi.getAll(),
  });

  // Use first organization as default, or user's selection
  const selectedOrgId = userSelectedOrgId || orgsData?.[0]?.id || null;

  // Update organization ID in apiClient when it changes
  useEffect(() => {
    if (selectedOrgId) {
      apiClient.setOrganizationId(selectedOrgId);
    }
  }, [selectedOrgId]);

  const { data: document, isLoading, refetch } = useQuery({
    queryKey: ["document", id, selectedOrgId],
    queryFn: () => documentsApi.getById(id),
    enabled: !!id && !!selectedOrgId,
  });

  const handleCopyId = useCallback(async () => {
    if (document?.id) {
      await navigator.clipboard.writeText(document.id);
      setCopiedId(true);
      success("ID документа скопирован в буфер обмена");
      setTimeout(() => setCopiedId(false), 2000);
    }
  }, [document?.id, success]);

  const handleDelete = useCallback(async () => {
    if (confirm("Вы уверены, что хотите удалить этот документ? Это действие невозможно отменить.")) {
      try {
        const orgId = apiClient.getOrganizationId();
        await documentsApi.delete(id, orgId || undefined);
        success("Документ удален");
        router.push("/dashboard/documents");
      } catch (err) {
        console.error("Delete error:", err);
      }
    }
  }, [id, success, router]);

  if (isLoading) {
    return (
      <DashboardLayout>
        <LoadingState message="Загрузка документа..." fullHeight />
      </DashboardLayout>
    );
  }

  if (!document) {
    return (
      <DashboardLayout>
        <ErrorState 
          title="Документ не найден" 
          message="Запрошенный документ не существует или был удален." 
          onRetry={() => router.push("/dashboard/documents")}
        />
      </DashboardLayout>
    );
  }

  if (isEditing) {
    return (
      <DashboardLayout>
        <div className="space-y-6">
          <div className="flex items-center gap-4">
            <button
              onClick={() => setIsEditing(false)}
              className="text-gray-600 hover:text-gray-900 transition-colors"
            >
              <ArrowLeft className="w-6 h-6" />
            </button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                Редактировать документ
              </h1>
            </div>
          </div>
          <DocumentForm
            initialData={document}
            isEditing={true}
            onCancel={() => {
              setIsEditing(false);
              refetch();
            }}
          />
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <button
              onClick={() => router.back()}
              className="text-gray-600 hover:text-gray-900 transition-colors"
            >
              <ArrowLeft className="w-6 h-6" />
            </button>
            <div className="flex-1">
              <h1 className="text-3xl font-bold text-gray-900">
                {document.foreignName || "Документ ЭСФ"}
              </h1>
              <p className="text-gray-600 mt-1">
                {document.createdAt
                  ? format(new Date(document.createdAt), "dd MMMM yyyy HH:mm", {
                      locale: ru,
                    })
                  : "—"}
              </p>
            </div>
          </div>
          <div className="flex flex-col gap-2">
            {orgsData && orgsData.length > 0 && (
              <select
                value={selectedOrgId || ""}
                onChange={(e) => {
                  setUserSelectedOrgId(e.target.value);
                }}
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
              >
                <option value="">Выберите организацию</option>
                {orgsData.map((org) => (
                  <option key={org.id} value={org.id}>
                    {org.name}
                  </option>
                ))}
              </select>
            )}
            <div className="flex gap-2">
              <button
                onClick={() => setIsEditing(true)}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
              >
                <Edit2 className="w-4 h-4" />
                Редактировать
              </button>
              <button
                onClick={handleDelete}
                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors flex items-center gap-2"
              >
                <Trash2 className="w-4 h-4" />
                Удалить
              </button>
            </div>
          </div>
        </div>

        {/* Main Info */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Основная информация */}
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center gap-2 mb-4">
                <FileText className="w-5 h-5 text-blue-600" />
                <h2 className="text-lg font-semibold text-gray-900">
                  Основная информация
                </h2>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    ИНН контрагента
                  </label>
                  <p className="text-gray-900 font-mono">
                    {document.contractorTin || "—"}
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Валюта
                  </label>
                  <p className="text-gray-900">{document.currencyCode || "—"}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Сумма
                  </label>
                  <p className="text-gray-900 font-semibold">
                    {document.totalCurrencyValue
                      ? document.totalCurrencyValue.toLocaleString("ru-RU")
                      : "—"}{" "}
                    {document.currencyCode || ""}
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Тип резидентности
                  </label>
                  <span
                    className={`px-3 py-1 text-sm rounded-full ${
                      document.isResident
                        ? "bg-green-100 text-green-800"
                        : "bg-blue-100 text-blue-800"
                    }`}
                  >
                    {document.isResident ? "Резидент КР" : "Нерезидент"}
                  </span>
                </div>
              </div>
            </div>

            {/* Даты и операции */}
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center gap-2 mb-4">
                <Calendar className="w-5 h-5 text-green-600" />
                <h2 className="text-lg font-semibold text-gray-900">
                  Даты и операции
                </h2>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Дата доставки
                  </label>
                  <p className="text-gray-900">
                    {document.deliveryDate
                      ? format(new Date(document.deliveryDate), "dd MMMM yyyy", {
                          locale: ru,
                        })
                      : "—"}
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Дата контракта
                  </label>
                  <p className="text-gray-900">
                    {document.contractStartDate
                      ? format(
                          new Date(document.contractStartDate),
                          "dd MMMM yyyy",
                          {
                            locale: ru,
                          }
                        )
                      : "—"}
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Код типа доставки
                  </label>
                  <p className="text-gray-900 font-mono">
                    {document.deliveryTypeCode || "—"}
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Код типа платежа
                  </label>
                  <p className="text-gray-900 font-mono">
                    {document.paymentCode || "—"}
                  </p>
                </div>
              </div>
            </div>

            {/* Финансовая информация */}
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center gap-2 mb-4">
                <DollarSign className="w-5 h-5 text-purple-600" />
                <h2 className="text-lg font-semibold text-gray-900">
                  Финансовая информация
                </h2>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Сумма без налогов
                  </label>
                  <p className="text-gray-900">
                    {document.totalCurrencyValueWithoutTaxes
                      ? document.totalCurrencyValueWithoutTaxes.toLocaleString(
                          "ru-RU"
                        )
                      : "—"}
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Курс валюты
                  </label>
                  <p className="text-gray-900">{document.currencyRate || "—"}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Оплачено
                  </label>
                  <p className="text-gray-900">{document.paidAmount || "—"}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    К оплате
                  </label>
                  <p className="text-gray-900 font-semibold">
                    {document.amountToBePaid || "—"}
                  </p>
                </div>
              </div>
            </div>

            {/* Комментарий */}
            {document.comment && (
              <div className="bg-white rounded-lg shadow p-6">
                <div className="flex items-center gap-2 mb-4">
                  <FileText className="w-5 h-5 text-indigo-600" />
                  <h2 className="text-lg font-semibold text-gray-900">
                    Комментарий
                  </h2>
                </div>
                <p className="text-gray-700 whitespace-pre-wrap">
                  {document.comment}
                </p>
              </div>
            )}

            {/* Записи товаров/услуг */}
            {document.catalogEntries && document.catalogEntries.length > 0 && (
              <div className="bg-white rounded-lg shadow p-6">
                <div className="flex items-center gap-2 mb-4">
                  <Package className="w-5 h-5 text-orange-600" />
                  <h2 className="text-lg font-semibold text-gray-900">
                    Товары и услуги ({document.catalogEntries.length})
                  </h2>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-gray-200">
                        <th className="px-4 py-2 text-left text-gray-700 font-semibold">
                          Код каталога
                        </th>
                        <th className="px-4 py-2 text-left text-gray-700 font-semibold">
                          Название
                        </th>
                        <th className="px-4 py-2 text-right text-gray-700 font-semibold">
                          Кол-во
                        </th>
                        <th className="px-4 py-2 text-right text-gray-700 font-semibold">
                          Цена
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {document.catalogEntries.map((entry, idx) => (
                        <tr key={idx} className="hover:bg-gray-50">
                          <td className="px-4 py-2 font-mono text-gray-900">
                            {entry.catalogCode}
                          </td>
                          <td className="px-4 py-2 text-gray-900">
                            {entry.catalogName}
                          </td>
                          <td className="px-4 py-2 text-right text-gray-900">
                            {entry.quantity}
                          </td>
                          <td className="px-4 py-2 text-right text-gray-900">
                            {entry.price?.toLocaleString("ru-RU")}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow p-6 space-y-4 sticky top-20">
              <div className="border-b pb-4">
                <div className="flex items-center gap-2 mb-2">
                  <FileText className="w-4 h-4 text-gray-600" />
                  <label className="block text-sm font-medium text-gray-700">
                    ID документа
                  </label>
                </div>
                <div className="flex gap-2">
                  <div className="flex-1 bg-gray-50 rounded p-3 font-mono text-xs text-gray-900 break-all">
                    {document.id}
                  </div>
                  <button
                    onClick={handleCopyId}
                    className="px-3 py-3 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors flex items-center justify-center"
                  >
                    {copiedId ? (
                      <Check className="w-4 h-4" />
                    ) : (
                      <Copy className="w-4 h-4" />
                    )}
                  </button>
                </div>
              </div>

              <div className="border-b pb-4">
                <div className="flex items-center gap-2 mb-2">
                  <Clock className="w-4 h-4 text-gray-600" />
                  <label className="block text-sm font-medium text-gray-700">
                    Создано
                  </label>
                </div>
                <p className="text-sm text-gray-600">
                  {document.createdAt
                    ? format(new Date(document.createdAt), "dd MMMM yyyy HH:mm", {
                        locale: ru,
                      })
                    : "—"}
                </p>
              </div>

              <div className="border-b pb-4">
                <div className="flex items-center gap-2 mb-2">
                  <Clock className="w-4 h-4 text-gray-600" />
                  <label className="block text-sm font-medium text-gray-700">
                    Обновлено
                  </label>
                </div>
                <p className="text-sm text-gray-600">
                  {document.updatedAt
                    ? format(new Date(document.updatedAt), "dd MMMM yyyy HH:mm", {
                        locale: ru,
                      })
                    : "—"}
                </p>
              </div>

              <div className="pt-2">
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between items-center">
                    <div className="flex items-center gap-2">
                      <TrendingUp className="w-4 h-4 text-gray-600" />
                      <span className="text-gray-600">Статус отправки:</span>
                    </div>
                    <span className="font-medium text-gray-900">
                      {document.isBranchDataSent ? "✓ Отправлен" : "Не отправлен"}
                    </span>
                  </div>
                  <div className="flex justify-between items-center">
                    <div className="flex items-center gap-2">
                      <DollarSign className="w-4 h-4 text-gray-600" />
                      <span className="text-gray-600">Без налогов:</span>
                    </div>
                    <span className="font-medium text-gray-900">
                      {document.isPriceWithoutTaxes ? "Да" : "Нет"}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
