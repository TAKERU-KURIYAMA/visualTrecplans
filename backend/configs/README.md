# configs ディレクトリ

## 概要
アプリケーションの設定ファイルを配置するディレクトリです。
環境別の設定やデフォルト設定を管理します。

## ファイル構成例
```
configs/
├── config.yaml         # デフォルト設定
├── config.dev.yaml     # 開発環境設定
├── config.prod.yaml    # 本番環境設定
└── config.test.yaml    # テスト環境設定
```

## 設定項目例
```yaml
# config.yaml
app:
  name: "VisualTrecplans"
  version: "1.0.0"
  port: 8080

database:
  host: "localhost"
  port: 5432
  name: "visualtrecplans"
  sslmode: "disable"

jwt:
  secret: "your-secret-key"
  expiry: 24h

cors:
  allowed_origins:
    - "http://localhost:3000"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
  allowed_headers:
    - "Content-Type"
    - "Authorization"
```

## 環境変数との連携
- 環境変数で上書き可能な設定
- 例: `DB_HOST`環境変数で`database.host`を上書き

## セキュリティ注意事項
- 秘密情報は環境変数で管理
- config.*.yamlはGitにコミットしない（.gitignoreに追加）
- config.example.yamlを用意してサンプルを提供