#!/bin/bash

# VisualTrecplans 開発環境セットアップスクリプト

set -e

echo "🚀 VisualTrecplans 開発環境セットアップを開始します..."

# 色付きの出力用
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 必要なコマンドのチェック
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}❌ $1 が見つかりません。インストールしてください。${NC}"
        exit 1
    fi
}

echo -e "${BLUE}📋 必要なコマンドをチェックしています...${NC}"
check_command "docker"
check_command "docker-compose"
check_command "git"

echo -e "${GREEN}✅ 必要なコマンドが揃っています${NC}"

# 環境変数ファイルの作成
echo -e "${BLUE}📝 環境変数ファイルを作成しています...${NC}"

if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  .env ファイルが見つかりません。サンプルからコピーします${NC}"
    
    # バックエンドの.envファイル
    if [ -f backend/.env.example ]; then
        cp backend/.env.example backend/.env
        echo -e "${GREEN}✅ backend/.env を作成しました${NC}"
    fi
    
    # フロントエンドの.envファイル（必要に応じて）
    if [ -f frontend/.env.example ]; then
        cp frontend/.env.example frontend/.env
        echo -e "${GREEN}✅ frontend/.env を作成しました${NC}"
    fi
else
    echo -e "${GREEN}✅ .env ファイルが既に存在します${NC}"
fi

# Dockerイメージのビルド
echo -e "${BLUE}🏗️  Dockerイメージをビルドしています...${NC}"
docker-compose build

# データベースの初期化
echo -e "${BLUE}🗄️  データベースを初期化しています...${NC}"
docker-compose up -d postgres

# PostgreSQLが起動するまで待機
echo -e "${YELLOW}⏳ PostgreSQLの起動を待機しています...${NC}"
sleep 30

# マイグレーションの実行（後で実装される予定）
echo -e "${BLUE}🔄 マイグレーションを実行します...${NC}"
# docker-compose exec backend make migrate-up

# 依存関係のインストール
echo -e "${BLUE}📦 依存関係をインストールしています...${NC}"

# フロントエンドの依存関係（後で実装される予定）
# if [ -f frontend/package.json ]; then
#     echo -e "${BLUE}📦 フロントエンドの依存関係をインストール中...${NC}"
#     docker-compose run --rm frontend npm install
# fi

# バックエンドの依存関係
if [ -f backend/go.mod ]; then
    echo -e "${BLUE}📦 バックエンドの依存関係をダウンロード中...${NC}"
    docker-compose run --rm backend go mod download
fi

# 開発用データの投入（後で実装される予定）
echo -e "${BLUE}🌱 開発用データを投入しています...${NC}"
# docker-compose exec backend make seed

echo -e "${GREEN}🎉 開発環境のセットアップが完了しました！${NC}"
echo ""
echo -e "${BLUE}📋 次のステップ:${NC}"
echo "1. 開発環境を起動: ./scripts/dev-start.sh"
echo "2. フロントエンド: http://localhost:3000"
echo "3. バックエンドAPI: http://localhost:8080"
echo "4. データベース管理: http://localhost:8081 (adminer)"
echo "5. 開発を停止: docker-compose down"
echo ""
echo -e "${YELLOW}💡 ヒント:${NC}"
echo "- ログを確認: docker-compose logs -f [service-name]"
echo "- コンテナに入る: docker-compose exec [service-name] /bin/sh"
echo "- データベースをリセット: ./scripts/dev-reset.sh"