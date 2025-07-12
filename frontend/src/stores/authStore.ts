import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User, LoginForm, RegisterForm, AuthState } from '../types/auth';
import { authService } from '../services/auth.service';

interface AuthStore extends AuthState {
  // Actions
  login: (credentials: LoginForm) => Promise<void>;
  register: (data: RegisterForm) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
  updateProfile: (data: Partial<User>) => Promise<void>;
  changePassword: (data: {
    current_password: string;
    new_password: string;
    password_confirm: string;
  }) => Promise<void>;
  clearError: () => void;
  setLoading: (loading: boolean) => void;
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set) => ({
      // Initial state
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      // Login action
      login: async (credentials: LoginForm) => {
        try {
          set({ isLoading: true, error: null });
          
          const response = await authService.login(credentials);
          
          set({
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error: error instanceof Error ? error.message : 'ログインに失敗しました',
          });
          throw error;
        }
      },

      // Register action
      register: async (data: RegisterForm) => {
        try {
          set({ isLoading: true, error: null });
          
          await authService.register(data);
          
          set({
            isLoading: false,
            error: null,
          });
        } catch (error) {
          set({
            isLoading: false,
            error: error instanceof Error ? error.message : '登録に失敗しました',
          });
          throw error;
        }
      },

      // Logout action
      logout: async () => {
        try {
          set({ isLoading: true });
          
          await authService.logout();
          
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          // Even if logout fails on server, clear local state
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error: null,
          });
        }
      },

      // Check authentication status
      checkAuth: async () => {
        try {
          set({ isLoading: true });
          
          const user = await authService.checkAuth();
          
          if (user) {
            set({
              user,
              isAuthenticated: true,
              isLoading: false,
              error: null,
            });
          } else {
            set({
              user: null,
              isAuthenticated: false,
              isLoading: false,
              error: null,
            });
          }
        } catch (error) {
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error: null, // Don't show error for auth check failures
          });
        }
      },

      // Update profile
      updateProfile: async (data: Partial<User>) => {
        try {
          set({ isLoading: true, error: null });
          
          const updatedUser = await authService.updateProfile(data);
          
          set({
            user: updatedUser,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          set({
            isLoading: false,
            error: error instanceof Error ? error.message : 'プロフィールの更新に失敗しました',
          });
          throw error;
        }
      },

      // Change password
      changePassword: async (data: {
        current_password: string;
        new_password: string;
        password_confirm: string;
      }) => {
        try {
          set({ isLoading: true, error: null });
          
          await authService.changePassword(data);
          
          set({
            isLoading: false,
            error: null,
          });
        } catch (error) {
          set({
            isLoading: false,
            error: error instanceof Error ? error.message : 'パスワードの変更に失敗しました',
          });
          throw error;
        }
      },

      // Clear error
      clearError: () => {
        set({ error: null });
      },

      // Set loading state
      setLoading: (loading: boolean) => {
        set({ isLoading: loading });
      },
    }),
    {
      name: 'auth-store',
      // Only persist user authentication state, not loading/error states
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);

// Selectors for convenience
export const useAuth = () => useAuthStore((state) => ({
  user: state.user,
  isAuthenticated: state.isAuthenticated,
  isLoading: state.isLoading,
  error: state.error,
}));

export const useAuthActions = () => useAuthStore((state) => ({
  login: state.login,
  register: state.register,
  logout: state.logout,
  checkAuth: state.checkAuth,
  updateProfile: state.updateProfile,
  changePassword: state.changePassword,
  clearError: state.clearError,
  setLoading: state.setLoading,
}));