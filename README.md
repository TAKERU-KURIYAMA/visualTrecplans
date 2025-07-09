# VisualTrecplans

言語依存度を最小化した、ビジュアル重視のトレーニング・サプリ管理アプリケーション

🌐 **Live Demo**: https://trecplans.com

## 🎯 プロジェクト概要

### コンセプト
- **言語に依存しないUI**: 人体図・アイコン・数値中心の操作
- **グローバル展開**: 非英語圏での既存サービス課題を解決
- **アクセシビリティ**: 色覚多様性対応、階層化UI設計

### ターゲット
- フィットネス愛好者（特に非英語圏）
- 既存英語サービスに言語障壁を感じるユーザー
- ビジュアル操作を好むユーザー

## 🏗️ アーキテクチャ

```
Frontend (React + TypeScript)
    ↕ REST API
Backend (Go + Gin)
    ↕ ORM (GORM)
PostgreSQL Database
```

## 🛠️ 技術スタック

### フロントエンド
- **React 18** + **TypeScript** - 型安全な開発
- **Tailwind CSS** - 高速スタイリング
- **Framer Motion** - アニメーション
- **React Hook Form** - フォーム管理
- **D3.js** - データ可視化
- **Zustand** - 状態管理

### バックエンド
- **Go 1.21** + **Gin Framework** - 高性能API
- **GORM** - PostgreSQL ORM
- **JWT認証** - ステートレス認証
- **Docker** - コンテナ化

### データベース
- **PostgreSQL 14+** - メインDB
- **JSON型** - カスタムデータ（低頻度更新のみ）

### インフラ
- **ConoHa VPS** - 本番環境
- **Nginx** - リバースプロキシ
- **systemd** - プロセス管理
- **Let's Encrypt** - SSL証明書

## 🚀 セットアップ

### 前提条件
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Docker & Docker Compose

### 開発環境構築

1. **リポジトリクローン**
```bash
git clone https://github.com/username/visualTrecplans.git
cd visualTrecplans
```

2. **Docker Compose起動**
```bash
docker-compose up -d
```

3. **フロントエンド起動**
```bash
cd frontend
npm install
npm run dev
```

4. **バックエンド起動**
```bash
cd backend
go mod tidy
go run main.go
```

### 環境変数
```bash
# backend/.env
DB_HOST=localhost
DB_PORT=5432
DB_USER=trecplans
DB_PASSWORD=password
DB_NAME=trecplans
JWT_SECRET=your-secret-key
PORT=8080
DOMAIN=trecplans.com
```

## 📁 プロジェクト構造

```
visualTrecplans/
├── frontend/                 # React アプリケーション
│   ├── src/
│   │   ├── components/      # UIコンポーネント
│   │   ├── pages/          # ページコンポーネント
│   │   ├── hooks/          # カスタムフック
│   │   ├── stores/         # Zustand ストア
│   │   └── types/          # TypeScript 型定義
├── backend/                 # Go API サーバー
│   ├── cmd/                # エントリーポイント
│   ├── internal/
│   │   ├── handlers/       # HTTPハンドラー
│   │   ├── models/         # データモデル
│   │   ├── services/       # ビジネスロジック
│   │   └── database/       # DB接続・マイグレーション
├── docker-compose.yml      # 開発環境
├── Dockerfile              # 本番環境用
└── docs/                   # ドキュメント
```

## 🎨 主要機能

### Phase 1: MVP
- [x] ユーザー認証（JWT）
- [x] ドロップダウン式トレーニング記録
- [x] 基本的な記録一覧表示
- [x] 人体図モック表示

### Phase 2: ビジュアル強化
- [ ] 人体図UI実装（SVG/Canvas選定中）
- [ ] アニメーション効果
- [ ] 進捗グラフ表示

### Phase 3: 機能拡張
- [ ] サプリ管理
- [ ] カレンダー表示
- [ ] データエクスポート

### Phase 4: スケーラビリティ
- [ ] ソーシャル機能
- [ ] 外部API連携
- [ ] モバイルアプリ対応

## 🔧 開発ガイド

### API エンドポイント

```go
// 認証
POST   /api/v1/auth/login
POST   /api/v1/auth/register

// トレーニング
GET    /api/v1/workouts           // 一覧取得
POST   /api/v1/workouts           // 記録作成
PUT    /api/v1/workouts/:id       // 記録更新
DELETE /api/v1/workouts/:id       // 記録削除

// マスタデータ
GET    /api/v1/body-parts         // 筋肉部位一覧
GET    /api/v1/exercises/icons    // エクササイズアイコン
```

### データベーススキーマ

```sql
-- ユーザー
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ワークアウト記録
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    muscle_group VARCHAR(50) NOT NULL,
    exercise_name VARCHAR(100) NOT NULL,
    exercise_icon VARCHAR(50),
    weight_kg DECIMAL(5,2),
    reps INTEGER,
    sets INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### コンポーネント設計

```typescript
// 人体図コンポーネント（Phase 2）
interface BodyMapProps {
  onPartSelect: (part: MusclePart) => void;
  selectedPart?: MusclePart;
  highlightedParts: MusclePart[];
}

// トレーニング記録コンポーネント
interface WorkoutFormProps {
  muscleGroup: string;
  onSubmit: (workout: WorkoutData) => void;
}
```

## 🧪 テスト

### フロントエンド
```bash
cd frontend
npm run test              # ユニットテスト
npm run test:e2e         # E2Eテスト
npm run test:coverage    # カバレッジレポート
```

### バックエンド
```bash
cd backend
go test ./...                    # 全テスト実行
go test -cover ./...            # カバレッジ付き
go test -race ./...             # レースコンディション検査
```

## 🚢 デプロイ

### 本番環境デプロイ
```bash
# Docker イメージビルド
docker build -t trecplans .

# ConoHa VPS デプロイ
ssh user@trecplans.com
docker pull trecplans
systemctl reload trecplans
```

### CI/CDパイプライン
GitHub Actions による自動デプロイ設定済み
- `main` ブランチへのpush時に自動デプロイ
- テスト → ビルド → デプロイの流れ

## 📊 パフォーマンス目標

- **初期表示**: 2秒以内
- **API応答時間**: 200ms以内
- **人体図操作**: 60fps維持
- **モバイル対応**: PWA対応

## 🔒 セキュリティ

- JWT認証
- bcryptパスワードハッシュ化
- CSPヘッダー設定
- レート制限（100req/min/IP）
- SQL injection対策（GORM使用）

## 🌍 多言語対応

### Phase 1対応言語
- 日本語
- 英語
- 簡体字中国語

### 言語依存最小化戦略
- 主要操作: ビジュアル・数値のみ
- カスタム名: ユーザーの母国語で自由入力
- UI要素: アイコン・色・パターンで識別

## 🤝 コントリビューション

1. Issue作成またはPull Request
2. 開発ブランチでの作業
3. テスト・リント通過確認
4. レビュー → マージ

### コーディング規約
- **Go**: gofmt + golint
- **TypeScript**: ESLint + Prettier
- **コミット**: Conventional Commits

## 📝 ライセンス

MIT License

## 📞 サポート

- **Issues**: GitHub Issues
- **Documentation**: `/docs` フォルダ
- **Wiki**: プロジェクトWiki参照

---

## 🏃‍♂️ クイックスタート

```bash
# 1. クローン
git clone https://github.com/username/visualTrecplans.git

# 2. 環境構築
docker-compose up -d

# 3. アクセス
# Frontend: http://localhost:3000
# Backend:  http://localhost:8080
# DB:       postgresql://localhost:5432
# Live:     https://trecplans.com
```

開発を始める準備が整いました！🎉