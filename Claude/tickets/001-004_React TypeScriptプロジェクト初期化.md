# チケット001-004: React TypeScriptプロジェクト初期化

## 概要
React 18 + TypeScriptの初期プロジェクトをセットアップし、必要な依存関係をインストールする

## 優先度
高

## 親チケット
001_プロジェクト初期セットアップ

## 詳細タスク

### 1. Reactプロジェクトの作成
- [ ] Viteを使用したプロジェクト作成
  ```bash
  cd frontend
  npm create vite@latest . -- --template react-ts
  ```
- [ ] 不要なデフォルトファイルの削除
- [ ] ディレクトリ構造の調整

### 2. 必要な依存関係のインストール
- [ ] 基本パッケージ
  ```json
  {
    "dependencies": {
      "react": "^18.2.0",
      "react-dom": "^18.2.0",
      "react-router-dom": "^6.20.0",
      "axios": "^1.6.0",
      "zustand": "^4.4.0"
    }
  }
  ```
- [ ] UI/スタイリング関連
  ```json
  {
    "dependencies": {
      "tailwindcss": "^3.3.0",
      "framer-motion": "^10.16.0",
      "@heroicons/react": "^2.0.0"
    }
  }
  ```
- [ ] フォーム・バリデーション
  ```json
  {
    "dependencies": {
      "react-hook-form": "^7.48.0",
      "zod": "^3.22.0"
    }
  }
  ```
- [ ] データ可視化
  ```json
  {
    "dependencies": {
      "d3": "^7.8.0",
      "@types/d3": "^7.4.0"
    }
  }
  ```

### 3. TypeScript設定
- [ ] tsconfig.jsonの最適化
  ```json
  {
    "compilerOptions": {
      "target": "ES2020",
      "module": "ESNext",
      "jsx": "react-jsx",
      "strict": true,
      "esModuleInterop": true,
      "skipLibCheck": true,
      "forceConsistentCasingInFileNames": true,
      "baseUrl": ".",
      "paths": {
        "@/*": ["src/*"]
      }
    }
  }
  ```
- [ ] 型定義ファイルの準備

### 4. 開発環境設定
- [ ] .env.development作成
- [ ] Vite設定のカスタマイズ
- [ ] プロキシ設定（API通信用）

### 5. 基本コンポーネントの作成
- [ ] App.tsx の再構築
- [ ] ルーティング設定
- [ ] エラーバウンダリーの実装

## 受け入れ条件
- `npm run dev`で開発サーバーが起動すること
- TypeScriptの型チェックが正しく動作すること
- 基本的なルーティングが機能すること
- ホットリロードが動作すること

## 見積もり工数
3時間