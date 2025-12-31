"use client";

import { memo, useCallback } from "react";
import { AlertCircle } from "lucide-react";

interface EmptyStateProps {
  icon?: React.ReactNode;
  title?: string;
  message?: string;
  action?: {
    label: string;
    onClick: () => void;
  };
}

const EmptyStateComponent = ({
  icon = <AlertCircle className="w-12 h-12 text-gray-400" aria-hidden="true" />,
  title = "Нет данных",
  message = "Похоже, здесь ничего нет",
  action,
}: EmptyStateProps) => {
  const handleAction = useCallback(() => {
    action?.onClick();
  }, [action]);

  return (
    <div
      role="status"
      aria-live="polite"
      aria-label={title}
      className="flex flex-col items-center justify-center py-12"
    >
      <div className="mb-4">{icon}</div>
      <h3 className="text-lg font-semibold text-gray-900 mb-2">{title}</h3>
      <p className="text-gray-600 mb-4">{message}</p>
      {action && (
        <button
          onClick={handleAction}
          aria-label={action.label}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2"
        >
          {action.label}
        </button>
      )}
    </div>
  );
};

export const EmptyState = memo(EmptyStateComponent);
