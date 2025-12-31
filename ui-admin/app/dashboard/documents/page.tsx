"use client";

import { useState, useEffect, useCallback, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { documentsApi, organizationsApi } from "@/lib/api";
import { EsfDocument } from "@/lib/types";
import { Search, FileText, Calendar, Eye, Trash2 } from "lucide-react";
import { format } from "date-fns";
import { ru } from "date-fns/locale";
import { useToast } from "@/hooks/useToast";
import { apiClient } from "@/lib/api-client";
import { DashboardLayout } from "@/components/DashboardLayout";
import { LoadingState } from "@/components/states/LoadingState";
import { ErrorState } from "@/components/states/ErrorState";
import { EmptyState } from "@/components/states/EmptyState";
import Link from "next/link";

export default function DocumentsPage() {
  const router = useRouter();
  const { success, error: showError } = useToast();
  const [searchTerm, setSearchTerm] = useState("");
  const [userSelectedOrgId, setUserSelectedOrgId] = useState<string | null>(null);
  const [isAuthorized, setIsAuthorized] = useState(true);

  // Check authentication on mount
  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      setIsAuthorized(false);
      router.push('/login');
    }
  }, [router]);

  // Fetch organizations
  const { data: orgsData, error: orgsError } = useQuery({
    queryKey: ["organizations"],
    queryFn: () => organizationsApi.getAll(),
    retry: 1,
  });

  // Handle auth error from organizations query
  useEffect(() => {
    if (orgsError) {
      const axiosError = orgsError as any;
      if (axiosError.response?.status === 401) {
        localStorage.removeItem('auth_token');
        router.push('/login');
      }
    }
  }, [orgsError, router]);

  // Use first organization as default, or user's selection
  const selectedOrgId = userSelectedOrgId || orgsData?.[0]?.id || null;

  // Update apiClient when organization changes
  useEffect(() => {
    if (selectedOrgId) {
      apiClient.setOrganizationId(selectedOrgId);
    }
  }, [selectedOrgId]);

  const { data, isLoading, error: docsError, refetch } = useQuery({
    queryKey: ["documents", selectedOrgId],
    queryFn: () => documentsApi.getAll(),
    enabled: !!selectedOrgId && isAuthorized,
    retry: 1,
  });

  // Handle auth error from documents query
  useEffect(() => {
    if (docsError) {
      const axiosError = docsError as any;
      if (axiosError.response?.status === 401) {
        localStorage.removeItem('auth_token');
        router.push('/login');
      }
    }
  }, [docsError, router]);

  if (!isAuthorized) {
    return null;
  }

  const filteredDocs = useMemo(
    () =>
      data?.data?.filter(
        (doc: EsfDocument) =>
          doc.contractorTin?.toLowerCase().includes(searchTerm.toLowerCase()) ||
          doc.foreignName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
          doc.comment?.toLowerCase().includes(searchTerm.toLowerCase())
      ),
    [data?.data, searchTerm]
  );

  const handleSearchChange = useCallback((value: string) => {
    setSearchTerm(value);
  }, []);

  const handleOrgChange = useCallback(
    (orgId: string) => {
      setUserSelectedOrgId(orgId);
    },
    []
  );

  const handleDocumentDelete = useCallback(
    async (doc: EsfDocument) => {
      if (confirm("Вы уверены, что хотите удалить этот документ?")) {
        try {
          const orgId = apiClient.getOrganizationId();
          await documentsApi.delete(doc.id, orgId || undefined);
          success("Документ удален");
          // Refresh the documents list
          window.location.reload();
        } catch (err) {
          showError("Ошибка при удалении документа");
        }
      }
    },
    [success, showError]
  );

  const handleRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isLoading) {
    return (
      <DashboardLayout>
        <LoadingState message="Загрузка документов..." fullHeight />
      </DashboardLayout>
    );
  }

  if (docsError) {
    return (
      <DashboardLayout>
        <ErrorState 
          title="Ошибка загрузки" 
          message="Не удалось загрузить документы. Попробуйте снова." 
          onRetry={handleRetry}
        />
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start gap-4">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Документы ЭСФ</h2>
          <p className="text-gray-600 mt-1">
            Всего документов: {data?.count || 0}
          </p>
        </div>
        <div className="flex flex-col gap-2">
          {orgsData && orgsData.length > 0 && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Организация
              </label>
              <select
                value={selectedOrgId || ""}
                onChange={(e) => {
                  handleOrgChange(e.target.value);
                }}
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">Выберите организацию</option>
                {orgsData.map((org) => (
                  <option key={org.id} value={org.id}>
                    {org.name}
                  </option>
                ))}
              </select>
            </div>
          )}
          <Link href="/dashboard/documents/create" className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-center">
            + Создать документ
          </Link>
        </div>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
        <input
          type="text"
          placeholder="Поиск по ИНН, названию или комментарию..."
          value={searchTerm}
          onChange={(e) => handleSearchChange(e.target.value)}
          className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>

      {/* Documents Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredDocs && filteredDocs.length > 0 ? (
          filteredDocs.map((doc: EsfDocument, index: number) => (
            <div
              key={`${doc.id}-${index}`}
              className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-6"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center flex-1">
                  <FileText className="w-5 h-5 text-blue-600 mr-2 flex-shrink-0" />
                  <h3 className="text-sm font-semibold text-gray-900 line-clamp-2">
                    {doc.foreignName || "Без названия"}
                  </h3>
                </div>
                <div className="flex space-x-2 ml-2">
                  <Link
                    href={`/dashboard/documents/${doc.id}`}
                    onClick={() => {
                      const orgToken = apiClient.getOrganizationToken();
                      if (orgToken) {
                        localStorage.setItem("doc_org_token", orgToken);
                      }
                    }}
                    className="text-green-600 hover:text-green-700 transition-colors flex-shrink-0"
                    title="Просмотреть детали"
                  >
                    <Eye className="w-4 h-4" />
                  </Link>
                  <button
                    onClick={() => handleDocumentDelete(doc)}
                    className="text-red-600 hover:text-red-700 transition-colors flex-shrink-0"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>

              <div className="space-y-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-500">ИНН контрагента:</span>
                  <span className="font-medium text-gray-900">
                    {doc.contractorTin || "—"}
                  </span>
                </div>

                <div className="flex justify-between">
                  <span className="text-gray-500">Валюта:</span>
                  <span className="font-medium text-gray-900">
                    {doc.currencyCode || "—"}
                  </span>
                </div>

                {doc.totalCurrencyValue && (
                  <div className="flex justify-between">
                    <span className="text-gray-500">Сумма:</span>
                    <span className="font-semibold text-gray-900">
                      {doc.totalCurrencyValue.toLocaleString("ru-RU")}{" "}
                      {doc.currencyCode}
                    </span>
                  </div>
                )}

                <div className="flex items-center text-gray-500 pt-2 border-t">
                  <Calendar className="w-4 h-4 mr-1" />
                  <span className="text-xs">
                    {doc.deliveryDate
                      ? format(new Date(doc.deliveryDate), "dd.MM.yyyy", {
                          locale: ru,
                        })
                      : "—"}
                  </span>
                </div>

                {doc.isResident !== undefined && (
                  <div className="pt-2">
                    <span
                      className={`px-2 py-1 text-xs rounded-full ${
                        doc.isResident
                          ? "bg-green-100 text-green-800"
                          : "bg-blue-100 text-blue-800"
                      }`}
                    >
                      {doc.isResident ? "Резидент КР" : "Нерезидент"}
                    </span>
                  </div>
                )}

                {doc.comment && (
                  <div className="pt-2">
                    <p className="text-xs text-gray-600 line-clamp-2">
                      {doc.comment}
                    </p>
                  </div>
                )}
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-3">
            <EmptyState 
              title="Документы не найдены" 
              message={searchTerm ? "Попробуйте изменить параметры поиска" : "Нет документов для отображения"}
            />
          </div>
        )}
      </div>
      </div>
    </DashboardLayout>
  );
}
