import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import { ApiError } from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiClient {
  private client: AxiosInstance;
  private isRedirecting = false;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
      timeout: 30000,
    });

    // Request interceptor для добавления токена и организации ID
    this.client.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        const token = this.getToken();
        if (token && config.headers) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        
        const orgId = this.getOrganizationId();
        if (orgId && config.headers) {
          config.headers['X-Org-Id'] = orgId;
        }
        
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor для обработки ошибок
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError<ApiError>) => {
        if (error.response?.status === 401) {
          // Токен истек или невалиден
          console.warn('Unauthorized error - clearing token and redirecting to login');
          this.clearToken();
          
          // Предотвратить множественные редиректы
          if (!this.isRedirecting && typeof window !== 'undefined') {
            this.isRedirecting = true;
            // Используем setTimeout чтобы дать время завершить текущую операцию
            setTimeout(() => {
              window.location.href = '/login';
            }, 100);
          }
        } else if (error.response) {
          // Логировать другие HTTP ошибки
          console.error(
            `HTTP Error ${error.response.status}:`,
            error.response.data || error.message
          );
        } else if (error.request) {
          // Запрос был отправлен, но нет ответа
          console.error('No response received:', error.request);
        } else {
          // Ошибка при подготовке запроса
          console.error('Error preparing request:', error.message);
        }
        return Promise.reject(error);
      }
    );
  }

  private getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('auth_token');
  }

  clearToken(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem('auth_token');
    localStorage.removeItem('org_id');
    localStorage.removeItem('org_token');
  }

  setToken(token: string): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem('auth_token', token);
  }

  // Organization token management
  setOrganizationToken(token: string): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem('org_token', token);
  }

  getOrganizationToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('org_token');
  }

  clearOrganizationToken(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem('org_token');
  }

  // Organization ID management for X-Org-Id header
  setOrganizationId(orgId: string): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem('org_id', orgId);
  }

  getOrganizationId(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('org_id');
  }

  clearOrganizationId(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem('org_id');
  }

  get axiosInstance(): AxiosInstance {
    return this.client;
  }

  // HTTP methods
  async get<T>(url: string, config?: InternalAxiosRequestConfig) {
    return this.client.get<T>(url, config);
  }

  async post<T>(url: string, data?: unknown, config?: InternalAxiosRequestConfig) {
    return this.client.post<T>(url, data, config);
  }

  async put<T>(url: string, data?: unknown, config?: InternalAxiosRequestConfig) {
    return this.client.put<T>(url, data, config);
  }

  async patch<T>(url: string, data?: unknown, config?: InternalAxiosRequestConfig) {
    return this.client.patch<T>(url, data, config);
  }

  async delete<T>(url: string, config?: InternalAxiosRequestConfig) {
    return this.client.delete<T>(url, config);
  }
}

export const apiClient = new ApiClient();
export const api = apiClient.axiosInstance;
