"use client";

import React, { ReactNode } from "react";

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error) => void;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends React.Component<ErrorBoundaryProps,ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error) {
    console.error("ErrorBoundary caught:", error);
    this.props.onError?.(error);
  }

  render() {
    if (this.state.hasError) {
      return (
        this.props.fallback || (
          <div className="bg-red-50 border border-red-200 rounded-lg p-6 m-4">
            <h2 className="text-xl font-bold text-red-900 mb-2">
              Что-то пошло не так
            </h2>
            <p className="text-red-800 mb-4">
              {this.state.error?.message || "Произошла неожиданная ошибка"}
            </p>
            <button
              onClick={() => window.location.reload()}
              className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
            >
              Перезагрузить страницу
            </button>
          </div>
        )
      );
    }

    return this.props.children;
  }
}
