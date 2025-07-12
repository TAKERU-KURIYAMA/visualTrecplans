import React, { useState, useRef, useEffect } from 'react';
import { ChevronDownIcon } from '@heroicons/react/24/outline';
import { MuscleGroup } from '@/types/workout';

interface MuscleGroupSelectorProps {
  value?: string;
  onChange: (value: string) => void;
  options: MuscleGroup[];
  error?: { message?: string };
  disabled?: boolean;
  placeholder?: string;
}

export const MuscleGroupSelector: React.FC<MuscleGroupSelectorProps> = ({
  value,
  onChange,
  options,
  error,
  disabled = false,
  placeholder = '筋肉部位を選択'
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Group options by category
  const groupedOptions = options.reduce((acc, option) => {
    const category = option.category;
    if (!acc[category]) {
      acc[category] = [];
    }
    acc[category].push(option);
    return acc;
  }, {} as Record<string, MuscleGroup[]>);

  // Category display names
  const categoryNames: Record<string, string> = {
    upper: '上半身',
    lower: '下半身',
    core: '体幹',
    full_body: '全身'
  };

  // Find selected option
  const selectedOption = options.find(option => option.code === value);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <div className="relative" ref={dropdownRef}>
      <label className="block text-sm font-medium text-gray-700 mb-1">
        筋肉部位 <span className="text-red-500">*</span>
      </label>
      
      <button
        type="button"
        onClick={() => !disabled && setIsOpen(!isOpen)}
        disabled={disabled}
        className={`
          w-full flex items-center justify-between p-3 border rounded-lg
          transition-colors duration-200
          ${disabled 
            ? 'bg-gray-100 cursor-not-allowed' 
            : isOpen 
              ? 'border-blue-500 ring-2 ring-blue-500 ring-opacity-20' 
              : 'border-gray-300 hover:border-gray-400'
          }
          ${error ? 'border-red-500' : ''}
        `}
      >
        <span className="flex items-center">
          {selectedOption ? (
            <>
              <div 
                className="w-4 h-4 rounded-full mr-3 border border-gray-300"
                style={{ backgroundColor: selectedOption.color_code }}
              />
              <span className="text-gray-900">{selectedOption.name_ja}</span>
            </>
          ) : (
            <span className="text-gray-500">{placeholder}</span>
          )}
        </span>
        <ChevronDownIcon 
          className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${
            isOpen ? 'transform rotate-180' : ''
          }`} 
        />
      </button>
      
      {isOpen && !disabled && (
        <div className="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
          {Object.entries(groupedOptions).map(([category, items]) => (
            <div key={category}>
              <div className="px-4 py-2 text-sm font-medium text-gray-500 bg-gray-50 border-b border-gray-200">
                {categoryNames[category] || category}
              </div>
              {items
                .sort((a, b) => a.sort_order - b.sort_order)
                .map((item) => (
                  <button
                    key={item.code}
                    type="button"
                    onClick={() => {
                      onChange(item.code);
                      setIsOpen(false);
                    }}
                    className={`
                      w-full flex items-center px-4 py-3 text-left transition-colors duration-150
                      hover:bg-blue-50 focus:bg-blue-50 focus:outline-none
                      ${value === item.code ? 'bg-blue-50 text-blue-700' : 'text-gray-900'}
                    `}
                  >
                    <div 
                      className="w-4 h-4 rounded-full mr-3 border border-gray-300"
                      style={{ backgroundColor: item.color_code }}
                    />
                    <span className="flex-1">{item.name_ja}</span>
                    {value === item.code && (
                      <div className="w-2 h-2 bg-blue-600 rounded-full" />
                    )}
                  </button>
                ))}
            </div>
          ))}
          
          {options.length === 0 && (
            <div className="px-4 py-3 text-sm text-gray-500 text-center">
              データが見つかりません
            </div>
          )}
        </div>
      )}
      
      {error && error.message && (
        <p className="mt-1 text-sm text-red-600">{error.message}</p>
      )}
    </div>
  );
};