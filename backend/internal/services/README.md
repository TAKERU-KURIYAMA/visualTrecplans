# services ディレクトリ

## 概要
ビジネスロジック層（ユースケース層）を配置するディレクトリです。
アプリケーションの中核となるビジネスルールとロジックを実装します。

## 責務
- ビジネスロジックの実装
- トランザクション管理
- 複数のリポジトリの調整
- ドメインモデルの操作
- 外部サービスとの連携

## ファイル構成例
```
services/
├── auth.go         # 認証サービス
├── user.go         # ユーザー管理サービス
├── workout.go      # ワークアウトサービス
├── stats.go        # 統計・分析サービス
└── export.go       # データエクスポートサービス
```

## 実装例
```go
// services/auth.go
package services

import (
    "errors"
    "time"
    "github.com/trecplans/backend/internal/models"
    "github.com/trecplans/backend/internal/database"
    "github.com/trecplans/backend/pkg/auth"
    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Register(email, username, password string) (*models.User, error)
    Login(email, password string) (string, *models.User, error)
    ValidateToken(token string) (*models.User, error)
}

type authService struct {
    userRepo database.UserRepository
    jwtUtil  auth.JWTUtil
}

func NewAuthService(userRepo database.UserRepository, jwtUtil auth.JWTUtil) AuthService {
    return &authService{
        userRepo: userRepo,
        jwtUtil:  jwtUtil,
    }
}

func (s *authService) Register(email, username, password string) (*models.User, error) {
    // ユーザーの存在確認
    existing, _ := s.userRepo.FindByEmail(email)
    if existing != nil {
        return nil, errors.New("email already exists")
    }
    
    // パスワードのハッシュ化
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // ユーザー作成
    user := &models.User{
        Email:        email,
        Username:     username,
        PasswordHash: string(hashedPassword),
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
    // ユーザー検索
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", nil, errors.New("invalid credentials")
    }
    
    // パスワード検証
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", nil, errors.New("invalid credentials")
    }
    
    // JWTトークン生成
    token, err := s.jwtUtil.GenerateToken(user.ID.String())
    if err != nil {
        return "", nil, err
    }
    
    return token, user, nil
}
```

## ベストプラクティス
- インターフェースを定義して依存性を注入
- 単一責任の原則に従う
- エラーハンドリングを適切に実装
- テスタブルな設計を心がける