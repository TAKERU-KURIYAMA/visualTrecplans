# middleware ディレクトリ

## 概要
HTTPミドルウェアを配置するディレクトリです。
リクエスト/レスポンスの前後処理、横断的関心事の実装を管理します。

## ミドルウェアの種類
- 認証・認可
- ログ記録
- CORS設定
- レート制限
- エラーハンドリング
- リクエストID付与

## ファイル構成例
```
middleware/
├── auth.go         # JWT認証ミドルウェア
├── cors.go         # CORS設定
├── logger.go       # アクセスログ
├── rate_limit.go   # レート制限
├── recovery.go     # パニックリカバリー
└── request_id.go   # リクエストID生成
```

## 実装例
```go
// middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/trecplans/backend/pkg/auth"
    "github.com/trecplans/backend/internal/services"
)

func AuthRequired(jwtUtil auth.JWTUtil, userService services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authorizationヘッダーの取得
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Bearerトークンの抽出
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        
        // トークン検証
        claims, err := jwtUtil.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // ユーザー情報の取得
        user, err := userService.GetByID(claims.UserID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }
        
        // コンテキストにユーザー情報を設定
        c.Set("currentUser", user)
        c.Next()
    }
}

// middleware/cors.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func CORS() gin.HandlerFunc {
    config := cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * 3600,
    }
    
    return cors.New(config)
}

// middleware/logger.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // リクエスト処理
        c.Next()
        
        // アクセスログ記録
        logger.Info("request",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", time.Since(start)),
            zap.String("ip", c.ClientIP()),
        )
    }
}
```

## ベストプラクティス
- ミドルウェアは単一責任を持つ
- エラー時は適切にAbort()を呼ぶ
- パフォーマンスへの影響を考慮
- 実行順序に注意（認証→認可→その他）