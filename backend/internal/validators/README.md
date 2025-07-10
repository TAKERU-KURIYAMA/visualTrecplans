# validators ディレクトリ

## 概要
入力値検証とリクエスト/レスポンスの型定義を配置するディレクトリです。
APIのデータ検証ルールとDTO（Data Transfer Object）を管理します。

## 責務
- リクエストボディの検証
- クエリパラメータの検証
- カスタムバリデーションルール
- エラーメッセージの管理

## ファイル構成例
```
validators/
├── auth.go         # 認証関連の検証
├── user.go         # ユーザー関連の検証
├── workout.go      # ワークアウト関連の検証
├── common.go       # 共通検証ルール
└── errors.go       # エラーメッセージ定義
```

## 実装例
```go
// validators/auth.go
package validators

import (
    "github.com/go-playground/validator/v10"
)

// LoginRequest ログインリクエスト
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// RegisterRequest ユーザー登録リクエスト
type RegisterRequest struct {
    Email           string `json:"email" binding:"required,email"`
    Username        string `json:"username" binding:"required,min=3,max=50"`
    Password        string `json:"password" binding:"required,min=8,containsany=!@#$%^&*"`
    PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

// validators/workout.go
package validators

import (
    "time"
    "github.com/google/uuid"
)

// CreateWorkoutRequest ワークアウト作成リクエスト
type CreateWorkoutRequest struct {
    Date      time.Time              `json:"date" binding:"required"`
    Exercises []ExerciseInput        `json:"exercises" binding:"required,min=1,dive"`
    Notes     string                 `json:"notes" binding:"max=500"`
}

// ExerciseInput エクササイズ入力
type ExerciseInput struct {
    ExerciseID uuid.UUID      `json:"exercise_id" binding:"required"`
    Sets       []SetInput     `json:"sets" binding:"required,min=1,dive"`
}

// SetInput セット入力
type SetInput struct {
    Reps   int     `json:"reps" binding:"required,min=1"`
    Weight float64 `json:"weight" binding:"min=0"`
    Unit   string  `json:"unit" binding:"required,oneof=kg lb"`
}

// validators/common.go
package validators

import (
    "regexp"
    "github.com/go-playground/validator/v10"
)

// カスタムバリデーション関数
func ValidatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    
    // 最低8文字
    if len(password) < 8 {
        return false
    }
    
    // 大文字を含む
    if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        return false
    }
    
    // 数字を含む
    if !regexp.MustCompile(`[0-9]`).MatchString(password) {
        return false
    }
    
    // 特殊文字を含む
    if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
        return false
    }
    
    return true
}

// RegisterCustomValidators カスタムバリデータの登録
func RegisterCustomValidators(v *validator.Validate) {
    v.RegisterValidation("strong_password", ValidatePassword)
}
```

## バリデーションタグ
- `required`: 必須フィールド
- `email`: メールアドレス形式
- `min`: 最小値/最小長
- `max`: 最大値/最大長
- `oneof`: 指定値のいずれか
- `dive`: ネストした構造体の検証

## エラーメッセージのカスタマイズ
```go
// validators/errors.go
package validators

import (
    "fmt"
    "github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
    errors := make(map[string]string)
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            field := e.Field()
            tag := e.Tag()
            
            switch tag {
            case "required":
                errors[field] = fmt.Sprintf("%s is required", field)
            case "email":
                errors[field] = "Invalid email format"
            case "min":
                errors[field] = fmt.Sprintf("%s must be at least %s", field, e.Param())
            default:
                errors[field] = fmt.Sprintf("%s is invalid", field)
            }
        }
    }
    
    return errors
}
```