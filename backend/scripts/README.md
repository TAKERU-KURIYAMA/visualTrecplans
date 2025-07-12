# scripts ディレクトリ

## 概要
ビルド、デプロイ、メンテナンスなどの自動化スクリプトを配置するディレクトリです。
開発効率を向上させるヘルパースクリプトを管理します。

## スクリプト例
```
scripts/
├── build.sh          # ビルドスクリプト
├── deploy.sh         # デプロイスクリプト
├── test.sh           # テスト実行スクリプト
├── setup-dev.sh      # 開発環境セットアップ
├── generate-docs.sh  # ドキュメント生成
└── backup-db.sh      # データベースバックアップ
```

## スクリプト例：setup-dev.sh
```bash
#!/bin/bash
set -e

echo "開発環境をセットアップしています..."

# 依存関係のインストール
echo "Go依存関係をインストール中..."
go mod download

# 環境変数ファイルのコピー
if [ ! -f .env ]; then
    echo ".envファイルを作成中..."
    cp .env.example .env
fi

# データベースの起動
echo "PostgreSQLを起動中..."
docker-compose up -d postgres

# マイグレーションの実行
echo "データベースマイグレーションを実行中..."
make migrate-up

# シードデータの投入
echo "シードデータを投入中..."
go run cmd/seeder/main.go

echo "セットアップが完了しました！"
```

## 命名規則
- スクリプト名は小文字とハイフンを使用
- 機能を明確に表す名前を付ける
- 実行権限を付与（chmod +x）

## 注意事項
- エラーハンドリングを適切に実装
- 冪等性を保つ（複数回実行しても安全）
- 必要に応じてヘルプメッセージを表示