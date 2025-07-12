# 環境変数設定ガイド

## 概要
VisualTrecplansプロジェクトでは、フロントエンドとバックエンドで環境変数を使用して設定を管理しています。

## フロントエンド環境変数

### 必須変数
| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| `VITE_API_URL` | バックエンドAPI URL | `http://localhost:8080/api/v1` |
| `VITE_APP_NAME` | アプリケーション名 | `VisualTrecplans` |
| `VITE_APP_VERSION` | アプリケーションバージョン | `1.0.0` |

### 任意変数
| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| `VITE_APP_TITLE` | ブラウザタイトル | `VisualTrecplans` |
| `VITE_APP_DESCRIPTION` | アプリケーション説明 | - |
| `VITE_FEATURE_SOCIAL` | ソーシャル機能有効化 | `true` |
| `VITE_FEATURE_EXPORT` | エクスポート機能有効化 | `true` |
| `VITE_FEATURE_ANALYTICS` | アナリティクス機能有効化 | `true` |
| `VITE_GOOGLE_ANALYTICS_ID` | Google Analytics ID | - |
| `VITE_SENTRY_DSN` | Sentry DSN | - |
| `VITE_DEBUG` | デバッグモード | `false` |
| `VITE_LOG_LEVEL` | ログレベル | `info` |

## バックエンド環境変数

### アプリケーション設定
| 変数名 | 説明 | デフォルト値 | 必須 |
|--------|------|-------------|------|
| `APP_NAME` | アプリケーション名 | `VisualTrecplans` | ❌ |
| `APP_ENV` | 環境名 | `development` | ❌ |
| `APP_VERSION` | バージョン | `1.0.0` | ❌ |
| `APP_PORT` | サーバーポート | `8080` | ❌ |
| `APP_DEBUG` | デバッグモード | `true` | ❌ |

### データベース設定
| 変数名 | 説明 | デフォルト値 | 必須 |
|--------|------|-------------|------|
| `DB_HOST` | データベースホスト | `localhost` | ✅ |
| `DB_PORT` | データベースポート | `5432` | ❌ |
| `DB_USER` | データベースユーザー | `trecplans` | ✅ |
| `DB_PASSWORD` | データベースパスワード | `password` | ❌ |
| `DB_NAME` | データベース名 | `trecplans_dev` | ✅ |
| `DB_SSLMODE` | SSL接続モード | `disable` | ❌ |

### JWT設定
| 変数名 | 説明 | デフォルト値 | 必須 |
|--------|------|-------------|------|
| `JWT_SECRET` | JWT秘密鍵 | - | ✅(本番環境) |
| `JWT_EXPIRY` | JWT有効期限 | `24h` | ❌ |
| `JWT_REFRESH_EXPIRY` | リフレッシュトークン有効期限 | `168h` | ❌ |

### CORS設定
| 変数名 | 説明 | デフォルト値 | 必須 |
|--------|------|-------------|------|
| `CORS_ALLOWED_ORIGINS` | 許可するオリジン（カンマ区切り） | `http://localhost:5173,http://localhost:3000` | ❌ |
| `CORS_ALLOWED_METHODS` | 許可するHTTPメソッド | `GET,POST,PUT,DELETE,OPTIONS` | ❌ |
| `CORS_ALLOWED_HEADERS` | 許可するヘッダー | `Origin,Content-Type,Accept,Authorization` | ❌ |

### その他
| 変数名 | 説明 | デフォルト値 | 必須 |
|--------|------|-------------|------|
| `LOG_LEVEL` | ログレベル | `debug` | ❌ |
| `LOG_FORMAT` | ログフォーマット | `json` | ❌ |
| `UPLOAD_PATH` | アップロードディレクトリ | `./uploads` | ❌ |
| `MAX_UPLOAD_SIZE` | 最大アップロードサイズ（バイト） | `10485760` | ❌ |

## 環境別設定ファイル

### フロントエンド
- `.env.development` - 開発環境用設定
- `.env.production` - 本番環境用設定
- `.env.example` - 設定テンプレート

### バックエンド
- `.env.example` - 設定テンプレート
- `.env.docker` - Docker環境用設定（開発）

## セットアップ手順

### 1. フロントエンド
```bash
cd frontend
cp .env.example .env.local
# .env.localを編集して必要な値を設定
```

### 2. バックエンド
```bash
cd backend
cp .env.example .env
# .envを編集して必要な値を設定
```

### 3. Docker環境
```bash
cp .env.docker.example .env.docker
# .env.dockerを編集して必要な値を設定
```

## セキュリティ注意事項

### 本番環境
1. **JWT_SECRET**: 強力なランダムキーを使用
2. **データベースパスワード**: 複雑なパスワードを設定
3. **CORS_ALLOWED_ORIGINS**: 実際のドメインのみを許可
4. **デバッグモード**: 本番環境では無効化

### 環境変数ファイルの管理
- `.env`ファイルはGitにコミットしない
- 本番環境では環境変数を直接設定
- シークレット情報は暗号化して管理

## トラブルシューティング

### よくある問題

#### 1. API接続エラー
```
Error: Failed to fetch from API
```
**解決方法**: `VITE_API_URL`が正しく設定されているか確認

#### 2. JWT認証エラー
```
Error: Invalid JWT token
```
**解決方法**: `JWT_SECRET`がフロントエンドとバックエンドで一致しているか確認

#### 3. データベース接続エラー
```
Error: Failed to connect to database
```
**解決方法**: データベース設定（`DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`）を確認

#### 4. CORS エラー
```
Error: CORS policy error
```
**解決方法**: `CORS_ALLOWED_ORIGINS`にフロントエンドのURLが含まれているか確認

### デバッグ方法

#### 1. 設定値の確認
```bash
# バックエンド
make run
# ログで設定値を確認

# フロントエンド
npm run dev
# ブラウザの開発者ツールでconsole.log(import.meta.env)
```

#### 2. 環境変数の表示
```bash
# Linux/Mac
printenv | grep VITE_
printenv | grep DB_

# Windows
set | findstr VITE_
set | findstr DB_
```

## 参考リンク

- [Vite 環境変数](https://vitejs.dev/guide/env-and-mode.html)
- [Gin Framework](https://gin-gonic.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [JWT](https://jwt.io/)