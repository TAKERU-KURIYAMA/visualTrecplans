import React, { useState, useRef, useEffect } from 'react';
import { ChevronDownIcon, PlusIcon } from '@heroicons/react/24/outline';
import { Exercise } from '@/types/workout';

interface ExerciseSelectorProps {
  value?: string;
  onChange: (value: string) => void;
  exercises: Exercise[];
  muscleGroup?: string;
  error?: { message?: string };
  disabled?: boolean;
  placeholder?: string;
  onAddCustom?: () => void;
}

export const ExerciseSelector: React.FC<ExerciseSelectorProps> = ({
  value,
  onChange,
  exercises,
  muscleGroup,
  error,
  disabled = false,
  placeholder = 'エクササイズを選択',
  onAddCustom
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const dropdownRef = useRef<HTMLDivElement>(null);
  const searchRef = useRef<HTMLInputElement>(null);

  // Filter exercises based on muscle group and search term
  const filteredExercises = exercises.filter(exercise => {
    const matchesMuscleGroup = !muscleGroup || exercise.muscle_group_code === muscleGroup;
    const matchesSearch = !searchTerm || 
      exercise.name_ja.toLowerCase().includes(searchTerm.toLowerCase()) ||
      exercise.name_en.toLowerCase().includes(searchTerm.toLowerCase());
    
    return matchesMuscleGroup && matchesSearch;
  });

  // Group exercises by type (standard vs custom)
  const standardExercises = filteredExercises.filter(ex => !ex.is_custom);
  const customExercises = filteredExercises.filter(ex => ex.is_custom);

  // Find selected exercise
  const selectedExercise = exercises.find(exercise => exercise.name_ja === value);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
        setSearchTerm('');
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // Focus search input when dropdown opens
  useEffect(() => {
    if (isOpen && searchRef.current) {
      searchRef.current.focus();
    }
  }, [isOpen]);

  const handleOpen = () => {
    if (!disabled) {
      setIsOpen(!isOpen);
      setSearchTerm('');
    }
  };

  const handleSelect = (exerciseName: string) => {
    onChange(exerciseName);
    setIsOpen(false);
    setSearchTerm('');
  };

  return (
    <div className="relative" ref={dropdownRef}>
      <label className="block text-sm font-medium text-gray-700 mb-1">
        エクササイズ <span className="text-red-500">*</span>
      </label>
      
      <button
        type="button"
        onClick={handleOpen}
        disabled={disabled || !muscleGroup}
        className={`
          w-full flex items-center justify-between p-3 border rounded-lg
          transition-colors duration-200
          ${disabled || !muscleGroup
            ? 'bg-gray-100 cursor-not-allowed' 
            : isOpen 
              ? 'border-blue-500 ring-2 ring-blue-500 ring-opacity-20' 
              : 'border-gray-300 hover:border-gray-400'
          }
          ${error ? 'border-red-500' : ''}
        `}
      >
        <span className="flex items-center">
          {selectedExercise ? (
            <>
              {selectedExercise.icon_name && (
                <div className="w-5 h-5 mr-3 text-gray-500">
                  {/* Icon placeholder - replace with actual icon */}
                  <div className="w-5 h-5 bg-gray-300 rounded" />
                </div>
              )}
              <span className="text-gray-900">{selectedExercise.name_ja}</span>
              {selectedExercise.is_custom && (
                <span className="ml-2 px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded">
                  カスタム
                </span>
              )}
            </>
          ) : (
            <span className="text-gray-500">
              {!muscleGroup ? '先に筋肉部位を選択してください' : placeholder}
            </span>
          )}
        </span>
        <ChevronDownIcon 
          className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${
            isOpen ? 'transform rotate-180' : ''
          }`} 
        />
      </button>
      
      {isOpen && !disabled && muscleGroup && (
        <div className="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-hidden">
          {/* Search input */}
          <div className="p-3 border-b border-gray-200">
            <input
              ref={searchRef}
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="エクササイズを検索..."
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          
          <div className="max-h-40 overflow-y-auto">
            {/* Standard exercises */}
            {standardExercises.length > 0 && (
              <div>
                <div className="px-4 py-2 text-sm font-medium text-gray-500 bg-gray-50">
                  標準エクササイズ
                </div>
                {standardExercises
                  .sort((a, b) => a.sort_order - b.sort_order)
                  .map((exercise) => (
                    <button
                      key={exercise.id}
                      type="button"
                      onClick={() => handleSelect(exercise.name_ja)}
                      className={`
                        w-full flex items-center px-4 py-3 text-left transition-colors duration-150
                        hover:bg-blue-50 focus:bg-blue-50 focus:outline-none
                        ${value === exercise.name_ja ? 'bg-blue-50 text-blue-700' : 'text-gray-900'}
                      `}
                    >
                      {exercise.icon_name && (
                        <div className="w-5 h-5 mr-3 text-gray-400">
                          {/* Icon placeholder */}
                          <div className="w-5 h-5 bg-gray-300 rounded" />
                        </div>
                      )}
                      <span className="flex-1">{exercise.name_ja}</span>
                      {value === exercise.name_ja && (
                        <div className="w-2 h-2 bg-blue-600 rounded-full" />
                      )}
                    </button>
                  ))}
              </div>
            )}
            
            {/* Custom exercises */}
            {customExercises.length > 0 && (
              <div>
                <div className="px-4 py-2 text-sm font-medium text-gray-500 bg-gray-50 border-t border-gray-200">
                  カスタムエクササイズ
                </div>
                {customExercises.map((exercise) => (
                  <button
                    key={exercise.id}
                    type="button"
                    onClick={() => handleSelect(exercise.name_ja)}
                    className={`
                      w-full flex items-center px-4 py-3 text-left transition-colors duration-150
                      hover:bg-blue-50 focus:bg-blue-50 focus:outline-none
                      ${value === exercise.name_ja ? 'bg-blue-50 text-blue-700' : 'text-gray-900'}
                    `}
                  >
                    {exercise.icon_name && (
                      <div className="w-5 h-5 mr-3 text-gray-400">
                        {/* Icon placeholder */}
                        <div className="w-5 h-5 bg-gray-300 rounded" />
                      </div>
                    )}
                    <span className="flex-1">{exercise.name_ja}</span>
                    <span className="ml-2 px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded">
                      カスタム
                    </span>
                    {value === exercise.name_ja && (
                      <div className="w-2 h-2 bg-blue-600 rounded-full ml-2" />
                    )}
                  </button>
                ))}
              </div>
            )}
            
            {/* Add custom exercise option */}
            {onAddCustom && (
              <div className="border-t border-gray-200">
                <button
                  type="button"
                  onClick={() => {
                    onAddCustom();
                    setIsOpen(false);
                  }}
                  className="w-full flex items-center px-4 py-3 text-left text-blue-600 hover:bg-blue-50 focus:bg-blue-50 focus:outline-none transition-colors duration-150"
                >
                  <PlusIcon className="w-5 h-5 mr-3" />
                  <span>カスタムエクササイズを追加</span>
                </button>
              </div>
            )}
            
            {/* No results */}
            {filteredExercises.length === 0 && !onAddCustom && (
              <div className="px-4 py-3 text-sm text-gray-500 text-center">
                エクササイズが見つかりません
              </div>
            )}
          </div>
        </div>
      )}
      
      {error && error.message && (
        <p className="mt-1 text-sm text-red-600">{error.message}</p>
      )}
    </div>
  );
};