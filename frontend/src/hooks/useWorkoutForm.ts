import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useWorkoutStore } from '@/stores/workoutStore';
import { WorkoutFormData } from '@/schemas/workout';
import { CreateWorkoutRequest, UpdateWorkoutRequest } from '@/types/workout';

interface UseWorkoutFormOptions {
  workoutId?: string;
  onSuccess?: () => void;
  onError?: (error: Error) => void;
  redirectOnSuccess?: string;
}

export const useWorkoutForm = ({
  workoutId,
  onSuccess,
  onError,
  redirectOnSuccess
}: UseWorkoutFormOptions = {}) => {
  const navigate = useNavigate();
  const { addWorkout, updateWorkout, isSubmitting } = useWorkoutStore();
  const [localError, setLocalError] = useState<string | null>(null);

  const handleSubmit = async (data: WorkoutFormData) => {
    setLocalError(null);
    
    try {
      if (workoutId) {
        // Update existing workout
        const updateData: UpdateWorkoutRequest = {
          muscle_group: data.muscle_group,
          exercise_name: data.exercise_name,
          exercise_icon: data.exercise_icon,
          weight_kg: data.weight_kg,
          reps: data.reps,
          sets: data.sets,
          notes: data.notes,
          performed_at: data.performed_at.toISOString()
        };
        
        await updateWorkout(workoutId, updateData);
      } else {
        // Create new workout
        const createData: CreateWorkoutRequest = {
          muscle_group: data.muscle_group,
          exercise_name: data.exercise_name,
          exercise_icon: data.exercise_icon,
          weight_kg: data.weight_kg,
          reps: data.reps,
          sets: data.sets,
          notes: data.notes,
          performed_at: data.performed_at.toISOString()
        };
        
        await addWorkout(createData);
      }

      // Success callback
      if (onSuccess) {
        onSuccess();
      }

      // Redirect if specified
      if (redirectOnSuccess) {
        navigate(redirectOnSuccess);
      }

    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'ワークアウトの保存に失敗しました';
      setLocalError(errorMessage);
      
      if (onError) {
        onError(error instanceof Error ? error : new Error(errorMessage));
      }
    }
  };

  const clearError = () => {
    setLocalError(null);
  };

  return {
    handleSubmit,
    isLoading: isSubmitting,
    error: localError,
    clearError
  };
};