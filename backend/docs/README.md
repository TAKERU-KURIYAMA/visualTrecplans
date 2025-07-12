# docs ディレクトリ

## 概要
バックエンドAPIドキュメントやアーキテクチャ設計書を配置するディレクトリです。
開発者向けの技術文書を管理します。

## ドキュメント構成例
```
docs/
├── api/              # API仕様書
│   ├── openapi.yaml  # OpenAPI仕様
│   └── postman/      # Postmanコレクション
├── architecture/     # アーキテクチャ設計書
│   ├── overview.md   # システム概要
│   └── database.md   # DB設計
├── guides/           # 開発ガイド
│   ├── setup.md      # セットアップガイド
│   └── testing.md    # テストガイド
└── adr/              # Architecture Decision Records
    ├── 001-use-gin-framework.md
    └── 002-adopt-clean-architecture.md
```

## API仕様書の例（OpenAPI）
```yaml
# docs/api/openapi.yaml
openapi: 3.0.0
info:
  title: VisualTrecplans API
  version: 1.0.0
  description: フィットネストラッキングアプリケーションのAPI

paths:
  /api/auth/login:
    post:
      summary: ユーザーログイン
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                email:
                  type: string
                password:
                  type: string
      responses:
        '200':
          description: ログイン成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string
                  user:
                    $ref: '#/components/schemas/User'
```

## ドキュメント管理方針
- APIドキュメントはコードと同期を保つ
- 重要な設計判断はADRとして記録
- 図表はPlantUMLやMermaidで管理