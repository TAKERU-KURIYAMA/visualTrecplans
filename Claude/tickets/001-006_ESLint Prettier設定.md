# チケット001-006: ESLint Prettier設定

## 概要
コード品質と一貫性を保つため、ESLintとPrettierを設定し、自動フォーマット環境を構築する

## 優先度
高

## 親チケット
001_プロジェクト初期セットアップ

## 詳細タスク

### 1. ESLintの設定
- [ ] 必要なパッケージのインストール
  ```bash
  npm install -D eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin
  npm install -D eslint-plugin-react eslint-plugin-react-hooks
  npm install -D eslint-plugin-import eslint-plugin-jsx-a11y
  ```
- [ ] .eslintrc.json の作成
  ```json
  {
    "extends": [
      "eslint:recommended",
      "plugin:@typescript-eslint/recommended",
      "plugin:react/recommended",
      "plugin:react-hooks/recommended",
      "plugin:jsx-a11y/recommended"
    ],
    "parser": "@typescript-eslint/parser",
    "parserOptions": {
      "ecmaVersion": "latest",
      "sourceType": "module",
      "ecmaFeatures": {
        "jsx": true
      }
    },
    "rules": {
      "react/react-in-jsx-scope": "off",
      "@typescript-eslint/no-unused-vars": ["warn", { "argsIgnorePattern": "^_" }],
      "import/order": ["error", {
        "groups": ["builtin", "external", "internal", "parent", "sibling", "index"],
        "newlines-between": "always"
      }]
    }
  }
  ```

### 2. Prettierの設定
- [ ] Prettierのインストール
  ```bash
  npm install -D prettier eslint-config-prettier eslint-plugin-prettier
  ```
- [ ] .prettierrc.json の作成
  ```json
  {
    "semi": false,
    "singleQuote": true,
    "tabWidth": 2,
    "trailingComma": "es5",
    "printWidth": 100,
    "arrowParens": "always",
    "endOfLine": "lf"
  }
  ```
- [ ] .prettierignore の作成

### 3. ESLintとPrettierの統合
- [ ] ESLint設定にPrettierを追加
- [ ] 競合ルールの解決
- [ ] フォーマットスクリプトの追加
  ```json
  "scripts": {
    "lint": "eslint src --ext .ts,.tsx",
    "lint:fix": "eslint src --ext .ts,.tsx --fix",
    "format": "prettier --write \"src/**/*.{ts,tsx,css,md}\""
  }
  ```

### 4. VSCode設定
- [ ] .vscode/settings.json の作成
  ```json
  {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.fixAll.eslint": true
    },
    "eslint.validate": [
      "javascript",
      "javascriptreact",
      "typescript",
      "typescriptreact"
    ]
  }
  ```
- [ ] 推奨拡張機能の設定

### 5. Git フック設定
- [ ] Huskyのインストール
  ```bash
  npm install -D husky lint-staged
  npx husky install
  ```
- [ ] pre-commitフックの設定
- [ ] lint-staged設定

## 受け入れ条件
- ESLintがTypeScriptとReactコードを正しくチェックすること
- Prettierが自動フォーマットを実行すること
- VSCodeで保存時に自動フォーマットされること
- Git commit時にリントチェックが実行されること

## 見積もり工数
2時間