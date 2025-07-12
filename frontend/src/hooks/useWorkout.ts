import { useEffect } from 'react';
import { useWorkoutStore } from '@/stores/workoutStore';
import { WorkoutFilter } from '@/types/workout';

export const useWorkouts = (filter?: WorkoutFilter) => {
  const {
    workouts,
    total,
    page,
    perPage,
    totalPages,
    isLoading,
    error,
    fetchWorkouts
  } = useWorkoutStore();

  useEffect(() => {
    fetchWorkouts(filter);
  }, [filter, fetchWorkouts]);

  return {
    workouts,
    total,
    page,
    perPage,
    totalPages,
    isLoading,
    error,
    refetch: () => fetchWorkouts(filter)
  };
};

export const useWorkout = (id?: string) => {
  const {
    currentWorkout,
    isLoading,
    error,
    fetchWorkout,
    setCurrentWorkout
  } = useWorkoutStore();

  useEffect(() => {
    if (id) {
      fetchWorkout(id);
    } else {
      setCurrentWorkout(null);
    }

    return () => {
      setCurrentWorkout(null);
    };
  }, [id, fetchWorkout, setCurrentWorkout]);

  return {
    workout: currentWorkout,
    isLoading,
    error
  };
};

export const useWorkoutStats = (period: 'week' | 'month' | 'year' = 'month') => {
  const {
    workoutStats,
    isLoading,
    error,
    fetchWorkoutStats
  } = useWorkoutStore();

  useEffect(() => {
    fetchWorkoutStats(period);
  }, [period, fetchWorkoutStats]);

  return {
    stats: workoutStats,
    isLoading,
    error,
    refetch: () => fetchWorkoutStats(period)
  };
};

export const useMuscleGroups = (lang: string = 'ja', category?: string) => {
  const {
    muscleGroups,
    error,
    fetchMuscleGroups
  } = useWorkoutStore();

  useEffect(() => {
    if (muscleGroups.length === 0 || category) {
      fetchMuscleGroups(lang, category);
    }
  }, [lang, category, muscleGroups.length, fetchMuscleGroups]);

  return {
    muscleGroups,
    error
  };
};

export const useExercises = (muscleGroupCode?: string, lang: string = 'ja') => {
  const {
    exercises,
    customExercises,
    error,
    fetchExercises,
    fetchCustomExercises
  } = useWorkoutStore();

  useEffect(() => {
    fetchExercises(muscleGroupCode, lang);
  }, [muscleGroupCode, lang, fetchExercises]);

  useEffect(() => {
    fetchCustomExercises();
  }, [fetchCustomExercises]);

  // Combine standard and custom exercises
  const allExercises = [
    ...exercises,
    ...customExercises.filter(ce => 
      !muscleGroupCode || ce.muscle_group_code === muscleGroupCode
    )
  ].sort((a, b) => a.sort_order - b.sort_order);

  return {
    exercises: allExercises,
    standardExercises: exercises,
    customExercises,
    error
  };
};

export const useExerciseIcons = (category?: string) => {
  const {
    exerciseIcons,
    error,
    fetchExerciseIcons
  } = useWorkoutStore();

  useEffect(() => {
    if (exerciseIcons.length === 0 || category) {
      fetchExerciseIcons(category);
    }
  }, [category, exerciseIcons.length, fetchExerciseIcons]);

  return {
    exerciseIcons,
    error
  };
};