# hooks ディレクトリ

## 概要
カスタムReactフックを配置するディレクトリです。
共通のロジックを再利用可能な形で管理します。

## 命名規則
- フック名は use で始まる camelCase を使用
- ファイル名は フック名.ts

## 想定されるフック例
```
hooks/
├── useAuth.ts          # 認証関連のロジック
├── useWorkout.ts       # トレーニング記録の操作
├── useBodyPart.ts      # 人体図の部位選択ロジック
├── useTranslation.ts   # 多言語対応
└── useLocalStorage.ts  # ローカルストレージ操作
```

## 例
```typescript
// hooks/useAuth.ts
export const useAuth = () => {
  const [user, setUser] = useState(null)
  const login = async (credentials) => {
    // ログイン処理
  }
  const logout = () => {
    // ログアウト処理
  }
  return { user, login, logout }
}
```