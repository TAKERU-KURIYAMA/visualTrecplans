# models ディレクトリ

## 概要
ドメインモデルとデータベースモデルを定義するディレクトリです。
アプリケーションの中核となるデータ構造を管理します。

## モデルの種類
- **ドメインモデル**: ビジネスロジックを表現
- **データベースモデル**: 永続化層のスキーマ定義
- **DTOモデル**: データ転送オブジェクト

## ファイル構成例
```
models/
├── user.go         # ユーザーモデル
├── workout.go      # ワークアウトモデル
├── exercise.go     # エクササイズモデル
├── supplement.go   # サプリメントモデル
└── common.go       # 共通モデル（タイムスタンプ等）
```

## 実装例
```go
// models/user.go
package models

import (
    "time"
    "github.com/google/uuid"
)

// User ドメインモデル
type User struct {
    ID           uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email        string        `json:"email" gorm:"unique;not null"`
    Username     string        `json:"username" gorm:"not null"`
    PasswordHash string        `json:"-" gorm:"not null"`
    Preferences  UserPreferences `json:"preferences" gorm:"embedded"`
    CreatedAt    time.Time     `json:"created_at"`
    UpdatedAt    time.Time     `json:"updated_at"`
}

// UserPreferences ユーザー設定
type UserPreferences struct {
    Language string `json:"language" gorm:"default:'ja'"`
    Theme    string `json:"theme" gorm:"default:'light'"`
    Units    string `json:"units" gorm:"default:'metric'"`
}

// TableName GORMテーブル名指定
func (User) TableName() string {
    return "users"
}

// Validate バリデーション
func (u *User) Validate() error {
    if u.Email == "" {
        return errors.New("email is required")
    }
    if u.Username == "" {
        return errors.New("username is required")
    }
    return nil
}
```

## タグの使用
- `json`: JSON変換時のフィールド名
- `gorm`: データベースマッピング
- `validate`: バリデーションルール

## ベストプラクティス
- モデルは単一責任の原則に従う
- ビジネスロジックはサービス層に配置
- プライベートフィールドは適切に隠蔽