import { z } from 'zod';

export const workoutSchema = z.object({
  muscle_group: z.string().min(1, '筋肉部位を選択してください'),
  exercise_name: z.string().min(1, 'エクササイズ名を入力してください').max(100, 'エクササイズ名は100文字以内で入力してください'),
  exercise_icon: z.string().optional(),
  weight_kg: z.number().min(0, '重量は0以上で入力してください').max(999.99, '重量は999.99kg以下で入力してください').optional(),
  reps: z.number().int().min(1, '回数は1以上で入力してください').max(999, '回数は999回以下で入力してください').optional(),
  sets: z.number().int().min(1, 'セット数は1以上で入力してください').max(99, 'セット数は99セット以下で入力してください').optional(),
  notes: z.string().max(500, 'メモは500文字以内で入力してください').optional(),
  performed_at: z.date().default(() => new Date())
});

export const workoutFilterSchema = z.object({
  muscle_group: z.string().optional(),
  start_date: z.string().optional(),
  end_date: z.string().optional(),
  exercise_name: z.string().optional(),
  page: z.number().int().min(1).default(1),
  per_page: z.number().int().min(1).max(100).default(20),
  order_by: z.enum(['performed_at', 'created_at', 'exercise_name']).default('performed_at'),
  order: z.enum(['asc', 'desc']).default('desc')
});

export type WorkoutFormData = z.infer<typeof workoutSchema>;
export type WorkoutFilterData = z.infer<typeof workoutFilterSchema>;