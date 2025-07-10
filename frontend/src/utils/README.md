# utils ディレクトリ

## 概要
ユーティリティ関数や共通のヘルパー関数を配置するディレクトリです。
アプリケーション全体で使用する汎用的な関数を管理します。

## ファイル構成例
```
utils/
├── date.ts         # 日付操作関連
├── validation.ts   # バリデーション関数
├── format.ts       # データフォーマット関連
├── constants.ts    # 定数定義
└── helpers.ts      # その他のヘルパー関数
```

## 命名規則
- ファイル名は機能を表す camelCase を使用
- 関数名は動作を表す動詞で始める

## 例
```typescript
// utils/date.ts
export const formatDate = (date: Date, format: string = 'YYYY-MM-DD'): string => {
  // 日付フォーマット処理
}

export const getDaysBetween = (start: Date, end: Date): number => {
  // 日数計算処理
}

// utils/validation.ts
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

export const isStrongPassword = (password: string): boolean => {
  return password.length >= 8 && /[A-Z]/.test(password) && /[0-9]/.test(password)
}
```