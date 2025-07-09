# チケット001-005: Tailwind CSS設定

## 概要
Tailwind CSSをプロジェクトに導入し、カスタムテーマとユーティリティクラスを設定する

## 優先度
高

## 親チケット
001_プロジェクト初期セットアップ

## 詳細タスク

### 1. Tailwind CSSのインストールと初期設定
- [ ] 必要なパッケージのインストール
  ```bash
  npm install -D tailwindcss postcss autoprefixer
  npx tailwindcss init -p
  ```
- [ ] CSS ファイルの設定
  ```css
  /* src/index.css */
  @tailwind base;
  @tailwind components;
  @tailwind utilities;
  ```

### 2. tailwind.config.jsのカスタマイズ
- [ ] コンテンツパスの設定
  ```javascript
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ]
  ```
- [ ] カスタムカラーパレット
  ```javascript
  colors: {
    primary: {
      50: '#eff6ff',
      500: '#3b82f6',
      900: '#1e3a8a'
    },
    muscle: {
      chest: '#ff6b6b',
      back: '#4ecdc4',
      legs: '#45b7d1',
      arms: '#96ceb4',
      shoulders: '#dda0dd',
      core: '#ffd93d'
    }
  }
  ```
- [ ] カスタムフォント設定
- [ ] ブレークポイントのカスタマイズ

### 3. ユーティリティクラスの拡張
- [ ] カスタムアニメーション
  ```javascript
  animation: {
    'fade-in': 'fadeIn 0.5s ease-in-out',
    'slide-up': 'slideUp 0.3s ease-out'
  }
  ```
- [ ] カスタムスペーシング
- [ ] グラデーション設定

### 4. コンポーネントクラスの定義
- [ ] ボタンスタイル
  ```css
  @layer components {
    .btn-primary {
      @apply px-4 py-2 bg-primary-500 text-white rounded-lg 
             hover:bg-primary-600 transition-colors duration-200;
    }
  }
  ```
- [ ] フォーム要素スタイル
- [ ] カードコンポーネント

### 5. アクセシビリティ対応
- [ ] フォーカススタイルの設定
- [ ] 色覚多様性対応のカラー調整
- [ ] ダークモード対応準備

## 受け入れ条件
- Tailwind CSSが正しく動作すること
- カスタムテーマが適用されること
- パージ設定により本番ビルドサイズが最適化されること
- VSCodeでのインテリセンスが動作すること

## 見積もり工数
2時間