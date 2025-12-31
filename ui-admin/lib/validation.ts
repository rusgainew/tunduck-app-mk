import { z } from "zod";

// Email validation
const emailSchema = z.string().email("Некорректный email");

// Phone validation (basic)
const phoneSchema = z.string().regex(/^\+?[\d\-\s()]+$/, "Некорректный номер телефона");

// TIN validation (numerical)
const tinSchema = z.string().regex(/^\d+$/, "ИНН должен содержать только цифры");

// Currency code validation (3 characters, uppercase)
const currencyCodeSchema = z.string().regex(/^[A-Z]{3}$/, "Код валюты должен быть 3 буквы");

// Date validation
const dateSchema = z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "Дата должна быть в формате YYYY-MM-DD");

// Positive number validation
const positiveNumberSchema = z.union([
  z.number().positive("Число должно быть положительным"),
  z.string().transform(Number).refine(n => n > 0, "Число должно быть положительным"),
]).optional().nullable();

// Non-negative number validation
const nonNegativeNumberSchema = z.union([
  z.number().nonnegative("Число не должно быть отрицательным"),
  z.string().transform(Number).refine(n => n >= 0, "Число не должно быть отрицательным"),
]).optional().nullable();

// Positive number validation - for coercing form input (always strings initially)
const formPositiveNumberSchema = z.union([
  z.number().positive().nullable(),
  z.null(),
]).optional();

// Non-negative number validation - for coercing form input  
const formNonNegativeNumberSchema = z.union([
  z.number().nonnegative().nullable(),
  z.null(),
]).optional();

// Document entry (товар/услуга) schema
export const entrySchema = z.object({
  catalogCode: z.string().min(1, "Код каталога обязателен"),
  catalogName: z.string().min(1, "Название обязательно"),
  measureCode: z.string().min(1, "Код единицы измерения обязателен"),
  quantity: z.union([
    z.number().positive("Количество должно быть больше нуля"),
    z.string().refine(val => val !== "", "Количество обязательно").transform(Number).refine(n => n > 0, "Количество должно быть больше нуля"),
  ]),
  price: z.union([
    z.number().nonnegative("Цена не должна быть отрицательной"),
    z.string().refine(val => val !== "", "Цена обязательна").transform(Number).refine(n => n >= 0, "Цена не должна быть отрицательной"),
  ]),
  turnoverSize: nonNegativeNumberSchema,
  vatRate: nonNegativeNumberSchema,
  vatAmount: nonNegativeNumberSchema,
  exciseRate: nonNegativeNumberSchema,
  exciseAmount: nonNegativeNumberSchema,
});

export type EntryFormData = z.infer<typeof entrySchema>;

// Main document schema with comprehensive validation
export const documentSchema = z.object({
  // Required field
  foreignName: z.string()
    .min(1, "Название документа обязательно")
    .max(255, "Название не должно превышать 255 символов"),

  // Optional fields with custom validation
  contractorTin: tinSchema.optional().nullable(),
  
  currencyCode: currencyCodeSchema.optional().nullable(),
  
  totalCurrencyValue: formPositiveNumberSchema,
  totalCurrencyValueWithoutTaxes: formNonNegativeNumberSchema,
  
  deliveryDate: dateSchema.optional().nullable(),
  deliveryTypeCode: z.string().max(10, "Код типа доставки не должен превышать 10 символов").optional().nullable(),
  
  isResident: z.boolean().optional(),
  
  contractStartDate: dateSchema.optional().nullable(),
  
  comment: z.string().max(1000, "Комментарий не должен превышать 1000 символов").optional().nullable(),
  
  currencyRate: formPositiveNumberSchema,
  
  supplyContractNumber: z.string().max(50, "Номер контракта не должен превышать 50 символов").optional().nullable(),
  
  affiliateTin: tinSchema.optional().nullable(),
  
  isBranchDataSent: z.boolean().optional(),
  isPriceWithoutTaxes: z.boolean().optional(),
  isIndustry: z.boolean().optional(),
  
  ownedCrmReceiptCode: z.string().max(50, "CRM код не должен превышать 50 символов").optional().nullable(),
  operationTypeCode: z.string().max(10, "Код типа операции не должен превышать 10 символов").optional().nullable(),
  
  supplierBankAccount: z.string().regex(/^[A-Z]{2}\d{2}[A-Z0-9]{1,30}$|^.{0}$/, "Некорректный номер банковского счета").optional().nullable(),
  contractorBankAccount: z.string().regex(/^[A-Z]{2}\d{2}[A-Z0-9]{1,30}$|^.{0}$/, "Некорректный номер банковского счета").optional().nullable(),
  
  countryCode: z.string().regex(/^[A-Z]{2}$|^.{0}$/, "Код страны должен быть 2 буквы").optional().nullable(),
  
  paymentCode: z.string().max(10, "Код платежа не должен превышать 10 символов").optional().nullable(),
  taxRateVATCode: z.string().max(10, "Код ставки НДС не должен превышать 10 символов").optional().nullable(),
  deliveryCode: z.string().max(10, "Код доставки не должен превышать 10 символов").optional().nullable(),
  
  paidAmount: formNonNegativeNumberSchema,
  penaltiesAmount: formNonNegativeNumberSchema,
  finesAmount: formNonNegativeNumberSchema,
  
  closingBalances: formNonNegativeNumberSchema,
  openingBalances: formNonNegativeNumberSchema,
  assessedContributionsAmount: formNonNegativeNumberSchema,
  amountToBePaid: formNonNegativeNumberSchema,
  
  personalAccountNumber: z.string().max(50, "Номер счета не должен превышать 50 символов").optional().nullable(),
});

export type DocumentFormData = z.infer<typeof documentSchema>;

// Login form schema
export const loginSchema = z.object({
  username: z.string()
    .min(3, "Username должен содержать минимум 3 символа")
    .max(50, "Username не должен превышать 50 символов"),
  password: z.string()
    .min(6, "Пароль должен содержать минимум 6 символов")
    .max(100, "Пароль не должен превышать 100 символов"),
});

export type LoginFormData = z.infer<typeof loginSchema>;

// Register form schema
export const registerSchema = z.object({
  email: emailSchema,
  password: z.string()
    .min(8, "Пароль должен содержать минимум 8 символов")
    .max(100, "Пароль не должен превышать 100 символов")
    .regex(/[A-Z]/, "Пароль должен содержать заглавную букву")
    .regex(/[a-z]/, "Пароль должен содержать строчную букву")
    .regex(/\d/, "Пароль должен содержать цифру"),
  confirmPassword: z.string(),
  fullName: z.string()
    .min(2, "Имя должно быть минимум 2 символа")
    .max(100, "Имя не должно превышать 100 символов"),
  username: z.string()
    .min(3, "Имя пользователя должно быть минимум 3 символа")
    .max(50, "Имя пользователя не должно превышать 50 символов")
    .regex(/^[a-zA-Z0-9_-]+$/, "Имя пользователя может содержать только буквы, цифры, дефис и подчеркивание"),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Пароли не совпадают",
  path: ["confirmPassword"],
});

export type RegisterFormData = z.infer<typeof registerSchema>;

// Organization form schema
export const organizationSchema = z.object({
  name: z.string()
    .min(1, "Название организации обязательно")
    .max(255, "Название не должно превышать 255 символов"),
  
  description: z.string()
    .max(1000, "Описание не должно превышать 1000 символов")
    .optional(),
  
  email: emailSchema.optional().or(z.literal("")),
  
  phone: phoneSchema.optional().or(z.literal("")),
  
  address: z.string()
    .max(500, "Адрес не должен превышать 500 символов")
    .optional(),
});

export type OrganizationFormData = z.infer<typeof organizationSchema>;

// User form schema
export const userSchema = z.object({
  email: emailSchema,
  username: z.string()
    .min(3, "Имя пользователя должно быть минимум 3 символа")
    .max(50, "Имя пользователя не должно превышать 50 символов")
    .regex(/^[a-zA-Z0-9_-]+$/, "Имя пользователя может содержать только буквы, цифры, дефис и подчеркивание"),
  fullName: z.string()
    .min(2, "Имя должно быть минимум 2 символа")
    .max(100, "Имя не должно превышать 100 символов"),
  phone: phoneSchema.optional().or(z.literal("")),
  role: z.enum(["admin", "user", "viewer"]),
});

export type UserFormData = z.infer<typeof userSchema>;
