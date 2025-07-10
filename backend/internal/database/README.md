# database ディレクトリ

## 概要
データベース接続管理とリポジトリパターンの実装を配置するディレクトリです。
永続化層の抽象化を提供し、データアクセスロジックを管理します。

## 責務
- データベース接続の管理
- リポジトリインターフェースの定義
- クエリの実装
- トランザクション管理
- マイグレーション実行

## ファイル構成例
```
database/
├── connection.go    # DB接続管理
├── repository.go    # リポジトリインターフェース定義
├── user_repo.go     # ユーザーリポジトリ実装
├── workout_repo.go  # ワークアウトリポジトリ実装
└── transaction.go   # トランザクション管理
```

## 実装例
```go
// database/connection.go
package database

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func Connect(config Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        return nil, err
    }
    
    return db, nil
}

// database/repository.go
package database

import (
    "github.com/visualtrecplans/backend/internal/models"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id string) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id string) error
}

type WorkoutRepository interface {
    Create(workout *models.Workout) error
    FindByID(id string) (*models.Workout, error)
    FindByUserID(userID string, limit, offset int) ([]*models.Workout, error)
    Update(workout *models.Workout) error
    Delete(id string) error
}

// database/user_repo.go
package database

import (
    "gorm.io/gorm"
    "github.com/visualtrecplans/backend/internal/models"
)

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

## ベストプラクティス
- リポジトリパターンでデータアクセスを抽象化
- インターフェースで依存関係を逆転
- エラーハンドリングを統一
- N+1問題に注意（Preloadを適切に使用）