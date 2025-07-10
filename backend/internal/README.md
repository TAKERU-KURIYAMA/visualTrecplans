# internal ディレクトリ

## 概要
プライベートなアプリケーションコードを配置するディレクトリです。
Goの仕様により、このディレクトリ内のパッケージは外部プロジェクトからインポートできません。

## ディレクトリ構造
```
internal/
├── handlers/      # HTTPハンドラー（コントローラー層）
├── models/        # ドメインモデル・データベースモデル
├── services/      # ビジネスロジック層
├── database/      # データベース接続・クエリ管理
├── middleware/    # ミドルウェア（認証、ログ等）
└── validators/    # 入力値検証
```

## クリーンアーキテクチャの適用
- **handlers**: プレゼンテーション層（HTTPリクエスト/レスポンス処理）
- **services**: ユースケース層（ビジネスロジック）
- **models**: エンティティ層（ドメインモデル）
- **database**: インフラストラクチャ層（永続化）

## 依存関係の方向
```
handlers → services → models
    ↓         ↓
middleware  database
    ↓
validators
```

## 注意事項
- 各層は下位層にのみ依存する
- ビジネスロジックはservices層に集約
- handlers層はHTTP固有の処理のみ