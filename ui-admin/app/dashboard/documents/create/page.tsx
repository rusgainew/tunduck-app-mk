"use client";

import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { organizationsApi } from "@/lib/api";
import { apiClient } from "@/lib/api-client";
import { DashboardLayout } from "@/components/DashboardLayout";
import { DocumentForm } from "@/components/DocumentForm";
import { ArrowLeft } from "lucide-react";
import { useState, useEffect } from "react";

export default function CreateDocumentPage() {
  const router = useRouter();
  const [selectedOrgId, setSelectedOrgId] = useState<string | null>(null);

  // Fetch organizations
  const { data: orgsData } = useQuery({
    queryKey: ["organizations"],
    queryFn: () => organizationsApi.getAll(),
  });

  // Set initial organization
  useEffect(() => {
    if (orgsData && orgsData.length > 0) {
      const orgId = selectedOrgId || orgsData[0].id;
      setSelectedOrgId(orgId);
      apiClient.setOrganizationId(orgId);
    }
  }, [orgsData, selectedOrgId]);

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
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                Создать новый документ ЭСФ
              </h1>
              <p className="text-gray-600 mt-1">
                Заполните форму для создания нового документа
              </p>
            </div>
          </div>
          {orgsData && orgsData.length > 0 && (
            <select
              value={selectedOrgId || ""}
              onChange={(e) => {
                setSelectedOrgId(e.target.value);
                apiClient.setOrganizationId(e.target.value);
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
        </div>

        {/* Form */}
        <DocumentForm onCancel={() => router.back()} />
      </div>
    </DashboardLayout>
  );
}
