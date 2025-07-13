# VisualTrecplans ドキュメント

## 概要
このディレクトリには、VisualTrecplansプロジェクトの包括的な設計書と開発ドキュメントが含まれています。

## ドキュメント一覧

### 📋 [API仕様書](./api-specification.md)
- RESTful APIエンドポイントの詳細仕様
- リクエスト/レスポンス形式
- 認証・認可方式
- エラーハンドリング
- レート制限とキャッシュ戦略

### 🔄 [シーケンス図](./sequence-diagrams.md)
- ユーザー登録・ログインフロー
- ワークアウト記録・管理フロー
- トークンリフレッシュ機能
- エラーハンドリングフロー
- 主要な業務プロセスの詳細な流れ

### 🎨 [画面ワイヤフレーム](./wireframes.md)
- 全主要画面のワイヤフレーム
- モバイル・デスクトップ対応デザイン
- UI/UXコンポーネント仕様
- レスポンシブデザイン設計
- カラーパレットとデザインシステム

### 🏗️ [システムアーキテクチャ](./architecture.md)
- システム全体構成図
- クリーンアーキテクチャ設計
- データベース設計とER図
- セキュリティ設計
- パフォーマンス・スケーラビリティ設計

### 💻 [開発ガイド](./development-guide.md)
- 開発環境セットアップ手順
- コーディング規約・ベストプラクティス
- テスト戦略と実装例
- デバッグ手順
- パフォーマンス最適化
- デプロイメント手順

## 技術スタック

### バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Gin
- **データベース**: PostgreSQL 15+
- **ORM**: GORM
- **認証**: JWT
- **ログ**: 構造化ログ (JSON)
- **コンテナ**: Docker & Docker Compose

### フロントエンド
- **言語**: TypeScript 5+
- **フレームワーク**: React 18+
- **ビルドツール**: Vite
- **スタイリング**: Tailwind CSS
- **状態管理**: Zustand
- **フォーム**: React Hook Form + Zod
- **ルーティング**: React Router v6

### 開発ツール
- **バージョン管理**: Git
- **CI/CD**: GitHub Actions
- **テスト**: Jest, React Testing Library, Playwright
- **リンター**: ESLint, Prettier, Go vet
- **エディタ**: VSCode (推奨設定含む)

## プロジェクト特徴

### 🌐 言語依存度最小化
- 筋肉部位の色分けによる視覚的識別
- アイコンベースのインターフェース
- インタラクティブ人体図
- 多言語対応基盤

### 🎯 ビジュアル重視設計
- 直感的な操作性
- レスポンシブデザイン
- アクセシビリティ対応
- ダークモード対応

### 🔒 セキュリティ重視
- JWT認証システム
- レート制限・ブルートフォース保護
- 包括的な監査ログ
- セキュリティヘッダー設定

### ⚡ パフォーマンス最適化
- 効率的なデータベース設計
- クライアントサイドキャッシュ
- 遅延読み込み
- バンドル最適化

## 開発フェーズ

### ✅ Phase 1: MVP (完了)
- [x] プロジェクト初期セットアップ
- [x] 認証機能実装
- [x] トレーニング記録機能MVP版

### 🚧 Phase 2: ビジュアル強化 (予定)
- [ ] リアルタイム人体図UI実装
- [ ] 進捗可視化機能強化
- [ ] アニメーション・インタラクション改善

### 📋 Phase 3: 機能拡張 (予定)
- [ ] サプリメント管理機能
- [ ] データエクスポート機能
- [ ] レポート・分析機能

### 🚀 Phase 4: スケーラビリティ (予定)
- [ ] ソーシャル機能
- [ ] 外部API連携
- [ ] モバイルアプリ対応

## クイックスタート

### Docker使用 (推奨)
```bash
# リポジトリクローン
git clone <repository-url>
cd trecplans

# 環境変数設定
cp .env.docker.example .env.docker

# 開発環境起動
docker-compose up -d

# データベース初期化
docker-compose exec backend go run cmd/migrate/main.go
docker-compose exec backend go run cmd/seed/main.go

# アクセス確認
# フロントエンド: http://localhost:3000
# バックエンドAPI: http://localhost:8080
# データベース管理: http://localhost:8081
```

### ローカル開発
詳細は [開発ガイド](./development-guide.md) を参照してください。

## API エンドポイント概要

### 認証
- `POST /api/v1/auth/register` - ユーザー登録
- `POST /api/v1/auth/login` - ログイン
- `POST /api/v1/auth/refresh` - トークンリフレッシュ
- `GET /api/v1/auth/profile` - プロフィール取得

### ワークアウト
- `POST /api/v1/workouts` - ワークアウト記録
- `GET /api/v1/workouts` - ワークアウト一覧
- `GET /api/v1/workouts/:id` - ワークアウト詳細
- `PUT /api/v1/workouts/:id` - ワークアウト更新
- `DELETE /api/v1/workouts/:id` - ワークアウト削除
- `GET /api/v1/workouts/stats` - 統計情報

### マスタデータ
- `GET /api/v1/muscle-groups` - 筋肉部位一覧
- `GET /api/v1/exercises` - エクササイズ一覧
- `POST /api/v1/exercises/custom` - カスタムエクササイズ作成
- `GET /api/v1/exercise-icons` - アイコン一覧

## データベース設計

### 主要テーブル
- **users** - ユーザー情報
- **workouts** - ワークアウト記録
- **muscle_groups** - 筋肉部位マスター
- **exercises** - エクササイズマスター
- **exercise_icons** - アイコンマスター
- **audit_logs** - 監査ログ

詳細は [システムアーキテクチャ](./architecture.md) を参照してください。

## 貢献ガイドライン

### コミット規約
```
feat: 新機能追加
fix: バグ修正
docs: ドキュメント更新
style: コードスタイル修正
refactor: リファクタリング
test: テスト追加・修正
chore: その他（依存関係更新など）
```

### プルリクエスト手順
1. 機能ブランチ作成 (`feature/新機能名`)
2. 実装・テスト作成
3. リンター・テスト通過確認
4. プルリクエスト作成
5. コードレビュー
6. マージ

## ライセンス
このプロジェクトは MIT License の下で公開されています。

## サポート
質問や問題がある場合は、以下を確認してください：
1. [開発ガイド](./development-guide.md) のトラブルシューティング
2. GitHubのIssue
3. プロジェクトドキュメント

---

このドキュメントは、開発チームがプロジェクトを効率的に理解し、貢献できるよう作成されています。定期的に更新され、最新の実装状況を反映します。