"use client";

import { memo, useCallback } from "react";
import { AlertCircle } from "lucide-react";

interface ErrorStateProps {
  title?: string;
  message?: string;
  onRetry?: () => void;
}

const ErrorStateComponent = ({
  title = "Ошибка загрузки",
  message = "Произошла ошибка при загрузке данных",
  onRetry,
}: ErrorStateProps) => {
  const handleRetry = useCallback(() => {
    onRetry?.();
  }, [onRetry]);

  return (
    <div
      role="status"
      aria-live="polite"
      aria-label={title}
      className="bg-red-50 border border-red-200 rounded-lg p-6"
    >
      <div className="flex items-start gap-3">
        <AlertCircle
          className="w-6 h-6 text-red-600 flex-shrink-0 mt-0.5"
          aria-hidden="true"
        />
        <div>
          <h3 className="text-lg font-semibold text-red-900 mb-2">{title}</h3>
          <p className="text-red-800 mb-4">{message}</p>
          {onRetry && (
            <button
              onClick={handleRetry}
              aria-label={`Повторить: ${message}`}
              className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors focus:outline-none focus:ring-2 focus:ring-red-600 focus:ring-offset-2"
            >
              Повторить
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export const ErrorState = memo(ErrorStateComponent);
