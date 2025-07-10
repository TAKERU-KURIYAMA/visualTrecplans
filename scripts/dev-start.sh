#!/bin/bash

# VisualTrecplans 開発環境起動スクリプト

set -e

# 色付きの出力用
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 VisualTrecplans 開発環境を起動しています...${NC}"

# Docker Composeが利用可能かチェック
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ docker-compose が見つかりません。インストールしてください。${NC}"
    exit 1
fi

# 既存のコンテナを停止
echo -e "${YELLOW}⏹️  既存のコンテナを停止しています...${NC}"
docker-compose down

# 最新のイメージでコンテナを起動
echo -e "${BLUE}🔄 コンテナを起動しています...${NC}"
docker-compose up -d

# サービスの起動を待機
echo -e "${YELLOW}⏳ サービスの起動を待機しています...${NC}"
sleep 10

# ヘルスチェック
echo -e "${BLUE}🏥 ヘルスチェックを実行しています...${NC}"

# PostgreSQL
echo -n "PostgreSQL: "
if docker-compose exec -T postgres pg_isready -U trecplans -d trecplans_dev >/dev/null 2>&1; then
    echo -e "${GREEN}✅ 正常${NC}"
else
    echo -e "${RED}❌ 接続エラー${NC}"
fi

# バックエンド（後で実装される予定）
echo -n "バックエンドAPI: "
if curl -f http://localhost:8080/health >/dev/null 2>&1; then
    echo -e "${GREEN}✅ 正常${NC}"
else
    echo -e "${YELLOW}⚠️  まだ起動中またはヘルスチェック未実装${NC}"
fi

# フロントエンド（後で実装される予定）
echo -n "フロントエンド: "
if curl -f http://localhost:3000 >/dev/null 2>&1; then
    echo -e "${GREEN}✅ 正常${NC}"
else
    echo -e "${YELLOW}⚠️  まだ起動中またはサーバー未実装${NC}"
fi

echo ""
echo -e "${GREEN}🎉 開発環境が起動しました！${NC}"
echo ""
echo -e "${BLUE}📋 アクセス情報:${NC}"
echo "🌐 フロントエンド: http://localhost:3000"
echo "🔗 バックエンドAPI: http://localhost:8080"
echo "🗄️  データベース管理: http://localhost:8081"
echo "   - システム: PostgreSQL"
echo "   - サーバー: postgres"
echo "   - ユーザー: trecplans"
echo "   - パスワード: password"
echo "   - データベース: trecplans_dev"
echo ""
echo -e "${BLUE}📋 便利なコマンド:${NC}"
echo "📊 ログを確認: docker-compose logs -f"
echo "🔍 特定のサービスのログ: docker-compose logs -f [frontend|backend|postgres]"
echo "🖥️  コンテナに入る: docker-compose exec [service-name] /bin/sh"
echo "⏹️  開発環境を停止: docker-compose down"
echo "🔄 データベースをリセット: ./scripts/dev-reset.sh"
echo ""
echo -e "${YELLOW}💡 開発中のファイル変更は自動で反映されます（ホットリロード）${NC}"

# ログの監視を開始するかどうか確認
echo ""
read -p "ログの監視を開始しますか？ (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}📊 ログの監視を開始します... (Ctrl+C で終了)${NC}"
    docker-compose logs -f
fi