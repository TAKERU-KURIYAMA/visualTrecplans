import { z } from 'zod';

// ログインフォームスキーマ
export const loginSchema = z.object({
  email: z
    .string()
    .min(1, 'メールアドレスを入力してください')
    .email('有効なメールアドレスを入力してください'),
  password: z
    .string()
    .min(1, 'パスワードを入力してください'),
});

// ユーザー登録フォームスキーマ
export const registerSchema = z.object({
  email: z
    .string()
    .min(1, 'メールアドレスを入力してください')
    .email('有効なメールアドレスを入力してください'),
  password: z
    .string()
    .min(8, 'パスワードは8文字以上で入力してください')
    .regex(
      /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{8,}$/,
      'パスワードは大文字、小文字、数字を含む必要があります'
    ),
  password_confirm: z
    .string()
    .min(1, 'パスワード確認を入力してください'),
  first_name: z
    .string()
    .optional()
    .or(z.literal('')),
  last_name: z
    .string()
    .optional()
    .or(z.literal('')),
}).refine((data) => data.password === data.password_confirm, {
  message: 'パスワードが一致しません',
  path: ['password_confirm'],
});

// パスワード変更フォームスキーマ
export const changePasswordSchema = z.object({
  current_password: z
    .string()
    .min(1, '現在のパスワードを入力してください'),
  new_password: z
    .string()
    .min(8, '新しいパスワードは8文字以上で入力してください')
    .regex(
      /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{8,}$/,
      'パスワードは大文字、小文字、数字を含む必要があります'
    ),
  password_confirm: z
    .string()
    .min(1, 'パスワード確認を入力してください'),
}).refine((data) => data.new_password === data.password_confirm, {
  message: 'パスワードが一致しません',
  path: ['password_confirm'],
});

// プロフィール更新フォームスキーマ
export const updateProfileSchema = z.object({
  first_name: z
    .string()
    .max(100, '名前は100文字以下で入力してください')
    .optional()
    .or(z.literal('')),
  last_name: z
    .string()
    .max(100, '名前は100文字以下で入力してください')
    .optional()
    .or(z.literal('')),
});

// 型エクスポート
export type LoginFormData = z.infer<typeof loginSchema>;
export type RegisterFormData = z.infer<typeof registerSchema>;
export type ChangePasswordFormData = z.infer<typeof changePasswordSchema>;
export type UpdateProfileFormData = z.infer<typeof updateProfileSchema>;