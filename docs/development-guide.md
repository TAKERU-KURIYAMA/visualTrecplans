# 開発ガイド

## 概要
VisualTrecplansプロジェクトの開発環境セットアップから実装までの包括的なガイドです。

## 開発環境要件

### システム要件
- **OS**: Windows 10/11, macOS 10.15+, Ubuntu 18.04+
- **RAM**: 8GB以上推奨
- **ストレージ**: 10GB以上の空き容量

### 必要なソフトウェア
```bash
# Backend
- Go 1.21+
- PostgreSQL 15+

# Frontend  
- Node.js 18+
- npm 9+ または yarn 1.22+

# Tools
- Docker 24+
- Docker Compose 2+
- Git 2.30+
- VSCode (推奨エディタ)
```

## セットアップ手順

### 1. リポジトリクローン
```bash
git clone https://github.com/your-username/visualtrecplans.git
cd visualtrecplans
```

### 2. Docker環境での開発 (推奨)
```bash
# 環境変数設定
cp .env.docker.example .env.docker

# Docker環境起動
docker-compose up -d

# データベース初期化
docker-compose exec backend go run cmd/migrate/main.go
docker-compose exec backend go run cmd/seed/main.go

# 動作確認
curl http://localhost:8080/health
curl http://localhost:3000
```

### 3. ローカル環境での開発

#### バックエンド
```bash
cd backend

# 依存関係インストール
go mod download

# 環境変数設定
cp .env.example .env
# .envファイルを編集してデータベース接続情報を設定

# データベース起動 (Docker使用)
docker run -d \
  --name postgres-dev \
  -e POSTGRES_USER=trecplans \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=trecplans \
  -p 5432:5432 \
  postgres:15-alpine

# マイグレーション実行
go run cmd/migrate/main.go

# シードデータ投入
go run cmd/seed/main.go

# 開発サーバー起動
go run cmd/main.go
```

#### フロントエンド
```bash
cd frontend

# 依存関係インストール
npm install

# 環境変数設定
cp .env.example .env.local
# .env.localを編集してAPI URLを設定

# 開発サーバー起動
npm run dev
```

## 開発ワークフロー

### ブランチ戦略
```
main ──────────────────────────────── (本番環境)
 │
 └── develop ──────────────────────── (統合環境)
      │
      ├── feature/auth-system ────── (機能開発)
      ├── feature/workout-form ───── (機能開発)
      └── hotfix/critical-bug ────── (緊急修正)
```

### Git フロー
```bash
# 新機能開発
git checkout develop
git pull origin develop
git checkout -b feature/new-feature

# 開発作業
git add .
git commit -m "feat: add new feature description"

# プッシュとプルリクエスト
git push origin feature/new-feature
# GitHub/GitLabでプルリクエスト作成
```

### コミットメッセージ規約
```bash
# 形式: <type>(<scope>): <description>

# 例:
feat(auth): add JWT token refresh functionality
fix(workout): resolve exercise dropdown filtering issue  
docs(api): update authentication endpoints documentation
style(frontend): apply consistent button styling
refactor(backend): extract workout validation logic
test(auth): add unit tests for login flow
chore(deps): update dependencies to latest versions
```

## コーディング規約

### Go (バックエンド)
```go
// パッケージ命名: 小文字、単数形
package user

// 関数命名: PascalCase (公開), camelCase (非公開)
func CreateUser(data *CreateUserRequest) (*User, error) {}
func validatePassword(password string) error {}

// 構造体命名: PascalCase
type UserService struct {
    repo UserRepository
    logger *Logger
}

// インターフェース命名: -er suffix
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
}

// エラーハンドリング
if err != nil {
    logger.Error("Failed to create user", "error", err)
    return nil, fmt.Errorf("failed to create user: %w", err)
}

// コンテキスト使用
func (s *userService) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 実装
}
```

### TypeScript (フロントエンド)
```typescript
// インターフェース命名: PascalCase
interface User {
  id: string;
  username: string;
  email: string;
  createdAt: string;
}

// 型エイリアス
type UserFormData = Pick<User, 'username' | 'email'>;

// 関数命名: camelCase
const createUser = async (data: UserFormData): Promise<User> => {
  // 実装
};

// Reactコンポーネント: PascalCase
const UserProfile: React.FC<UserProfileProps> = ({ user }) => {
  // 実装
};

// カスタムフック: use- prefix
const useAuth = () => {
  // 実装
};

// 定数: UPPER_SNAKE_CASE
const API_BASE_URL = 'http://localhost:8080/api/v1';
```

### CSS/Tailwind
```css
/* クラス命名: BEM記法またはTailwind */
.workout-form__input--error {
  @apply border-red-500 text-red-900;
}

/* Tailwindクラス順序 */
.button {
  @apply 
    /* Layout */
    flex items-center justify-center
    /* Size */
    w-full h-12 px-4 py-2
    /* Typography */
    text-base font-medium text-white
    /* Background & Border */
    bg-blue-600 border border-transparent rounded-lg
    /* Effects */
    hover:bg-blue-700 focus:ring-2 focus:ring-blue-500
    /* Transitions */
    transition-colors duration-200;
}
```

## テスト戦略

### バックエンドテスト
```go
// ユニットテスト例
func TestUserService_Create(t *testing.T) {
    // テストセットアップ
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, logger)
    
    // テストケース
    tests := []struct {
        name    string
        input   *CreateUserRequest
        want    *User
        wantErr bool
    }{
        {
            name: "valid user creation",
            input: &CreateUserRequest{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "Password123!",
            },
            want: &User{
                Username: "testuser",
                Email:    "test@example.com",
            },
            wantErr: false,
        },
    }
    
    // テスト実行
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := service.Create(context.Background(), tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got.Username, tt.want.Username) {
                t.Errorf("Create() = %v, want %v", got, tt.want)
            }
        })
    }
}

// 統合テスト
func TestUserHandler_Create_Integration(t *testing.T) {
    // テストDBセットアップ
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // HTTPリクエストテスト
    router := setupRouter(db)
    w := httptest.NewRecorder()
    
    payload := `{"username":"testuser","email":"test@example.com","password":"Password123!"}`
    req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

### フロントエンドテスト
```typescript
// React Testing Library + Jest
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { LoginForm } from './LoginForm';

describe('LoginForm', () => {
  test('renders login form correctly', () => {
    render(<LoginForm onSubmit={vi.fn()} />);
    
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
  });

  test('submits form with valid data', async () => {
    const mockSubmit = vi.fn();
    render(<LoginForm onSubmit={mockSubmit} />);
    
    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: 'test@example.com' }
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: 'password123' }
    });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    
    await waitFor(() => {
      expect(mockSubmit).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123'
      });
    });
  });

  test('shows validation errors for invalid input', async () => {
    render(<LoginForm onSubmit={vi.fn()} />);
    
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    
    await waitFor(() => {
      expect(screen.getByText(/email is required/i)).toBeInTheDocument();
      expect(screen.getByText(/password is required/i)).toBeInTheDocument();
    });
  });
});

// E2Eテスト (Playwright)
import { test, expect } from '@playwright/test';

test('user can register and login', async ({ page }) => {
  // 登録
  await page.goto('/register');
  await page.fill('[data-testid=username]', 'testuser');
  await page.fill('[data-testid=email]', 'test@example.com');
  await page.fill('[data-testid=password]', 'Password123!');
  await page.click('[data-testid=submit]');
  
  await expect(page).toHaveURL('/dashboard');
  await expect(page.locator('[data-testid=welcome]')).toContainText('Welcome, testuser');
  
  // ログアウト
  await page.click('[data-testid=logout]');
  await expect(page).toHaveURL('/login');
  
  // ログイン
  await page.fill('[data-testid=email]', 'test@example.com');
  await page.fill('[data-testid=password]', 'Password123!');
  await page.click('[data-testid=login]');
  
  await expect(page).toHaveURL('/dashboard');
});
```

## デバッグ手順

### バックエンドデバッグ
```bash
# ログレベル設定
export LOG_LEVEL=debug

# Delve デバッガー使用
go install github.com/go-delve/delve/cmd/dlv@latest

# デバッグサーバー起動
dlv debug cmd/main.go -- --config=.env

# VSCode設定 (.vscode/launch.json)
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Go Debug",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "cmd/main.go",
      "envFile": "${workspaceFolder}/.env",
      "args": []
    }
  ]
}
```

### フロントエンドデバッグ
```bash
# ブラウザ開発者ツール
# - Console: エラーログ確認
# - Network: API通信確認  
# - Application: ローカルストレージ確認

# React Developer Tools
# - Component tree inspection
# - State management debugging

# Redux DevTools (Zustand)
npm install @redux-devtools/extension
```

## パフォーマンス最適化

### バックエンド最適化
```go
// データベースクエリ最適化
func (r *workoutRepository) GetUserWorkouts(ctx context.Context, userID uuid.UUID) ([]*Workout, error) {
    var workouts []*Workout
    
    // N+1問題を避けるためのPreload
    err := r.db.WithContext(ctx).
        Preload("MuscleGroup").
        Where("user_id = ?", userID).
        Order("performed_at DESC").
        Limit(20).  // ページネーション
        Find(&workouts).Error
        
    return workouts, err
}

// コネクションプール設定
func setupDatabase() *gorm.DB {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database")
    }
    
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db
}
```

### フロントエンド最適化
```typescript
// React.memo for component optimization
const WorkoutCard = React.memo<WorkoutCardProps>(({ workout }) => {
  return (
    <div className="workout-card">
      {/* コンポーネント内容 */}
    </div>
  );
});

// useMemo for expensive calculations
const WorkoutStats = ({ workouts }: { workouts: Workout[] }) => {
  const stats = useMemo(() => {
    return calculateWorkoutStats(workouts);
  }, [workouts]);

  return <div>{/* 統計表示 */}</div>;
};

// Lazy loading for route components
const Dashboard = lazy(() => import('./pages/Dashboard'));
const WorkoutHistory = lazy(() => import('./pages/WorkoutHistory'));

// Virtualization for large lists
import { FixedSizeList as List } from 'react-window';

const VirtualizedWorkoutList = ({ workouts }: { workouts: Workout[] }) => {
  const Row = ({ index, style }) => (
    <div style={style}>
      <WorkoutCard workout={workouts[index]} />
    </div>
  );

  return (
    <List
      height={600}
      itemCount={workouts.length}
      itemSize={120}
    >
      {Row}
    </List>
  );
};
```

## トラブルシューティング

### よくある問題と解決法

#### 1. データベース接続エラー
```bash
# 問題: connection refused
# 解決: PostgreSQLサービス状態確認
docker ps | grep postgres
docker-compose logs postgres

# 環境変数確認
echo $DB_HOST $DB_USER $DB_PASSWORD $DB_NAME
```

#### 2. フロントエンドビルドエラー
```bash
# 問題: Module not found
# 解決: 依存関係再インストール
rm -rf node_modules package-lock.json
npm install

# TypeScriptエラー確認
npm run type-check
```

#### 3. API通信エラー
```bash
# 問題: CORS error
# 解決: バックエンドCORS設定確認
# internal/middleware/cors.go

# 問題: 401 Unauthorized
# 解決: JWT トークン確認
# ブラウザ開発者ツール > Application > Local Storage
```

#### 4. Docker関連問題
```bash
# コンテナ状態確認
docker-compose ps

# ログ確認
docker-compose logs backend
docker-compose logs frontend

# ボリューム再作成
docker-compose down -v
docker-compose up -d
```

## 本番環境デプロイ

### 環境変数設定
```bash
# 本番環境変数 (.env.production)
NODE_ENV=production
VITE_API_URL=https://api.trecplans.com/api/v1

# バックエンド環境変数
APP_ENV=production
DB_HOST=production-db-host
JWT_SECRET=secure-random-secret
LOG_LEVEL=info
```

### ビルド手順
```bash
# フロントエンドビルド
cd frontend
npm run build

# バックエンドビルド
cd backend
go build -o main cmd/main.go

# Dockerイメージビルド
docker build -t trecplans-frontend:latest -f frontend/Dockerfile.prod ./frontend
docker build -t trecplans-backend:latest -f backend/Dockerfile.prod ./backend
```

### デプロイメント
```bash
# Docker Compose本番環境
docker-compose -f docker-compose.prod.yml up -d

# ヘルスチェック
curl https://api.trecplans.com/health
curl https://trecplans.com
```

## 監視・ログ

### ログ設定
```go
// 構造化ログ設定
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

// 使用例
logger.Info("User logged in",
    slog.String("user_id", userID),
    slog.String("ip", clientIP),
    slog.Duration("duration", time.Since(start)),
)
```

### メトリクス収集
```go
// Prometheus メトリクス例
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"method", "endpoint"},
    )
)
```

このガイドにより、開発者がVisualTrecplansプロジェクトに効率的に参加し、高品質なコードを継続的に提供できるようになります。