import axios, { AxiosInstance, AxiosError } from 'axios';
import type {
  LoginForm,
  RegisterForm,
  LoginResponse,
  RegisterResponse,
  ErrorResponse,
  User,
} from '../types/auth';

class AuthService {
  private api: AxiosInstance;
  private readonly baseURL: string;

  constructor() {
    this.baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
    
    this.api = axios.create({
      baseURL: this.baseURL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor to add auth token
    this.api.interceptors.request.use(
      (config) => {
        const token = this.getAccessToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor to handle token refresh
    this.api.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as any;

        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            const refreshToken = this.getRefreshToken();
            if (refreshToken) {
              const response = await this.refreshToken();
              this.saveTokens(response.access_token, refreshToken);
              originalRequest.headers.Authorization = `Bearer ${response.access_token}`;
              return this.api(originalRequest);
            }
          } catch (refreshError) {
            this.clearTokens();
            window.location.href = '/login';
            return Promise.reject(refreshError);
          }
        }

        return Promise.reject(error);
      }
    );
  }

  // ログイン
  async login(credentials: LoginForm): Promise<LoginResponse> {
    try {
      const response = await this.api.post<LoginResponse>('/auth/login', credentials);
      const { access_token, refresh_token } = response.data;
      
      // トークンを保存
      this.saveTokens(access_token, refresh_token);
      
      return response.data;
    } catch (error) {
      throw this.handleError(error as AxiosError);
    }
  }

  // ユーザー登録
  async register(data: RegisterForm): Promise<RegisterResponse> {
    try {
      const response = await this.api.post<RegisterResponse>('/auth/register', data);
      return response.data;
    } catch (error) {
      throw this.handleError(error as AxiosError);
    }
  }

  // ログアウト
  async logout(): Promise<void> {
    try {
      await this.api.post('/auth/logout');
    } catch (error) {
      console.warn('Logout request failed:', error);
    } finally {
      this.clearTokens();
    }
  }

  // プロフィール取得
  async getProfile(): Promise<User> {
    try {
      const response = await this.api.get<User>('/auth/profile');
      return response.data;
    } catch (error) {
      throw this.handleError(error as AxiosError);
    }
  }

  // プロフィール更新
  async updateProfile(data: Partial<User>): Promise<User> {
    try {
      const response = await this.api.put<User>('/auth/profile', data);
      return response.data;
    } catch (error) {
      throw this.handleError(error as AxiosError);
    }
  }

  // パスワード変更
  async changePassword(data: {
    current_password: string;
    new_password: string;
    password_confirm: string;
  }): Promise<void> {
    try {
      await this.api.put('/auth/password', data);
    } catch (error) {
      throw this.handleError(error as AxiosError);
    }
  }

  // トークンリフレッシュ
  async refreshToken(): Promise<LoginResponse> {
    const refreshToken = this.getRefreshToken();
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    try {
      const response = await this.api.post<LoginResponse>('/auth/refresh', {
        refresh_token: refreshToken,
      });
      return response.data;
    } catch (error) {
      throw this.handleError(error as AxiosError);
    }
  }

  // 認証状態チェック
  async checkAuth(): Promise<User | null> {
    const token = this.getAccessToken();
    if (!token) {
      return null;
    }

    try {
      return await this.getProfile();
    } catch (error) {
      this.clearTokens();
      return null;
    }
  }

  // トークン管理
  private saveTokens(accessToken: string, refreshToken?: string): void {
    localStorage.setItem('auth_access_token', accessToken);
    if (refreshToken) {
      localStorage.setItem('auth_refresh_token', refreshToken);
    }
  }

  private getAccessToken(): string | null {
    return localStorage.getItem('auth_access_token');
  }

  private getRefreshToken(): string | null {
    return localStorage.getItem('auth_refresh_token');
  }

  private clearTokens(): void {
    localStorage.removeItem('auth_access_token');
    localStorage.removeItem('auth_refresh_token');
  }

  // 認証状態確認
  isAuthenticated(): boolean {
    return !!this.getAccessToken();
  }

  // エラーハンドリング
  private handleError(error: AxiosError): Error {
    if (error.response?.data) {
      const errorData = error.response.data as ErrorResponse;
      return new Error(errorData.message || errorData.error || 'リクエストに失敗しました');
    }

    if (error.code === 'ECONNABORTED') {
      return new Error('リクエストがタイムアウトしました');
    }

    if (error.message === 'Network Error') {
      return new Error('ネットワークエラーが発生しました');
    }

    return new Error('予期しないエラーが発生しました');
  }
}

// シングルトンインスタンス
export const authService = new AuthService();
export default authService;