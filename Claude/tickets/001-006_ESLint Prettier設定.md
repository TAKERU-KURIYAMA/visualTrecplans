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

## 完了報告
### 作業日時: 2025-07-11

### 実施内容
1. **ESLint設定の構築**
   - **包括的なESLintルール**: TypeScript、React、jsx-a11y、Prettier統合
   - **適切なparser設定**: @typescript-eslint/parserでTypeScript完全対応
   - **プラグイン統合**: 品質とアクセシビリティを重視したルール設定
   - **カスタムルール**: 未使用変数警告、JSX scope設定、prop-types無効化

2. **Prettier設定の最適化**
   - **詳細なフォーマット設定**: セミコロンなし、シングルクォート、LF改行
   - **プロジェクト統一**: 2スペースインデント、100文字幅制限
   - **包括的な.prettierignore**: ビルド成果物、依存関係、環境ファイルを除外
   - **VS Code統合**: 言語別デフォルトフォーマッタ設定

3. **package.jsonスクリプトの拡充**
   - **lint**: ESLintによるコード品質チェック
   - **lint:fix**: 自動修正機能付きリント
   - **format**: Prettierによる一括フォーマット
   - **format:check**: フォーマット確認（CI用）
   - **type-check**: TypeScript型チェック

4. **lint-stagedとHusky統合**
   - **pre-commitフック**: コミット前の自動リント・フォーマット
   - **段階的チェック**: TypeScript/JSXファイルとその他ファイルの分離処理
   - **Git統合**: プロジェクトルートでのHusky設定

5. **GitHub Actions CI/CD**
   - **包括的な品質チェック**: ESLint、Prettier、TypeScript、ビルド確認
   - **Node.js 18環境**: 最新LTS版でのテスト実行
   - **キャッシュ最適化**: npm依存関係の高速化

6. **VSCode設定の強化**
   - **ESLint統合**: workingDirectories設定でフロントエンド特化
   - **自動修正**: 保存時のESLintルール適用
   - **拡張機能**: TypeScript Importerの追加推奨

7. **コード品質の向上**
   - **アクセシビリティ修正**: 不適切なhref="#"をbutton要素に変更
   - **TypeScript最適化**: 未使用React importの削除
   - **フォーマット統一**: 全ファイルでの一貫したスタイル適用

### 技術的特徴
- **完全なTypeScript対応**: 型チェックとリントの統合
- **アクセシビリティ重視**: jsx-a11yプラグインによる包括的チェック
- **開発効率向上**: 保存時自動フォーマット・リント
- **CI/CD統合**: GitHub Actionsでの品質ゲート

### 成果物
- 0エラー、0警告でのESLintチェック通過
- 全ファイルのPrettierフォーマット適用
- TypeScript型チェック完全通過
- CI/CD自動化による品質保証体制

### 未消化作業
なし - 全てのタスクを完了し、高品質なコード管理体制を構築