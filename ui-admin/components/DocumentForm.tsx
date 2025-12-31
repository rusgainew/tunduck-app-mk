"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { EsfDocument, EsfEntry } from "@/lib/types";
import { documentsApi } from "@/lib/api";
import { useToast } from "@/hooks/useToast";
import { useRouter } from "next/navigation";
import { apiClient } from "@/lib/api-client";
import { Save, X, Plus, Trash2 } from "lucide-react";
import { documentSchema, type DocumentFormData, entrySchema } from "@/lib/validation";
import { ErrorMessage } from "./FormField";

interface DocumentFormProps {
  initialData?: Partial<EsfDocument>;
  isEditing?: boolean;
  onCancel?: () => void;
}

export function DocumentForm({
  initialData,
  isEditing = false,
  onCancel,
}: DocumentFormProps) {
  const { success, error: showError } = useToast();
  const router = useRouter();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [entries, setEntries] = useState<Partial<EsfEntry>[]>(
    initialData?.catalogEntries || []
  );
  const [showEntryForm, setShowEntryForm] = useState(false);
  const [editingEntryIndex, setEditingEntryIndex] = useState<number | null>(null);
  const [currentEntry, setCurrentEntry] = useState<Partial<EsfEntry>>({});

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<DocumentFormData>({
    resolver: zodResolver(documentSchema),
    defaultValues: initialData
      ? {
          foreignName: initialData.foreignName || "",
          contractorTin: initialData.contractorTin || "",
          currencyCode: initialData.currencyCode || "",
          totalCurrencyValue: initialData.totalCurrencyValue || 0,
          deliveryDate: initialData.deliveryDate || "",
          contractStartDate: initialData.contractStartDate || "",
          comment: initialData.comment || "",
          isResident: initialData.isResident || false,
          isBranchDataSent: initialData.isBranchDataSent || false,
          isPriceWithoutTaxes: initialData.isPriceWithoutTaxes || false,
        }
      : {},
  });

  const onSubmit = async (data: DocumentFormData) => {
    try {
      setIsSubmitting(true);
      const orgId = apiClient.getOrganizationId();

      // Convert nulls to undefined for API
      const cleanData = Object.entries(data).reduce((acc, [key, value]) => {
        if (value !== null && value !== "") {
          const typedAcc = acc as Record<string, unknown>;
          typedAcc[key] = value;
        }
        return acc;
      }, {} as Partial<EsfDocument>);

      // Add entries to the document
      if (entries.length > 0) {
        cleanData.catalogEntries = entries as EsfEntry[];
      }

      if (isEditing && initialData?.id) {
        await documentsApi.update(initialData.id, cleanData, orgId || undefined);
        success("Документ обновлён успешно");
      } else {
        const newDoc = await documentsApi.create(cleanData, orgId || undefined);
        success("Документ создан успешно");
        router.push(`/dashboard/documents/${newDoc.id}`);
      }
    } catch (error: Error | unknown) {
      const errorMessage = error instanceof Error ? error.message : "Неизвестная ошибка";
      showError(
        isEditing
          ? `Ошибка при обновлении документа: ${errorMessage}`
          : `Ошибка при создании документа: ${errorMessage}`
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleAddEntry = () => {
    setCurrentEntry({});
    setEditingEntryIndex(null);
    setShowEntryForm(true);
  };

  const handleSaveEntry = () => {
    // Validate entry before saving
    const entryValidation = entrySchema.safeParse(currentEntry);
    
    if (!entryValidation.success) {
      const errors = entryValidation.error.flatten().fieldErrors;
      const errorMessages = Object.entries(errors)
        .map(([field, msgs]) => {
          const fieldNames: { [key: string]: string } = {
            catalogCode: "Код каталога",
            catalogName: "Название",
            measureCode: "Код измерения",
            quantity: "Количество",
            price: "Цена",
          };
          const fieldLabel = fieldNames[field] || field;
          return `${fieldLabel}: ${msgs.join(", ")}`;
        })
        .join("; ");
      showError(`Ошибка в товаре: ${errorMessages}`);
      return;
    }

    if (editingEntryIndex !== null) {
      const newEntries = [...entries];
      newEntries[editingEntryIndex] = currentEntry;
      setEntries(newEntries);
    } else {
      setEntries([...entries, currentEntry]);
    }

    setShowEntryForm(false);
    setCurrentEntry({});
    success(editingEntryIndex !== null ? "Товар обновлён" : "Товар добавлен");
  };

  const handleEditEntry = (index: number) => {
    setCurrentEntry(entries[index]);
    setEditingEntryIndex(index);
    setShowEntryForm(true);
  };

  const handleDeleteEntry = (index: number) => {
    setEntries(entries.filter((_, i) => i !== index));
  };

  const handleCancelEntryForm = () => {
    setShowEntryForm(false);
    setCurrentEntry({});
    setEditingEntryIndex(null);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Основная информация */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          Основная информация
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Название */}
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Название документа *
            </label>
            <input
              type="text"
              {...register("foreignName")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.foreignName
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Введите название"
            />
            <ErrorMessage error={errors.foreignName} />
          </div>

          {/* ИНН контрагента */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              ИНН контрагента
            </label>
            <input
              type="text"
              {...register("contractorTin")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.contractorTin
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Введите ИНН"
            />
            <ErrorMessage error={errors.contractorTin} />
          </div>

          {/* Валюта */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Код валюты
            </label>
            <input
              type="text"
              {...register("currencyCode")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.currencyCode
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="USD, EUR, KGS"
              maxLength={3}
            />
            <ErrorMessage error={errors.currencyCode} />
          </div>

          {/* Сумма */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Сумма
            </label>
            <input
              type="number"
              step="0.01"
              {...register("totalCurrencyValue")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.totalCurrencyValue
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.00"
            />
            <ErrorMessage error={errors.totalCurrencyValue} />
          </div>

          {/* Сумма без налогов */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Сумма без налогов
            </label>
            <input
              type="number"
              step="0.01"
              {...register("totalCurrencyValueWithoutTaxes")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.totalCurrencyValueWithoutTaxes
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.00"
            />
            <ErrorMessage error={errors.totalCurrencyValueWithoutTaxes} />
          </div>

          {/* Курс валюты */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Курс валюты
            </label>
            <input
              type="number"
              step="0.0001"
              {...register("currencyRate")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.currencyRate
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.0000"
            />
            <ErrorMessage error={errors.currencyRate} />
          </div>

          {/* Дата доставки */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Дата доставки
            </label>
            <input
              type="date"
              {...register("deliveryDate")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.deliveryDate
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
            />
            <ErrorMessage error={errors.deliveryDate} />
          </div>

          {/* Дата контракта */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Дата контракта
            </label>
            <input
              type="date"
              {...register("contractStartDate")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.contractStartDate
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
            />
            <ErrorMessage error={errors.contractStartDate} />
          </div>

          {/* Номер контракта */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Номер контракта
            </label>
            <input
              type="text"
              {...register("supplyContractNumber")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.supplyContractNumber
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Введите номер"
            />
            <ErrorMessage error={errors.supplyContractNumber} />
          </div>

          {/* Резидентность */}
          <div className="flex items-center">
            <input
              type="checkbox"
              {...register("isResident")}
              className="w-4 h-4 border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
            />
            <label className="ml-2 text-sm font-medium text-gray-700">
              Резидент КР
            </label>
          </div>

          {/* Данные отправлены */}
          <div className="flex items-center">
            <input
              type="checkbox"
              {...register("isBranchDataSent")}
              className="w-4 h-4 border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
            />
            <label className="ml-2 text-sm font-medium text-gray-700">
              Данные отправлены
            </label>
          </div>

          {/* Цена без налогов */}
          <div className="flex items-center">
            <input
              type="checkbox"
              {...register("isPriceWithoutTaxes")}
              className="w-4 h-4 border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
            />
            <label className="ml-2 text-sm font-medium text-gray-700">
              Цена без налогов
            </label>
          </div>
        </div>
      </div>

      {/* Коды и классификация */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          Коды и классификация
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Код типа доставки
            </label>
            <input
              type="text"
              {...register("deliveryTypeCode")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.deliveryTypeCode
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Код"
            />
            <ErrorMessage error={errors.deliveryTypeCode} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Код типа платежа
            </label>
            <input
              type="text"
              {...register("paymentCode")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.paymentCode
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Код"
            />
            <ErrorMessage error={errors.paymentCode} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Код типа операции
            </label>
            <input
              type="text"
              {...register("operationTypeCode")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.operationTypeCode
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Код"
            />
            <ErrorMessage error={errors.operationTypeCode} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Код страны
            </label>
            <input
              type="text"
              {...register("countryCode")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.countryCode
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Код"
              maxLength={2}
            />
            <ErrorMessage error={errors.countryCode} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              ТИН филиала
            </label>
            <input
              type="text"
              {...register("affiliateTin")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.affiliateTin
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Введите ТИН"
            />
            <ErrorMessage error={errors.affiliateTin} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Код налога на НДС
            </label>
            <input
              type="text"
              {...register("taxRateVATCode")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.taxRateVATCode
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Код"
            />
            <ErrorMessage error={errors.taxRateVATCode} />
          </div>
        </div>
      </div>

      {/* Финансовые данные */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          Финансовые данные
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Оплачено
            </label>
            <input
              type="number"
              step="0.01"
              {...register("paidAmount")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.paidAmount
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.00"
            />
            <ErrorMessage error={errors.paidAmount} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              К оплате
            </label>
            <input
              type="number"
              step="0.01"
              {...register("amountToBePaid")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.amountToBePaid
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.00"
            />
            <ErrorMessage error={errors.amountToBePaid} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Штрафы
            </label>
            <input
              type="number"
              step="0.01"
              {...register("penaltiesAmount")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.penaltiesAmount
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.00"
            />
            <ErrorMessage error={errors.penaltiesAmount} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Штрафные санкции
            </label>
            <input
              type="number"
              step="0.01"
              {...register("finesAmount")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.finesAmount
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.00"
            />
            <ErrorMessage error={errors.finesAmount} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Входящие остатки
            </label>
            <input
              type="number"
              step="0.01"
              {...register("openingBalances")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.openingBalances
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.00"
            />
            <ErrorMessage error={errors.openingBalances} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Исходящие остатки
            </label>
            <input
              type="number"
              step="0.01"
              {...register("closingBalances")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.closingBalances
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.00"
            />
            <ErrorMessage error={errors.closingBalances} />
          </div>
        </div>
      </div>

      {/* Дополнительные поля */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          Дополнительная информация
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Банковский счет поставщика
            </label>
            <input
              type="text"
              {...register("supplierBankAccount")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.supplierBankAccount
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Счет"
            />
            <ErrorMessage error={errors.supplierBankAccount} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Банковский счет контрагента
            </label>
            <input
              type="text"
              {...register("contractorBankAccount")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.contractorBankAccount
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Счет"
            />
            <ErrorMessage error={errors.contractorBankAccount} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Оценённые взносы
            </label>
            <input
              type="number"
              step="0.01"
              {...register("assessedContributionsAmount")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.assessedContributionsAmount
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="0.00"
            />
            <ErrorMessage error={errors.assessedContributionsAmount} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Лицевой счёт
            </label>
            <input
              type="text"
              {...register("personalAccountNumber")}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
                errors.personalAccountNumber
                  ? "border-red-500 bg-red-50"
                  : "border-gray-300 bg-white"
              }`}
              placeholder="Номер счёта"
            />
            <ErrorMessage error={errors.personalAccountNumber} />
          </div>

          <div className="flex items-center">
            <input
              type="checkbox"
              {...register("isIndustry")}
              className="w-4 h-4 border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
            />
            <label className="ml-2 text-sm font-medium text-gray-700">
              Промышленность
            </label>
          </div>
        </div>
      </div>

      {/* Товары и услуги */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold text-gray-900">
            Товары и услуги ({entries.length})
          </h2>
          <button
            type="button"
            onClick={handleAddEntry}
            className="px-3 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors flex items-center gap-2 text-sm"
          >
            <Plus className="w-4 h-4" />
            Добавить товар
          </button>
        </div>

        {showEntryForm && (
          <div className="border-2 border-blue-300 rounded-lg p-4 mb-4 bg-blue-50">
            <h3 className="font-semibold text-gray-900 mb-3">
              {editingEntryIndex !== null ? "Редактировать товар" : "Новый товар"}
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              <input
                type="text"
                placeholder="Код товара"
                value={currentEntry.catalogCode || ""}
                onChange={(e) =>
                  setCurrentEntry({ ...currentEntry, catalogCode: e.target.value })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
              <input
                type="text"
                placeholder="Название товара"
                value={currentEntry.catalogName || ""}
                onChange={(e) =>
                  setCurrentEntry({ ...currentEntry, catalogName: e.target.value })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
              <input
                type="text"
                placeholder="Код единицы измерения"
                value={currentEntry.measureCode || ""}
                onChange={(e) =>
                  setCurrentEntry({ ...currentEntry, measureCode: e.target.value })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
              <input
                type="number"
                step="0.01"
                placeholder="Количество"
                value={currentEntry.quantity || ""}
                onChange={(e) =>
                  setCurrentEntry({
                    ...currentEntry,
                    quantity: parseFloat(e.target.value) || 0,
                  })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
              <input
                type="number"
                step="0.01"
                placeholder="Цена"
                value={currentEntry.price || ""}
                onChange={(e) =>
                  setCurrentEntry({
                    ...currentEntry,
                    price: parseFloat(e.target.value) || 0,
                  })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
              <input
                type="number"
                step="0.01"
                placeholder="Размер оборота"
                value={currentEntry.turnoverSize || ""}
                onChange={(e) =>
                  setCurrentEntry({
                    ...currentEntry,
                    turnoverSize: parseFloat(e.target.value) || 0,
                  })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
              <input
                type="number"
                step="0.01"
                placeholder="Налоговая ставка НДС"
                value={currentEntry.vatRate || ""}
                onChange={(e) =>
                  setCurrentEntry({
                    ...currentEntry,
                    vatRate: parseFloat(e.target.value) || 0,
                  })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
              <input
                type="number"
                step="0.01"
                placeholder="Сумма НДС"
                value={currentEntry.vatAmount || ""}
                onChange={(e) =>
                  setCurrentEntry({
                    ...currentEntry,
                    vatAmount: parseFloat(e.target.value) || 0,
                  })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
              <input
                type="number"
                step="0.01"
                placeholder="Налоговая ставка акциза"
                value={currentEntry.exciseRate || ""}
                onChange={(e) =>
                  setCurrentEntry({
                    ...currentEntry,
                    exciseRate: parseFloat(e.target.value) || 0,
                  })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
              <input
                type="number"
                step="0.01"
                placeholder="Сумма акциза"
                value={currentEntry.exciseAmount || ""}
                onChange={(e) =>
                  setCurrentEntry({
                    ...currentEntry,
                    exciseAmount: parseFloat(e.target.value) || 0,
                  })
                }
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div className="flex gap-2 justify-end mt-3">
              <button
                type="button"
                onClick={handleCancelEntryForm}
                className="px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Отмена
              </button>
              <button
                type="button"
                onClick={handleSaveEntry}
                className="px-3 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
              >
                {editingEntryIndex !== null ? "Обновить" : "Добавить"}
              </button>
            </div>
          </div>
        )}

        {entries.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="px-4 py-2 text-left text-gray-700 font-semibold">
                    Код
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
                  <th className="px-4 py-2 text-center text-gray-700 font-semibold">
                    Действия
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {entries.map((entry, idx) => (
                  <tr key={idx} className="hover:bg-gray-50">
                    <td className="px-4 py-2 font-mono text-gray-900">
                      {entry.catalogCode}
                    </td>
                    <td className="px-4 py-2 text-gray-900">{entry.catalogName}</td>
                    <td className="px-4 py-2 text-right text-gray-900">
                      {entry.quantity}
                    </td>
                    <td className="px-4 py-2 text-right text-gray-900">
                      {entry.price?.toLocaleString("ru-RU")}
                    </td>
                    <td className="px-4 py-2 text-center">
                      <div className="flex gap-2 justify-center">
                        <button
                          type="button"
                          onClick={() => handleEditEntry(idx)}
                          className="text-blue-600 hover:text-blue-700 transition-colors"
                        >
                          ✏️
                        </button>
                        <button
                          type="button"
                          onClick={() => handleDeleteEntry(idx)}
                          className="text-red-600 hover:text-red-700 transition-colors"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="text-gray-500 text-center py-4">
            {showEntryForm ? "" : "Товары не добавлены"}
          </p>
        )}
      </div>

      {/* Комментарий */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          Комментарий
        </h2>
        <textarea
          {...register("comment")}
          rows={4}
          className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition ${
            errors.comment
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white"
          }`}
          placeholder="Введите комментарий (опционально)"
        />
        <ErrorMessage error={errors.comment} />
      </div>

      {/* Actions */}
      <div className="flex gap-2 justify-end">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors flex items-center gap-2"
          >
            <X className="w-4 h-4" />
            Отмена
          </button>
        )}
        <button
          type="submit"
          disabled={isSubmitting}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2 disabled:opacity-50"
        >
          <Save className="w-4 h-4" />
          {isSubmitting ? "Сохранение..." : isEditing ? "Обновить" : "Создать"}
        </button>
      </div>
    </form>
  );
}
