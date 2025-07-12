# services ディレクトリ

## 概要
API通信ロジックや外部サービスとの連携を管理するディレクトリです。
axiosを使用したHTTPクライアントの実装を配置します。

## ファイル構成例
```
services/
├── api.ts          # axiosインスタンスの設定
├── auth.ts         # 認証関連のAPI呼び出し
├── workout.ts      # トレーニング関連のAPI呼び出し
├── supplement.ts   # サプリメント関連のAPI呼び出し
└── export.ts       # データエクスポート関連
```

## 命名規則
- ファイル名は camelCase を使用
- 関数名は動詞で始める（get, post, update, delete等）

## 例
```typescript
// services/api.ts
import axios from 'axios'

export const api = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// services/auth.ts
import { api } from './api'

export const authService = {
  login: async (email: string, password: string) => {
    const response = await api.post('/auth/login', { email, password })
    return response.data
  },
  
  register: async (userData: RegisterData) => {
    const response = await api.post('/auth/register', userData)
    return response.data
  },
}
```