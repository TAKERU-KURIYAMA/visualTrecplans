export interface MuscleGroup {
  code: string;
  name_ja: string;
  name_en: string;
  category: string;
  color_code: string;
  sort_order: number;
}

export interface Exercise {
  id: number;
  muscle_group_code: string;
  name_ja: string;
  name_en: string;
  icon_name: string;
  is_custom: boolean;
  sort_order: number;
}

export interface ExerciseIcon {
  name: string;
  svg_path: string;
  category: string;
}

export interface Workout {
  id: string;
  muscle_group: string;
  exercise_name: string;
  exercise_icon?: string;
  weight_kg?: number;
  reps?: number;
  sets?: number;
  notes?: string;
  performed_at: string;
  created_at: string;
  updated_at: string;
}

export interface WorkoutFormData {
  muscle_group: string;
  exercise_name: string;
  exercise_icon?: string;
  weight_kg?: number;
  reps?: number;
  sets?: number;
  notes?: string;
  performed_at: Date;
}

export interface CreateWorkoutRequest {
  muscle_group: string;
  exercise_name: string;
  exercise_icon?: string;
  weight_kg?: number;
  reps?: number;
  sets?: number;
  notes?: string;
  performed_at: string;
}

export interface UpdateWorkoutRequest {
  muscle_group?: string;
  exercise_name?: string;
  exercise_icon?: string;
  weight_kg?: number;
  reps?: number;
  sets?: number;
  notes?: string;
  performed_at?: string;
}

export interface WorkoutListResponse {
  workouts: Workout[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface WorkoutFilter {
  muscle_group?: string;
  start_date?: string;
  end_date?: string;
  exercise_name?: string;
  page?: number;
  per_page?: number;
  order_by?: string;
  order?: 'asc' | 'desc';
}

export interface WorkoutStats {
  total_workouts: number;
  total_sets: number;
  total_reps: number;
  total_weight_lifted: number;
  workouts_by_muscle: Record<string, number>;
  most_used_exercises: ExerciseCount[];
  recent_workouts: Workout[];
  weekly_progress: WeeklyProgress[];
}

export interface ExerciseCount {
  exercise_name: string;
  count: number;
}

export interface WeeklyProgress {
  week: string;
  workout_count: number;
  total_sets: number;
  total_reps: number;
  total_weight: number;
}