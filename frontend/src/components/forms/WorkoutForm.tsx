import React, { useState } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { workoutSchema, WorkoutFormData } from '@/schemas/workout';
import { useMuscleGroups, useExercises } from '@/hooks/useWorkout';
import { MuscleGroupSelector } from './MuscleGroupSelector';
import { ExerciseSelector } from './ExerciseSelector';
import { WeightInput, RepsInput, SetsInput, NotesInput, DateTimeInput } from './inputs/WorkoutInputs';

interface WorkoutFormProps {
  onSubmit: (data: WorkoutFormData) => Promise<void>;
  initialData?: Partial<WorkoutFormData>;
  isLoading?: boolean;
  submitText?: string;
  showCancel?: boolean;
  onCancel?: () => void;
}

export const WorkoutForm: React.FC<WorkoutFormProps> = ({
  onSubmit,
  initialData,
  isLoading = false,
  submitText = '保存',
  showCancel = false,
  onCancel
}) => {
  const [showCustomExerciseModal, setShowCustomExerciseModal] = useState(false);
  
  const {
    control,
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isValid }
  } = useForm<WorkoutFormData>({
    resolver: zodResolver(workoutSchema),
    defaultValues: {
      muscle_group: '',
      exercise_name: '',
      exercise_icon: '',
      weight_kg: undefined,
      reps: undefined,
      sets: undefined,
      notes: '',
      performed_at: new Date(),
      ...initialData
    },
    mode: 'onChange'
  });

  const selectedMuscleGroup = watch('muscle_group');
  const { muscleGroups } = useMuscleGroups();
  const { exercises } = useExercises(selectedMuscleGroup);

  const handleFormSubmit = async (data: WorkoutFormData) => {
    try {
      await onSubmit(data);
    } catch (error) {
      // Error handling is done in the parent component
      console.error('Form submission error:', error);
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-6">
        {/* Muscle Group Selection */}
        <Controller
          name="muscle_group"
          control={control}
          render={({ field: { value, onChange }, fieldState: { error } }) => (
            <MuscleGroupSelector
              value={value}
              onChange={(newValue) => {
                onChange(newValue);
                // Reset exercise when muscle group changes
                if (watch('exercise_name')) {
                  setValue('exercise_name', '');
                }
              }}
              options={muscleGroups}
              error={error}
              disabled={isLoading}
            />
          )}
        />

        {/* Exercise Selection */}
        <Controller
          name="exercise_name"
          control={control}
          render={({ field: { value, onChange }, fieldState: { error } }) => (
            <ExerciseSelector
              value={value}
              onChange={onChange}
              exercises={exercises}
              muscleGroup={selectedMuscleGroup}
              error={error}
              disabled={isLoading}
              onAddCustom={() => setShowCustomExerciseModal(true)}
            />
          )}
        />

        {/* Workout Details Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Controller
            name="weight_kg"
            control={control}
            render={({ field: { value, onChange }, fieldState: { error } }) => (
              <WeightInput
                value={value}
                onChange={(e) => onChange(e.target.value ? Number(e.target.value) : undefined)}
                error={error}
                disabled={isLoading}
              />
            )}
          />

          <Controller
            name="reps"
            control={control}
            render={({ field: { value, onChange }, fieldState: { error } }) => (
              <RepsInput
                value={value}
                onChange={(e) => onChange(e.target.value ? Number(e.target.value) : undefined)}
                error={error}
                disabled={isLoading}
              />
            )}
          />

          <Controller
            name="sets"
            control={control}
            render={({ field: { value, onChange }, fieldState: { error } }) => (
              <SetsInput
                value={value}
                onChange={(e) => onChange(e.target.value ? Number(e.target.value) : undefined)}
                error={error}
                disabled={isLoading}
              />
            )}
          />
        </div>

        {/* Notes */}
        <Controller
          name="notes"
          control={control}
          render={({ field: { value, onChange }, fieldState: { error } }) => (
            <NotesInput
              value={value}
              onChange={onChange}
              error={error}
              disabled={isLoading}
            />
          )}
        />

        {/* Date and Time */}
        <Controller
          name="performed_at"
          control={control}
          render={({ field: { value, onChange }, fieldState: { error } }) => (
            <DateTimeInput
              value={value}
              onChange={(e) => onChange(new Date(e.target.value))}
              error={error}
              disabled={isLoading}
            />
          )}
        />

        {/* Form Actions */}
        <div className="flex flex-col sm:flex-row gap-3 pt-6">
          <button
            type="submit"
            disabled={isLoading || !isValid}
            className={`
              flex-1 px-6 py-3 rounded-lg font-medium transition-all duration-200
              ${isLoading || !isValid
                ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                : 'bg-blue-600 text-white hover:bg-blue-700 active:bg-blue-800 transform active:scale-95'
              }
            `}
          >
            {isLoading ? (
              <div className="flex items-center justify-center">
                <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                保存中...
              </div>
            ) : (
              submitText
            )}
          </button>

          {showCancel && onCancel && (
            <button
              type="button"
              onClick={onCancel}
              disabled={isLoading}
              className="px-6 py-3 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              キャンセル
            </button>
          )}
        </div>
      </form>

      {/* Custom Exercise Modal */}
      {showCustomExerciseModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-medium mb-4">カスタムエクササイズ作成</h3>
            <p className="text-gray-600 mb-4">
              この機能は今後のアップデートで実装予定です。
            </p>
            <div className="flex justify-end">
              <button
                onClick={() => setShowCustomExerciseModal(false)}
                className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
              >
                閉じる
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};