import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import workoutService from '@/services/workout.service';
import { 
  Workout, 
  WorkoutListResponse, 
  CreateWorkoutRequest, 
  UpdateWorkoutRequest, 
  WorkoutFilter,
  WorkoutStats,
  MuscleGroup,
  Exercise,
  ExerciseIcon
} from '@/types/workout';

interface WorkoutState {
  // Workout data
  workouts: Workout[];
  currentWorkout: Workout | null;
  workoutStats: WorkoutStats | null;
  total: number;
  page: number;
  perPage: number;
  totalPages: number;

  // Master data
  muscleGroups: MuscleGroup[];
  exercises: Exercise[];
  customExercises: Exercise[];
  exerciseIcons: ExerciseIcon[];

  // Loading states
  isLoading: boolean;
  isSubmitting: boolean;
  error: string | null;

  // Actions
  fetchWorkouts: (filter?: WorkoutFilter) => Promise<void>;
  fetchWorkout: (id: string) => Promise<void>;
  addWorkout: (data: CreateWorkoutRequest) => Promise<void>;
  updateWorkout: (id: string, data: UpdateWorkoutRequest) => Promise<void>;
  deleteWorkout: (id: string) => Promise<void>;
  fetchWorkoutStats: (period?: 'week' | 'month' | 'year') => Promise<void>;

  // Master data actions
  fetchMuscleGroups: (lang?: string, category?: string) => Promise<void>;
  fetchExercises: (muscleGroupCode?: string, lang?: string) => Promise<void>;
  fetchCustomExercises: () => Promise<void>;
  createCustomExercise: (data: { name: string; muscle_group_code: string; icon_name?: string }) => Promise<void>;
  fetchExerciseIcons: (category?: string) => Promise<void>;

  // Utility actions
  clearError: () => void;
  setCurrentWorkout: (workout: Workout | null) => void;
}

export const useWorkoutStore = create<WorkoutState>()(
  persist(
    (set, get) => ({
      // Initial state
      workouts: [],
      currentWorkout: null,
      workoutStats: null,
      total: 0,
      page: 1,
      perPage: 20,
      totalPages: 0,

      muscleGroups: [],
      exercises: [],
      customExercises: [],
      exerciseIcons: [],

      isLoading: false,
      isSubmitting: false,
      error: null,

      // Workout actions
      fetchWorkouts: async (filter?: WorkoutFilter) => {
        set({ isLoading: true, error: null });
        try {
          const response: WorkoutListResponse = await workoutService.getWorkouts(filter);
          set({
            workouts: response.workouts,
            total: response.total,
            page: response.page,
            perPage: response.per_page,
            totalPages: response.total_pages,
            isLoading: false
          });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to fetch workouts',
            isLoading: false 
          });
        }
      },

      fetchWorkout: async (id: string) => {
        set({ isLoading: true, error: null });
        try {
          const workout = await workoutService.getWorkout(id);
          set({ currentWorkout: workout, isLoading: false });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to fetch workout',
            isLoading: false 
          });
        }
      },

      addWorkout: async (data: CreateWorkoutRequest) => {
        set({ isSubmitting: true, error: null });
        try {
          const newWorkout = await workoutService.createWorkout(data);
          set(state => ({
            workouts: [newWorkout, ...state.workouts],
            total: state.total + 1,
            isSubmitting: false
          }));
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to create workout',
            isSubmitting: false 
          });
          throw error;
        }
      },

      updateWorkout: async (id: string, data: UpdateWorkoutRequest) => {
        set({ isSubmitting: true, error: null });
        try {
          const updatedWorkout = await workoutService.updateWorkout(id, data);
          set(state => ({
            workouts: state.workouts.map(w => w.id === id ? updatedWorkout : w),
            currentWorkout: state.currentWorkout?.id === id ? updatedWorkout : state.currentWorkout,
            isSubmitting: false
          }));
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to update workout',
            isSubmitting: false 
          });
          throw error;
        }
      },

      deleteWorkout: async (id: string) => {
        set({ isLoading: true, error: null });
        try {
          await workoutService.deleteWorkout(id);
          set(state => ({
            workouts: state.workouts.filter(w => w.id !== id),
            total: Math.max(0, state.total - 1),
            currentWorkout: state.currentWorkout?.id === id ? null : state.currentWorkout,
            isLoading: false
          }));
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to delete workout',
            isLoading: false 
          });
          throw error;
        }
      },

      fetchWorkoutStats: async (period: 'week' | 'month' | 'year' = 'month') => {
        set({ isLoading: true, error: null });
        try {
          const stats = await workoutService.getWorkoutStats(period);
          set({ workoutStats: stats, isLoading: false });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to fetch workout stats',
            isLoading: false 
          });
        }
      },

      // Master data actions
      fetchMuscleGroups: async (lang: string = 'ja', category?: string) => {
        try {
          const muscleGroups = await workoutService.getMuscleGroups(lang, category);
          set({ muscleGroups });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to fetch muscle groups'
          });
        }
      },

      fetchExercises: async (muscleGroupCode?: string, lang: string = 'ja') => {
        try {
          const exercises = await workoutService.getExercises(muscleGroupCode, lang);
          set({ exercises });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to fetch exercises'
          });
        }
      },

      fetchCustomExercises: async () => {
        try {
          const customExercises = await workoutService.getCustomExercises();
          set({ customExercises });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to fetch custom exercises'
          });
        }
      },

      createCustomExercise: async (data: { name: string; muscle_group_code: string; icon_name?: string }) => {
        set({ isSubmitting: true, error: null });
        try {
          const newExercise = await workoutService.createCustomExercise(data);
          set(state => ({
            customExercises: [...state.customExercises, newExercise],
            exercises: [...state.exercises, newExercise],
            isSubmitting: false
          }));
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to create custom exercise',
            isSubmitting: false 
          });
          throw error;
        }
      },

      fetchExerciseIcons: async (category?: string) => {
        try {
          const exerciseIcons = await workoutService.getExerciseIcons(category);
          set({ exerciseIcons });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to fetch exercise icons'
          });
        }
      },

      // Utility actions
      clearError: () => set({ error: null }),
      setCurrentWorkout: (workout: Workout | null) => set({ currentWorkout: workout })
    }),
    {
      name: 'workout-store',
      partialize: (state) => ({
        muscleGroups: state.muscleGroups,
        exerciseIcons: state.exerciseIcons
      })
    }
  )
);