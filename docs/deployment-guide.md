# GitHub Actions 自動デプロイ機能 使い方ガイド

## 概要
VisualTrecplansプロジェクトでは、GitHub Actionsを使用したCI/CDパイプラインにより、コードのプッシュからデプロイまでを自動化しています。

## 🚀 デプロイフロー

### 自動デプロイの流れ
```
コード変更 → Push/PR → CI/CD実行 → 自動デプロイ
    ↓
1. バックエンドテスト・リント
2. フロントエンドテスト・リント  
3. セキュリティスキャン
4. Dockerイメージビルド
5. 環境別デプロイ
6. E2Eテスト
7. Slack通知
```

### ブランチ別デプロイ戦略
- **develop** → 開発環境 (dev.trecplans.com)
- **main** → 本番環境 (trecplans.com)
- **PR** → テスト実行のみ（デプロイなし）

## ⚙️ セットアップ手順

### 1. GitHub Secrets設定
リポジトリの Settings → Secrets and variables → Actions で以下を設定：

#### サーバー接続情報
```
DEV_HOST=dev.example.com
DEV_USERNAME=deploy
DEV_SSH_KEY=-----BEGIN OPENSSH PRIVATE KEY-----...

PROD_HOST=prod.example.com  
PROD_USERNAME=deploy
PROD_SSH_KEY=-----BEGIN OPENSSH PRIVATE KEY-----...
```

#### データベース設定
```
DB_USER=trecplans
DB_PASSWORD=secure_password
DB_NAME=trecplans
```

#### JWT設定
```
JWT_SECRET=your_jwt_secret_key
JWT_REFRESH_SECRET=your_refresh_secret_key
```

#### 外部サービス
```
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
REDIS_PASSWORD=redis_password
GRAFANA_PASSWORD=grafana_admin_password
```

### 2. GitHub Variables設定
Variables タブで以下を設定：

```
VITE_API_URL=https://api.trecplans.com/api/v1
CORS_ALLOWED_ORIGINS=https://trecplans.com,https://dev.trecplans.com
GITHUB_REPOSITORY=your-username/visualtrecplans
```

### 3. サーバー環境準備

#### Docker & Docker Compose インストール
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### デプロイディレクトリ作成
```bash
sudo mkdir -p /opt/trecplans
sudo chown $USER:$USER /opt/trecplans
cd /opt/trecplans

# 設定ファイル配置
git clone https://github.com/your-username/visualtrecplans.git .
cp .env.example .env
# .envファイルを編集して本番設定を入力
```

#### SSH鍵設定
```bash
# デプロイ用ユーザー作成
sudo useradd -m -s /bin/bash deploy
sudo usermod -aG docker deploy

# SSH鍵設定
sudo mkdir -p /home/deploy/.ssh
sudo touch /home/deploy/.ssh/authorized_keys
sudo chown -R deploy:deploy /home/deploy/.ssh
sudo chmod 700 /home/deploy/.ssh
sudo chmod 600 /home/deploy/.ssh/authorized_keys

# 公開鍵を authorized_keys に追加
echo "ssh-rsa YOUR_PUBLIC_KEY deploy@trecplans" | sudo tee -a /home/deploy/.ssh/authorized_keys
```

## 🔄 使用方法

### 開発環境デプロイ
```bash
# developブランチにプッシュで自動デプロイ
git checkout develop
git add .
git commit -m "feat: new feature"
git push origin develop

# → 開発環境 (dev.trecplans.com) に自動デプロイ
```

### 本番環境デプロイ
```bash
# mainブランチにマージで自動デプロイ
git checkout main
git merge develop
git push origin main

# → 本番環境 (trecplans.com) に自動デプロイ
```

### プルリクエスト時の動作
```bash
# feature ブランチからPR作成
git checkout -b feature/new-feature
git add .
git commit -m "feat: implement new feature"
git push origin feature/new-feature

# → テスト・リント実行（デプロイなし）
# → PRマージ後に対象環境へデプロイ
```

## 🔧 手動デプロイ・管理

### 手動ワークフロー実行
GitHub Actions タブから手動実行可能：

1. **ロールバック実行**
   - Actions → "Rollback Deployment" → "Run workflow"
   - 環境選択（development/production）
   - コミットSHA指定（任意）

2. **個別コンポーネントデプロイ**
   - 必要に応じてカスタムワークフローを作成

### ローカルからの手動デプロイ
```bash
# 開発環境
ssh deploy@dev.example.com
cd /opt/trecplans
docker-compose -f docker-compose.dev.yml pull
docker-compose -f docker-compose.dev.yml up -d

# 本番環境
ssh deploy@prod.example.com
cd /opt/trecplans
docker-compose pull
docker-compose up -d
```

## 📊 監視・ログ

### デプロイ状況確認
- **GitHub Actions**: リポジトリの Actions タブ
- **Slack通知**: #deployments チャンネル
- **Grafana**: https://monitoring.trecplans.com:3001

### ログ確認方法
```bash
# アプリケーションログ
docker-compose logs -f backend
docker-compose logs -f frontend

# システムログ
journalctl -u docker -f

# Nginx アクセスログ
docker-compose exec frontend tail -f /var/log/nginx/access.log
```

### ヘルスチェック
```bash
# APIヘルスチェック
curl https://api.trecplans.com/health

# フロントエンドヘルスチェック  
curl https://trecplans.com/health

# データベース接続確認
docker-compose exec postgres pg_isready -U trecplans
```

## 🔄 ロールバック手順

### 自動ロールバック（推奨）
1. GitHub → Actions → "Rollback Deployment"
2. "Run workflow" クリック
3. 環境選択（production/development）
4. ロールバック先コミットSHA入力（任意）
5. "Run workflow" 実行

### 手動ロールバック
```bash
# サーバーにSSH接続
ssh deploy@prod.example.com
cd /opt/trecplans

# 現在のバージョンをバックアップ
docker tag ghcr.io/your-repo/backend:latest ghcr.io/your-repo/backend:backup-$(date +%Y%m%d_%H%M%S)

# 以前のバージョンを取得
docker pull ghcr.io/your-repo/backend:main-COMMIT_SHA
docker tag ghcr.io/your-repo/backend:main-COMMIT_SHA ghcr.io/your-repo/backend:latest

# 再デプロイ
docker-compose down
docker-compose up -d

# ヘルスチェック
curl -f https://api.trecplans.com/health
```

## 🔒 セキュリティ機能

### 自動セキュリティチェック
- **依存関係脆弱性スキャン**: npm audit, govulncheck
- **コンテナ脆弱性スキャン**: Trivy
- **コード品質チェック**: ESLint, golangci-lint
- **SAST/DAST**: CodeQL integration

### セキュリティアラート対応
```bash
# 脆弱性発見時の対応フロー
1. GitHub Security Advisory確認
2. 依存関係アップデート
3. テスト実行・動作確認
4. セキュリティパッチデプロイ
5. 本番環境動作確認
```

## 🚨 トラブルシューティング

### よくある問題と解決法

#### デプロイ失敗時
```bash
# 1. ワークフローログ確認
GitHub Actions → 失敗したジョブクリック → ログ詳細確認

# 2. サーバー状態確認
ssh deploy@server
docker-compose ps
docker-compose logs

# 3. ディスク容量確認
df -h
docker system prune -a
```

#### データベース接続エラー
```bash
# PostgreSQL状態確認
docker-compose exec postgres pg_isready
docker-compose logs postgres

# 接続テスト
docker-compose exec backend go run cmd/migrate/main.go --dry-run
```

#### SSL証明書エラー
```bash
# Let's Encrypt証明書更新
certbot renew --dry-run
docker-compose restart frontend
```

### エラー別対応表

| エラー | 原因 | 解決方法 |
|--------|------|----------|
| Docker image pull失敗 | レジストリ認証エラー | GitHub Container Registry設定確認 |
| Database migration失敗 | スキーマ競合 | 手動マイグレーション実行 |
| Health check失敗 | アプリケーション起動遅延 | タイムアウト時間延長 |
| Disk space不足 | 古いイメージ蓄積 | `docker system prune -a` 実行 |

## 📈 パフォーマンス最適化

### ビルド時間短縮
- **並列実行**: テスト・ビルドジョブ並列化
- **キャッシュ活用**: Docker layer cache, npm cache
- **差分ビルド**: 変更されたコンポーネントのみビルド

### デプロイ時間短縮
- **Rolling Update**: ゼロダウンタイムデプロイ
- **Health Check**: 適切なタイムアウト設定
- **Container Registry**: 最適化されたイメージ使用

## 📞 サポート・問い合わせ

### 緊急時連絡先
- **Slack**: #ops-emergency
- **Email**: ops@trecplans.com
- **On-call**: PagerDuty integration

### ドキュメント参照
- **API仕様**: `/docs/api-specification.md`
- **システム設計**: `/docs/architecture.md`
- **開発ガイド**: `/docs/development-guide.md`

---

このガイドにより、開発チームはGitHub Actionsを活用した効率的なCI/CDパイプラインを運用できます。定期的な見直しと改善により、より安定したデプロイ環境を維持しましょう。