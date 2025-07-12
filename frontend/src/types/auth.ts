export interface User {
  id: string;
  email: string;
  first_name?: string;
  last_name?: string;
  is_active: boolean;
  email_verified: boolean;
  email_verified_at?: string;
  last_login_at?: string;
  login_count: number;
  created_at: string;
}

export interface LoginForm {
  email: string;
  password: string;
}

export interface RegisterForm {
  email: string;
  password: string;
  password_confirm: string;
  first_name?: string;
  last_name?: string;
}

export interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
  message: string;
}

export interface RegisterResponse {
  id: string;
  email: string;
  first_name?: string;
  last_name?: string;
  created_at: string;
  message: string;
}

export interface ErrorResponse {
  error: string;
  message: string;
  code?: string;
  details?: Record<string, any>;
}

export interface ValidationErrorResponse {
  error: string;
  message: string;
  fields: Record<string, string>;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface TokenInfo {
  accessToken: string;
  refreshToken?: string;
  expiresIn: number;
  tokenType: string;
}