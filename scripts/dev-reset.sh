#!/bin/bash

# VisualTrecplans 開発環境リセットスクリプト

set -e

# 色付きの出力用
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${YELLOW}⚠️  開発環境をリセットします${NC}"
echo -e "${RED}⚠️  この操作により、すべてのデータが失われます！${NC}"
echo ""

# 確認プロンプト
read -p "本当にリセットしますか？ (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}ℹ️  リセットをキャンセルしました${NC}"
    exit 0
fi

echo -e "${BLUE}🔄 開発環境をリセットしています...${NC}"

# すべてのコンテナを停止
echo -e "${YELLOW}⏹️  コンテナを停止しています...${NC}"
docker-compose down

# ボリュームの削除
echo -e "${YELLOW}🗑️  データベースボリュームを削除しています...${NC}"
docker-compose down -v

# 未使用のDockerリソースをクリーンアップ
echo -e "${YELLOW}🧹 未使用のDockerリソースをクリーンアップしています...${NC}"
docker system prune -f

# イメージの再ビルド
echo -e "${BLUE}🏗️  イメージを再ビルドしています...${NC}"
docker-compose build --no-cache

# データベースの再起動と初期化
echo -e "${BLUE}🗄️  データベースを再起動しています...${NC}"
docker-compose up -d postgres

# PostgreSQLが起動するまで待機
echo -e "${YELLOW}⏳ PostgreSQLの起動を待機しています...${NC}"
sleep 30

# マイグレーションの実行（後で実装される予定）
echo -e "${BLUE}🔄 マイグレーションを実行しています...${NC}"
# docker-compose exec backend make migrate-up

# 開発用データの投入（後で実装される予定）
echo -e "${BLUE}🌱 開発用データを投入しています...${NC}"
# docker-compose exec backend make seed

# 全サービスの起動
echo -e "${BLUE}🚀 全サービスを起動しています...${NC}"
docker-compose up -d

# 起動確認
echo -e "${YELLOW}⏳ サービスの起動を確認しています...${NC}"
sleep 15

# ヘルスチェック
echo -e "${BLUE}🏥 ヘルスチェックを実行しています...${NC}"

# PostgreSQL
echo -n "PostgreSQL: "
if docker-compose exec -T postgres pg_isready -U trecplans -d trecplans_dev >/dev/null 2>&1; then
    echo -e "${GREEN}✅ 正常${NC}"
else
    echo -e "${RED}❌ 接続エラー${NC}"
fi

# バックエンド
echo -n "バックエンドAPI: "
if curl -f http://localhost:8080/health >/dev/null 2>&1; then
    echo -e "${GREEN}✅ 正常${NC}"
else
    echo -e "${YELLOW}⚠️  まだ起動中またはヘルスチェック未実装${NC}"
fi

# フロントエンド
echo -n "フロントエンド: "
if curl -f http://localhost:3000 >/dev/null 2>&1; then
    echo -e "${GREEN}✅ 正常${NC}"
else
    echo -e "${YELLOW}⚠️  まだ起動中またはサーバー未実装${NC}"
fi

echo ""
echo -e "${GREEN}🎉 開発環境のリセットが完了しました！${NC}"
echo ""
echo -e "${BLUE}📋 アクセス情報:${NC}"
echo "🌐 フロントエンド: http://localhost:3000"
echo "🔗 バックエンドAPI: http://localhost:8080"
echo "🗄️  データベース管理: http://localhost:8081"
echo ""
echo -e "${BLUE}📋 次のステップ:${NC}"
echo "1. 開発を開始: ./scripts/dev-start.sh"
echo "2. ログを確認: docker-compose logs -f"
echo "3. 開発を停止: docker-compose down"