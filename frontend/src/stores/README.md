# stores ディレクトリ

## 概要
Zustandを使用した状態管理ストアを配置するディレクトリです。
グローバルな状態管理が必要なデータを管理します。

## ディレクトリ構造
```
stores/
├── authStore.ts      # 認証情報の管理
├── workoutStore.ts   # トレーニングデータの管理
├── uiStore.ts        # UIの状態管理（モーダル、ローディング等）
└── settingsStore.ts  # ユーザー設定の管理
```

## 命名規則
- ストア名は camelCase + Store を使用
- ファイル名は ストア名.ts

## 例
```typescript
// stores/authStore.ts
import { create } from 'zustand'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  login: (user: User) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  login: (user) => set({ user, isAuthenticated: true }),
  logout: () => set({ user: null, isAuthenticated: false }),
}))
```