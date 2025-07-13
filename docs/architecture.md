# システム設計書

## 概要
VisualTrecplansは、フィットネス・トレーニング管理を目的とした言語依存度を最小化したWebアプリケーションです。ビジュアル重視のUI、筋肉部位の色分け、人体図を活用したインタラクティブな操作により、言語の壁を越えたユーザビリティを提供します。

## システム全体構成

### アーキテクチャ概要
```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Layer                            │
├─────────────────────────────────────────────────────────────────┤
│ React + TypeScript Frontend                                     │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │ Components  │   Hooks     │   Stores    │   Services          │ │
│ │             │             │  (Zustand)  │                     │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                    │ HTTPS/REST API
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API Layer                               │
├─────────────────────────────────────────────────────────────────┤
│ Go + Gin Backend                                                │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │ Handlers    │ Middleware  │ Services    │ Repositories        │ │
│ │             │             │             │                     │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                    │ GORM
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Database Layer                             │
├─────────────────────────────────────────────────────────────────┤
│ PostgreSQL                                                      │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │   Users     │  Workouts   │   Master    │   Audit Logs       │ │
│ │             │             │    Data     │                     │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## バックエンドアーキテクチャ

### クリーンアーキテクチャ採用
```
┌───────────────────────────────────────────────────────────────────┐
│                      Presentation Layer                          │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │   Auth      │   Workout   │   Master    │     Middleware      │ │
│ │ Handlers    │  Handlers   │  Handlers   │                     │ │
│ │             │             │             │  - Authentication   │ │
│ │ - Login     │ - Create    │ - Muscle    │  - CORS            │ │
│ │ - Register  │ - List      │   Groups    │  - Rate Limiting   │ │
│ │ - Profile   │ - Update    │ - Exercises │  - Security        │ │
│ │ - Logout    │ - Delete    │ - Icons     │  - Audit Logging   │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌───────────────────────────────────────────────────────────────────┐
│                       Use Case Layer                             │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │    Auth     │   Workout   │   Master    │       Common        │ │
│ │  Service    │  Service    │  Service    │                     │ │
│ │             │             │             │ - JWT Service       │ │
│ │ - Register  │ - Create    │ - Get       │ - Email Service     │ │
│ │ - Login     │ - Update    │   Groups    │ - Cache Service     │ │
│ │ - Refresh   │ - Delete    │ - Get       │                     │ │
│ │ - Profile   │ - Stats     │   Exercises │                     │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌───────────────────────────────────────────────────────────────────┐
│                       Entity Layer                               │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │    User     │   Workout   │   Master    │      Audit          │ │
│ │   Model     │   Model     │   Models    │      Model          │ │
│ │             │             │             │                     │ │
│ │ - ID        │ - ID        │ - Muscle    │ - ID               │ │
│ │ - Username  │ - UserID    │   Group     │ - UserID           │ │
│ │ - Email     │ - Exercise  │ - Exercise  │ - Action           │ │
│ │ - Password  │ - Weight    │ - Icon      │ - Timestamp        │ │
│ │ - Timestamps│ - Reps/Sets │ │ │ │ │ │ │ │ - Details          │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌───────────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                          │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │    Auth     │   Workout   │   Master    │      Common         │ │
│ │ Repository  │ Repository  │ Repository  │                     │ │
│ │             │             │             │ - Database          │ │
│ │ - Create    │ - Create    │ - Get All   │   Connection        │ │
│ │ - FindByID  │ - FindByID  │ - Get By    │ - Migration         │ │
│ │ - FindBy    │ - FindBy    │   Category  │ - Logger            │ │
│ │   Email     │   UserID    │ - Create    │ - Config            │ │
│ │ - Update    │ - Update    │   Custom    │                     │ │
│ │ - Delete    │ - Delete    │             │                     │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
```

### ディレクトリ構造
```
backend/
├── cmd/
│   ├── main.go              # アプリケーションエントリーポイント
│   ├── migrate/             # データベースマイグレーション
│   └── seed/                # データシード
├── internal/
│   ├── handlers/            # HTTP ハンドラー
│   │   ├── auth/            # 認証関連エンドポイント
│   │   ├── workout/         # ワークアウト関連エンドポイント
│   │   └── master/          # マスタデータエンドポイント
│   ├── services/            # ビジネスロジック
│   │   ├── auth_service.go
│   │   ├── jwt_service.go
│   │   ├── workout_service.go
│   │   └── master_service.go
│   ├── repositories/        # データアクセス層
│   │   ├── workout_repository.go
│   │   └── master_repository.go
│   ├── models/              # ドメインモデル
│   │   ├── user.go
│   │   └── workout.go
│   ├── middleware/          # HTTP ミドルウェア
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── rate_limit.go
│   │   └── security.go
│   ├── validators/          # バリデーション
│   ├── database/            # データベース設定
│   └── audit/               # 監査ログ
├── pkg/                     # 再利用可能パッケージ
│   ├── config/              # 設定管理
│   ├── logger/              # ログ管理
│   └── jwt/                 # JWT ユーティリティ
├── migrations/              # データベースマイグレーション
├── seeds/                   # 初期データ
├── configs/                 # 設定ファイル
└── docs/                    # API ドキュメント
```

## フロントエンドアーキテクチャ

### React + TypeScript構成
```
┌───────────────────────────────────────────────────────────────────┐
│                      Presentation Layer                          │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │    Pages    │ Components  │   Layout    │      Common         │ │
│ │             │             │             │                     │ │
│ │ - Login     │ - Forms     │ - Header    │ - ErrorBoundary     │ │
│ │ - Register  │ - Workout   │ - Footer    │ - ProtectedRoute    │ │
│ │ - Dashboard │   History   │ - Sidebar   │ - LoadingSpinner    │ │
│ │ - Workout   │ - Body      │ - Navigation│ - Toast             │ │
│ │   Record    │   Diagram   │             │                     │ │
│ │ - History   │ - Stats     │             │                     │ │
│ │ - Profile   │   Charts    │             │                     │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌───────────────────────────────────────────────────────────────────┐
│                        Hook Layer                                │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │    Auth     │   Workout   │   Master    │      Common         │ │
│ │   Hooks     │   Hooks     │   Hooks     │     Hooks           │ │
│ │             │             │             │                     │ │
│ │ - useAuth   │ - useWorkout│ - useMuscle │ - useLocalStorage   │ │
│ │ - useLogin  │ - useWorkout│   Groups    │ - useDebounce       │ │
│ │ - useRegister│   Form     │ - useExer   │ - useClickOutside   │ │
│ │             │ - useWorkout│   cises     │ - useMediaQuery     │ │
│ │             │   History   │             │                     │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌───────────────────────────────────────────────────────────────────┐
│                       State Layer                                │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │    Auth     │   Workout   │     UI      │      Router         │ │
│ │   Store     │   Store     │   Store     │                     │ │
│ │ (Zustand)   │ (Zustand)   │ (Zustand)   │ (React Router)      │ │
│ │             │             │             │                     │ │
│ │ - user      │ - workouts  │ - theme     │ - Public Routes     │ │
│ │ - tokens    │ - stats     │ - sidebar   │ - Protected Routes  │ │
│ │ - loading   │ - filters   │ - modal     │ - Nested Routes     │ │
│ │ - error     │ - pagination│ - toast     │                     │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌───────────────────────────────────────────────────────────────────┐
│                      Service Layer                               │
│ ┌─────────────┬─────────────┬─────────────┬─────────────────────┐ │
│ │    Auth     │   Workout   │     API     │     Utils           │ │
│ │  Service    │  Service    │   Client    │                     │ │
│ │             │             │             │                     │ │
│ │ - login     │ - create    │ - axios     │ - formatters        │ │
│ │ - register  │ - update    │   instance  │ - validators        │ │
│ │ - refresh   │ - delete    │ - intercep  │ - constants         │ │
│ │ - logout    │ - list      │   tors      │ - helpers           │ │
│ │ - profile   │ - stats     │ - error     │                     │ │
│ │             │             │   handling  │                     │ │
│ └─────────────┴─────────────┴─────────────┴─────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
```

### ディレクトリ構造
```
frontend/
├── public/
│   ├── index.html
│   └── assets/
├── src/
│   ├── components/         # 再利用可能コンポーネント
│   │   ├── common/         # 汎用コンポーネント
│   │   ├── forms/          # フォーム関連
│   │   ├── layout/         # レイアウト関連
│   │   └── workout/        # ワークアウト固有
│   ├── pages/              # ページコンポーネント
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   ├── Dashboard.tsx
│   │   ├── CreateWorkout.tsx
│   │   └── WorkoutHistory.tsx
│   ├── hooks/              # カスタムフック
│   │   ├── useAuth.ts
│   │   ├── useWorkout.ts
│   │   └── useWorkoutForm.ts
│   ├── stores/             # Zustand ストア
│   │   ├── authStore.ts
│   │   └── workoutStore.ts
│   ├── services/           # API サービス
│   │   ├── api.ts
│   │   ├── auth.service.ts
│   │   └── workout.service.ts
│   ├── types/              # TypeScript 型定義
│   │   ├── auth.ts
│   │   └── workout.ts
│   ├── schemas/            # Zod バリデーション
│   │   ├── auth.ts
│   │   └── workout.ts
│   ├── utils/              # ユーティリティ
│   ├── styles/             # CSS/Tailwind
│   └── assets/             # 静的ファイル
├── package.json
├── tsconfig.json
├── tailwind.config.js
└── vite.config.ts
```

## データベース設計

### ER図
```
┌─────────────────────┐     ┌─────────────────────┐
│       Users         │     │      Workouts       │
├─────────────────────┤     ├─────────────────────┤
│ id (UUID) PK        │────<│ user_id (UUID) FK   │
│ username (VARCHAR)  │     │ id (UUID) PK        │
│ email (VARCHAR)     │     │ muscle_group (VARCHAR)│
│ password_hash       │     │ exercise_name (VARCHAR)│
│ created_at          │     │ exercise_icon (VARCHAR)│
│ updated_at          │     │ weight_kg (DECIMAL) │
│ deleted_at          │     │ reps (INTEGER)      │
└─────────────────────┘     │ sets (INTEGER)      │
                            │ notes (TEXT)        │
                            │ performed_at        │
                            │ created_at          │
                            │ updated_at          │
                            │ deleted_at          │
                            └─────────────────────┘
                                      │
                                      │ (参照)
                                      ▼
┌─────────────────────┐     ┌─────────────────────┐
│   Muscle_Groups     │     │     Exercises       │
├─────────────────────┤     ├─────────────────────┤
│ id (SERIAL) PK      │────<│ muscle_group_code FK│
│ code (VARCHAR)      │     │ id (SERIAL) PK      │
│ name_ja (VARCHAR)   │     │ name_ja (VARCHAR)   │
│ name_en (VARCHAR)   │     │ name_en (VARCHAR)   │
│ category (VARCHAR)  │     │ icon_name (VARCHAR) │
│ color_code (VARCHAR)│     │ is_custom (BOOLEAN) │
│ sort_order (INTEGER)│     │ created_by (UUID)   │
└─────────────────────┘     │ sort_order (INTEGER)│
                            └─────────────────────┘

┌─────────────────────┐     ┌─────────────────────┐
│  Exercise_Icons     │     │    Audit_Logs      │
├─────────────────────┤     ├─────────────────────┤
│ id (SERIAL) PK      │     │ id (UUID) PK        │
│ name (VARCHAR)      │     │ user_id (UUID)      │
│ svg_path (TEXT)     │     │ action (VARCHAR)    │
│ category (VARCHAR)  │     │ resource_type       │
└─────────────────────┘     │ resource_id (VARCHAR)│
                            │ details (JSONB)     │
                            │ ip_address (VARCHAR)│
                            │ user_agent (TEXT)   │
                            │ created_at          │
                            └─────────────────────┘
```

### テーブル定義

#### Users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- インデックス
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;

-- トリガー
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

#### Workouts
```sql
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    muscle_group VARCHAR(50) NOT NULL,
    exercise_name VARCHAR(100) NOT NULL,
    exercise_icon VARCHAR(50),
    weight_kg DECIMAL(5,2) CHECK (weight_kg >= 0),
    reps INTEGER CHECK (reps > 0),
    sets INTEGER CHECK (sets > 0),
    notes TEXT,
    performed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- インデックス
CREATE INDEX idx_workouts_user_id ON workouts(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_workouts_performed_at ON workouts(performed_at DESC);
CREATE INDEX idx_workouts_muscle_group ON workouts(muscle_group);
CREATE INDEX idx_workouts_exercise_name ON workouts(exercise_name);
```

#### Muscle_Groups
```sql
CREATE TABLE muscle_groups (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name_ja VARCHAR(100) NOT NULL,
    name_en VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    color_code VARCHAR(7),
    sort_order INTEGER DEFAULT 0
);

-- 初期データ
INSERT INTO muscle_groups (code, name_ja, name_en, category, color_code, sort_order) VALUES
('chest', '胸', 'Chest', 'upper', '#ff6b6b', 1),
('back', '背中', 'Back', 'upper', '#4ecdc4', 2),
('shoulders', '肩', 'Shoulders', 'upper', '#dda0dd', 3),
('arms', '腕', 'Arms', 'upper', '#96ceb4', 4),
('core', '腹', 'Core', 'core', '#ffd93d', 5),
('legs', '脚', 'Legs', 'lower', '#45b7d1', 6),
('glutes', '臀部', 'Glutes', 'lower', '#ff9999', 7),
('full_body', '全身', 'Full Body', 'full_body', '#b8b8b8', 8);
```

## セキュリティ設計

### 認証・認可
```
┌─────────────────────────────────────────────────────────────────┐
│                     Authentication Flow                        │
├─────────────────────────────────────────────────────────────────┤
│ 1. User Registration                                            │
│    ├─ Password Strength Validation (8+ chars, mixed case)      │
│    ├─ Email Uniqueness Check                                   │
│    ├─ Password Hashing (bcrypt, cost=12)                       │
│    └─ User Creation + Audit Log                                │
│                                                                 │
│ 2. User Login                                                   │
│    ├─ Rate Limiting (5 attempts per 5 minutes)                 │
│    ├─ Brute Force Protection (progressive delays)              │
│    ├─ Credential Validation                                     │
│    ├─ JWT Generation (Access + Refresh tokens)                 │
│    └─ Audit Log Creation                                        │
│                                                                 │
│ 3. Token Management                                             │
│    ├─ Access Token (15 minutes expiry)                         │
│    ├─ Refresh Token (7 days expiry)                            │
│    ├─ Automatic Token Refresh                                  │
│    └─ Secure Token Storage (HTTP-only cookies)                 │
└─────────────────────────────────────────────────────────────────┘
```

### セキュリティ対策
```
┌─────────────────────────────────────────────────────────────────┐
│                     Security Measures                          │
├─────────────────────────────────────────────────────────────────┤
│ HTTP Security Headers                                           │
│ ├─ Content-Security-Policy                                      │
│ ├─ X-Frame-Options: DENY                                        │
│ ├─ X-Content-Type-Options: nosniff                              │
│ ├─ Referrer-Policy: strict-origin-when-cross-origin            │
│ └─ Strict-Transport-Security                                    │
│                                                                 │
│ Input Validation                                                │
│ ├─ Server-side Validation (Go validators)                      │
│ ├─ Client-side Validation (Zod schemas)                        │
│ ├─ SQL Injection Prevention (GORM ORM)                         │
│ └─ XSS Prevention (input sanitization)                         │
│                                                                 │
│ Rate Limiting                                                   │
│ ├─ General API: 100 requests/minute                            │
│ ├─ Auth Endpoints: 20 requests/minute                          │
│ └─ Login Endpoint: 5 requests/minute                           │
│                                                                 │
│ Data Protection                                                 │
│ ├─ Password Hashing (bcrypt)                                   │
│ ├─ JWT Secret Management                                        │
│ ├─ Environment Variable Protection                              │
│ └─ Database Connection Encryption                               │
└─────────────────────────────────────────────────────────────────┘
```

## パフォーマンス設計

### キャッシュ戦略
```
┌─────────────────────────────────────────────────────────────────┐
│                      Caching Strategy                          │
├─────────────────────────────────────────────────────────────────┤
│ Client-Side Caching                                             │
│ ├─ React Query / SWR (future implementation)                   │
│ ├─ Zustand Persistence (authentication state)                  │
│ ├─ LocalStorage (user preferences)                             │
│ └─ Service Worker (offline capability - future)                │
│                                                                 │
│ HTTP Caching                                                    │
│ ├─ Master Data: 1 hour cache                                   │
│ ├─ Exercise Icons: 24 hour cache                               │
│ ├─ Static Assets: Long-term caching                            │
│ └─ ETag Support for conditional requests                       │
│                                                                 │
│ Database Optimization                                           │
│ ├─ Strategic Indexing                                           │
│ ├─ Query Optimization (GORM)                                   │
│ ├─ Connection Pooling                                           │
│ └─ Pagination for large datasets                               │
└─────────────────────────────────────────────────────────────────┘
```

### スケーラビリティ設計
```
┌─────────────────────────────────────────────────────────────────┐
│                   Scalability Design                           │
├─────────────────────────────────────────────────────────────────┤
│ Horizontal Scaling (Future)                                     │
│ ├─ Stateless API Design                                         │
│ ├─ Load Balancer Ready                                          │
│ ├─ Database Read Replicas                                       │
│ └─ Microservices Architecture (Phase 4)                        │
│                                                                 │
│ Resource Optimization                                           │
│ ├─ Efficient SQL Queries                                       │
│ ├─ Lazy Loading (React components)                             │
│ ├─ Image Optimization                                           │
│ └─ Bundle Splitting (Webpack/Vite)                             │
│                                                                 │
│ Monitoring & Observability                                     │
│ ├─ Structured Logging (JSON format)                            │
│ ├─ Performance Metrics                                          │
│ ├─ Error Tracking                                               │
│ └─ Health Check Endpoints                                       │
└─────────────────────────────────────────────────────────────────┘
```

## デプロイメント設計

### Docker構成
```yaml
# docker-compose.yml
version: '3.8'
services:
  frontend:
    build: 
      context: ./frontend
      dockerfile: Dockerfile.dev
    ports:
      - "3000:3000"
    environment:
      - VITE_API_URL=http://localhost:8080/api/v1
    volumes:
      - ./frontend:/app
      - /app/node_modules

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_USER=trecplans
      - DB_PASSWORD=password
      - DB_NAME=trecplans
      - JWT_SECRET=your-secret-key
    depends_on:
      - postgres
    volumes:
      - ./backend:/app

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=trecplans
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=trecplans
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  adminer:
    image: adminer:latest
    ports:
      - "8081:8080"
    depends_on:
      - postgres

volumes:
  postgres_data:
```

### CI/CD パイプライン (GitHub Actions)
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - run: cd backend && go test ./...
      - run: cd backend && go vet ./...

  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: cd frontend && npm ci
      - run: cd frontend && npm run test
      - run: cd frontend && npm run lint
      - run: cd frontend && npm run type-check

  build-and-deploy:
    needs: [test-backend, test-frontend]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build and Deploy
        run: |
          # Docker build and push
          # Deployment steps
```

## 監視・ログ設計

### ログ設計
```go
// 構造化ログ形式
type LogEntry struct {
    Timestamp time.Time   `json:"timestamp"`
    Level     string      `json:"level"`
    Message   string      `json:"message"`
    UserID    string      `json:"user_id,omitempty"`
    RequestID string      `json:"request_id,omitempty"`
    Method    string      `json:"method,omitempty"`
    Path      string      `json:"path,omitempty"`
    IP        string      `json:"ip,omitempty"`
    UserAgent string      `json:"user_agent,omitempty"`
    Duration  int64       `json:"duration_ms,omitempty"`
    Error     string      `json:"error,omitempty"`
    Details   interface{} `json:"details,omitempty"`
}
```

### 監査ログ
```go
type AuditLog struct {
    ID           uuid.UUID   `json:"id"`
    UserID       uuid.UUID   `json:"user_id"`
    Action       string      `json:"action"`
    ResourceType string      `json:"resource_type"`
    ResourceID   string      `json:"resource_id"`
    Details      interface{} `json:"details"`
    IPAddress    string      `json:"ip_address"`
    UserAgent    string      `json:"user_agent"`
    CreatedAt    time.Time   `json:"created_at"`
}
```

## 今後の拡張設計

### Phase 2: ビジュアル強化
- リアルタイム人体図更新
- アニメーション効果
- 3D人体モデル (Three.js)
- プログレッシブWebアプリ (PWA)

### Phase 3: 機能拡張
- サプリメント管理
- データエクスポート機能
- レポート生成
- モバイルアプリ (React Native)

### Phase 4: スケーラビリティ
- マイクロサービス分割
- イベント駆動アーキテクチャ
- リアルタイム機能 (WebSocket)
- 外部API連携

この設計書は、VisualTrecplansの現在の実装状況と将来の拡張性を考慮した包括的なシステム設計を示しています。