import React, { memo } from "react";
import { FieldError } from "react-hook-form";
import { AlertCircle } from "lucide-react";

interface ErrorMessageProps {
  error?: FieldError;
  className?: string;
}

const ErrorMessageComponent = ({ error, className = "" }: ErrorMessageProps) => {
  if (!error) return null;

  return (
    <div className={`flex items-center gap-2 text-red-600 text-sm mt-1 ${className}`}>
      <AlertCircle size={16} className="flex-shrink-0" />
      <span>{error.message}</span>
    </div>
  );
};

export const ErrorMessage = memo(ErrorMessageComponent);

interface FormFieldProps {
  label: string;
  error?: FieldError;
  required?: boolean;
  children: React.ReactNode;
  hint?: string;
}

const FormFieldComponent = ({
  label,
  error,
  required = false,
  children,
  hint,
}: FormFieldProps) => {
  return (
    <div className="flex flex-col gap-1">
      <label className="text-sm font-medium text-gray-700">
        {label}
        {required && <span className="text-red-600 ml-1">*</span>}
      </label>
      {children}
      {hint && <p className="text-xs text-gray-500 mt-1">{hint}</p>}
      <ErrorMessage error={error} />
    </div>
  );
};

export const FormField = memo(FormFieldComponent);

interface TextInputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: FieldError;
  hint?: string;
}

export const TextInput = React.forwardRef<HTMLInputElement, TextInputProps>(
  ({ label, error, hint, className = "", ...props }, ref) => {
    return (
      <FormField label={label} error={error} hint={hint} required={props.required}>
        <input
          ref={ref}
          className={`px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
            error
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white hover:border-gray-400"
          } ${className}`}
          {...props}
        />
      </FormField>
    );
  }
);

TextInput.displayName = "TextInput";

interface TextAreaProps
  extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label: string;
  error?: FieldError;
  hint?: string;
}

export const TextArea = React.forwardRef<HTMLTextAreaElement, TextAreaProps>(
  ({ label, error, hint, className = "", ...props }, ref) => {
    return (
      <FormField label={label} error={error} hint={hint} required={props.required}>
        <textarea
          ref={ref}
          className={`px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition resize-none ${
            error
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white hover:border-gray-400"
          } ${className}`}
          {...props}
        />
      </FormField>
    );
  }
);

TextArea.displayName = "TextArea";

interface SelectProps
  extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label: string;
  error?: FieldError;
  options: Array<{ value: string; label: string }>;
  hint?: string;
}

export const Select = React.forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, error, options, hint, className = "", ...props }, ref) => {
    return (
      <FormField label={label} error={error} hint={hint} required={props.required}>
        <select
          ref={ref}
          className={`px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
            error
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white hover:border-gray-400"
          } ${className}`}
          {...props}
        >
          <option value="">Выберите...</option>
          {options.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      </FormField>
    );
  }
);

Select.displayName = "Select";

interface CheckboxProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: FieldError;
}

export const Checkbox = React.forwardRef<HTMLInputElement, CheckboxProps>(
  ({ label, error, className = "", ...props }, ref) => {
    return (
      <FormField label={label} error={error}>
        <div className="flex items-center gap-2">
          <input
            ref={ref}
            type="checkbox"
            className={`w-4 h-4 border rounded cursor-pointer ${
              error
                ? "border-red-500 accent-red-600"
                : "border-gray-300 accent-blue-600"
            } ${className}`}
            {...props}
          />
        </div>
      </FormField>
    );
  }
);

Checkbox.displayName = "Checkbox";

interface NumberInputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: FieldError;
  hint?: string;
}

export const NumberInput = React.forwardRef<HTMLInputElement, NumberInputProps>(
  ({ label, error, hint, className = "", ...props }, ref) => {
    return (
      <FormField label={label} error={error} hint={hint} required={props.required}>
        <input
          ref={ref}
          type="number"
          step="any"
          className={`px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition ${
            error
              ? "border-red-500 bg-red-50"
              : "border-gray-300 bg-white hover:border-gray-400"
          } ${className}`}
          {...props}
        />
      </FormField>
    );
  }
);

NumberInput.displayName = "NumberInput";
