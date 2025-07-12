# types ディレクトリ

## 概要
TypeScriptの型定義ファイルを配置するディレクトリです。
アプリケーション全体で使用する型を一元管理します。

## ファイル構成例
```
types/
├── user.ts         # ユーザー関連の型定義
├── workout.ts      # トレーニング関連の型定義
├── api.ts          # API レスポンスの型定義
├── common.ts       # 共通の型定義
└── index.ts        # 型定義のエクスポート
```

## 命名規則
- ファイル名は camelCase を使用
- 型名は PascalCase を使用
- インターフェースは I プレフィックスは使用しない

## 例
```typescript
// types/user.ts
export interface User {
  id: string
  email: string
  username: string
  createdAt: Date
  preferences: UserPreferences
}

export interface UserPreferences {
  language: 'ja' | 'en' | 'es' | 'fr'
  theme: 'light' | 'dark'
  units: 'metric' | 'imperial'
}
```