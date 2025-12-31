"use client";

import React, { memo } from "react";
import { Loader2 } from "lucide-react";

interface LoadingStateProps {
  fullHeight?: boolean;
  message?: string;
}

/**
 * Компонент состояния загрузки
 * Мемоизирован для оптимизации перерендеров
 */
const LoadingStateComponent: React.FC<LoadingStateProps> = ({
  fullHeight = false,
  message = "Загрузка...",
}) => {
  const containerClass = fullHeight ? "h-64" : "py-12";

  return (
    <div
      className={`flex flex-col items-center justify-center ${containerClass}`}
      role="status"
      aria-live="polite"
      aria-label={message}
    >
      <Loader2
        className="w-12 h-12 text-blue-600 animate-spin mb-4"
        aria-hidden="true"
      />
      <p className="text-gray-600">{message}</p>
    </div>
  );
};

export const LoadingState = memo(LoadingStateComponent);
