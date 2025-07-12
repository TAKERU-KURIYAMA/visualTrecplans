# cmd ディレクトリ

## 概要
アプリケーションのエントリーポイントを配置するディレクトリです。
main.goファイルや、複数のサービスがある場合は各サービスのエントリーポイントを管理します。

## ディレクトリ構造例
```
cmd/
├── api/
│   └── main.go      # APIサーバーのエントリーポイント
├── migration/
│   └── main.go      # マイグレーション実行用
└── seeder/
    └── main.go      # シードデータ投入用
```

## 命名規則
- 各サービスのディレクトリ名は小文字
- main.goをエントリーポイントとする

## 例
```go
// cmd/api/main.go
package main

import (
    "log"
    "github.com/visualtrecplans/backend/internal/handlers"
    "github.com/visualtrecplans/backend/internal/database"
)

func main() {
    // データベース接続
    db := database.Connect()
    defer db.Close()
    
    // サーバー起動
    router := handlers.SetupRouter(db)
    log.Fatal(router.Run(":8080"))
}
```