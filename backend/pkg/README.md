# pkg ディレクトリ

## 概要
再利用可能な公開パッケージを配置するディレクトリです。
他のプロジェクトからインポート可能なユーティリティや共通機能を管理します。

## 想定されるパッケージ例
```
pkg/
├── logger/        # ロギングユーティリティ
├── errors/        # カスタムエラー定義
├── response/      # 共通レスポンス形式
├── auth/          # 認証関連ユーティリティ
└── utils/         # その他の汎用ユーティリティ
```

## 命名規則
- パッケージ名は小文字
- 機能を表す簡潔な名前を使用

## 例
```go
// pkg/response/response.go
package response

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func Success(data interface{}) APIResponse {
    return APIResponse{
        Success: true,
        Data:    data,
    }
}

func Error(err string) APIResponse {
    return APIResponse{
        Success: false,
        Error:   err,
    }
}
```

## 注意事項
- アプリケーション固有のロジックは含めない
- 汎用性の高い機能のみを配置
- 十分なテストカバレッジを確保