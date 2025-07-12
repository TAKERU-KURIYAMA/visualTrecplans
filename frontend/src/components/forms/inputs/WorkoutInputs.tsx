import React, { forwardRef } from 'react';

interface InputProps {
  value?: number | string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  error?: { message?: string };
  disabled?: boolean;
  placeholder?: string;
  className?: string;
}

export const WeightInput = forwardRef<HTMLInputElement, InputProps>(
  ({ value, onChange, error, disabled = false, placeholder = "0", className = "", ...props }, ref) => {
    return (
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          重量 (kg)
        </label>
        <div className="relative">
          <input
            ref={ref}
            type="number"
            step="0.5"
            min="0"
            max="999.99"
            value={value ?? ''}
            onChange={onChange}
            disabled={disabled}
            placeholder={placeholder}
            className={`
              w-full p-3 pr-12 border rounded-lg transition-colors duration-200
              ${disabled 
                ? 'bg-gray-100 cursor-not-allowed' 
                : 'focus:ring-2 focus:ring-blue-500 focus:border-transparent'
              }
              ${error ? 'border-red-500' : 'border-gray-300'}
              ${className}
            `}
            {...props}
          />
          <span className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 text-sm pointer-events-none">
            kg
          </span>
        </div>
        {error && error.message && (
          <p className="mt-1 text-sm text-red-600">{error.message}</p>
        )}
      </div>
    );
  }
);

WeightInput.displayName = 'WeightInput';

export const RepsInput = forwardRef<HTMLInputElement, InputProps>(
  ({ value, onChange, error, disabled = false, placeholder = "0", className = "", ...props }, ref) => {
    return (
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          回数
        </label>
        <div className="relative">
          <input
            ref={ref}
            type="number"
            min="1"
            max="999"
            value={value ?? ''}
            onChange={onChange}
            disabled={disabled}
            placeholder={placeholder}
            className={`
              w-full p-3 pr-12 border rounded-lg transition-colors duration-200
              ${disabled 
                ? 'bg-gray-100 cursor-not-allowed' 
                : 'focus:ring-2 focus:ring-blue-500 focus:border-transparent'
              }
              ${error ? 'border-red-500' : 'border-gray-300'}
              ${className}
            `}
            {...props}
          />
          <span className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 text-sm pointer-events-none">
            回
          </span>
        </div>
        {error && error.message && (
          <p className="mt-1 text-sm text-red-600">{error.message}</p>
        )}
      </div>
    );
  }
);

RepsInput.displayName = 'RepsInput';

export const SetsInput = forwardRef<HTMLInputElement, InputProps>(
  ({ value, onChange, error, disabled = false, placeholder = "0", className = "", ...props }, ref) => {
    return (
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          セット数
        </label>
        <div className="relative">
          <input
            ref={ref}
            type="number"
            min="1"
            max="99"
            value={value ?? ''}
            onChange={onChange}
            disabled={disabled}
            placeholder={placeholder}
            className={`
              w-full p-3 pr-16 border rounded-lg transition-colors duration-200
              ${disabled 
                ? 'bg-gray-100 cursor-not-allowed' 
                : 'focus:ring-2 focus:ring-blue-500 focus:border-transparent'
              }
              ${error ? 'border-red-500' : 'border-gray-300'}
              ${className}
            `}
            {...props}
          />
          <span className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 text-sm pointer-events-none">
            セット
          </span>
        </div>
        {error && error.message && (
          <p className="mt-1 text-sm text-red-600">{error.message}</p>
        )}
      </div>
    );
  }
);

SetsInput.displayName = 'SetsInput';

interface TextAreaProps {
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
  error?: { message?: string };
  disabled?: boolean;
  placeholder?: string;
  className?: string;
  rows?: number;
}

export const NotesInput = forwardRef<HTMLTextAreaElement, TextAreaProps>(
  ({ value, onChange, error, disabled = false, placeholder = "メモを入力...", className = "", rows = 3, ...props }, ref) => {
    return (
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          メモ
        </label>
        <textarea
          ref={ref}
          value={value ?? ''}
          onChange={onChange}
          disabled={disabled}
          placeholder={placeholder}
          rows={rows}
          maxLength={500}
          className={`
            w-full p-3 border rounded-lg transition-colors duration-200 resize-none
            ${disabled 
              ? 'bg-gray-100 cursor-not-allowed' 
              : 'focus:ring-2 focus:ring-blue-500 focus:border-transparent'
            }
            ${error ? 'border-red-500' : 'border-gray-300'}
            ${className}
          `}
          {...props}
        />
        <div className="flex justify-between items-center mt-1">
          {error && error.message ? (
            <p className="text-sm text-red-600">{error.message}</p>
          ) : (
            <div />
          )}
          <p className="text-sm text-gray-400">
            {(value?.length ?? 0)} / 500
          </p>
        </div>
      </div>
    );
  }
);

NotesInput.displayName = 'NotesInput';

interface DateTimeInputProps {
  value?: string | Date;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  error?: { message?: string };
  disabled?: boolean;
  className?: string;
}

export const DateTimeInput = forwardRef<HTMLInputElement, DateTimeInputProps>(
  ({ value, onChange, error, disabled = false, className = "", ...props }, ref) => {
    // Convert Date to datetime-local string format
    const formatDateTimeLocal = (date: Date | string) => {
      const d = new Date(date);
      const pad = (num: number) => num.toString().padStart(2, '0');
      
      return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
    };

    const inputValue = value 
      ? (value instanceof Date ? formatDateTimeLocal(value) : value)
      : formatDateTimeLocal(new Date());

    return (
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          実施日時
        </label>
        <input
          ref={ref}
          type="datetime-local"
          value={inputValue}
          onChange={onChange}
          disabled={disabled}
          className={`
            w-full p-3 border rounded-lg transition-colors duration-200
            ${disabled 
              ? 'bg-gray-100 cursor-not-allowed' 
              : 'focus:ring-2 focus:ring-blue-500 focus:border-transparent'
            }
            ${error ? 'border-red-500' : 'border-gray-300'}
            ${className}
          `}
          {...props}
        />
        {error && error.message && (
          <p className="mt-1 text-sm text-red-600">{error.message}</p>
        )}
      </div>
    );
  }
);

DateTimeInput.displayName = 'DateTimeInput';