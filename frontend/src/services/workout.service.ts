import api from './api';
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

class WorkoutService {
  // Workout CRUD operations
  async createWorkout(data: CreateWorkoutRequest): Promise<Workout> {
    const response = await api.post('/workouts', data);
    return response.data;
  }

  async getWorkouts(filter?: WorkoutFilter): Promise<WorkoutListResponse> {
    const params = new URLSearchParams();
    
    if (filter) {
      Object.entries(filter).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          params.append(key, value.toString());
        }
      });
    }

    const response = await api.get(`/workouts?${params.toString()}`);
    return response.data;
  }

  async getWorkout(id: string): Promise<Workout> {
    const response = await api.get(`/workouts/${id}`);
    return response.data;
  }

  async updateWorkout(id: string, data: UpdateWorkoutRequest): Promise<Workout> {
    const response = await api.put(`/workouts/${id}`, data);
    return response.data;
  }

  async deleteWorkout(id: string): Promise<void> {
    await api.delete(`/workouts/${id}`);
  }

  async getWorkoutStats(period: 'week' | 'month' | 'year' = 'month'): Promise<WorkoutStats> {
    const response = await api.get(`/workouts/stats?period=${period}`);
    return response.data;
  }

  // Master data operations
  async getMuscleGroups(lang: string = 'ja', category?: string): Promise<MuscleGroup[]> {
    const params = new URLSearchParams();
    params.append('lang', lang);
    if (category) {
      params.append('category', category);
    }

    const response = await api.get(`/muscle-groups?${params.toString()}`);
    return response.data.data;
  }

  async getExercises(muscleGroupCode?: string, lang: string = 'ja'): Promise<Exercise[]> {
    const params = new URLSearchParams();
    params.append('lang', lang);
    if (muscleGroupCode) {
      params.append('muscle_group', muscleGroupCode);
    }

    const response = await api.get(`/exercises?${params.toString()}`);
    return response.data.data;
  }

  async getCustomExercises(): Promise<Exercise[]> {
    const response = await api.get('/exercises/custom');
    return response.data.data;
  }

  async createCustomExercise(data: {
    name: string;
    muscle_group_code: string;
    icon_name?: string;
  }): Promise<Exercise> {
    const response = await api.post('/exercises/custom', data);
    return response.data;
  }

  async getExerciseIcons(category?: string): Promise<ExerciseIcon[]> {
    const params = new URLSearchParams();
    if (category) {
      params.append('category', category);
    }

    const response = await api.get(`/exercise-icons?${params.toString()}`);
    return response.data.data;
  }
}

export default new WorkoutService();