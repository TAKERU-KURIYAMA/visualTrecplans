# handlers ディレクトリ

## 概要
HTTPハンドラー（コントローラー層）を配置するディレクトリです。
HTTPリクエストの受信、レスポンスの返却、およびサービス層への委譲を担当します。

## 責務
- HTTPリクエストのパース
- 入力値の基本的な検証
- サービス層の呼び出し
- HTTPレスポンスの生成
- エラーハンドリング

## ファイル構成例
```
handlers/
├── auth.go         # 認証関連のハンドラー
├── user.go         # ユーザー関連のハンドラー
├── workout.go      # トレーニング関連のハンドラー
├── router.go       # ルーティング設定
└── middleware.go   # ハンドラー用ミドルウェア
```

## 実装例
```go
// handlers/auth.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/trecplans/backend/internal/services"
    "github.com/trecplans/backend/internal/validators"
)

type AuthHandler struct {
    authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req validators.LoginRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    token, user, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user":  user,
    })
}
```

## ベストプラクティス
- ビジネスロジックは含めない
- HTTPステータスコードを適切に使用
- エラーレスポンスの形式を統一
- 依存性注入を使用