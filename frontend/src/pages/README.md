# pages ディレクトリ

## 概要
ページレベルのコンポーネントを配置するディレクトリです。
React Routerのルートに対応する各ページコンポーネントを管理します。

## ディレクトリ構造例
```
pages/
├── Home/           # ホーム画面
├── Login/          # ログイン画面
├── Register/       # ユーザー登録画面
├── Dashboard/      # ダッシュボード
├── Workout/        # トレーニング記録
├── Progress/       # 進捗確認
└── Settings/       # 設定画面
```

## 命名規則
- ページ名は PascalCase を使用
- 各ページは独自のディレクトリを持つ
- index.tsx をエントリポイントとする

## 例
```typescript
// pages/Dashboard/index.tsx
export const Dashboard: React.FC = () => {
  return (
    <Layout>
      <h1>ダッシュボード</h1>
      {/* ページコンテンツ */}
    </Layout>
  )
}
```