# migrations ディレクトリ

## 概要
データベースマイグレーションファイルを配置するディレクトリです。
スキーマの変更履歴を管理し、データベースのバージョン管理を行います。

## ファイル命名規則
```
{タイムスタンプ}_{説明}.{up|down}.sql

例:
20250110_001_create_users_table.up.sql
20250110_001_create_users_table.down.sql
```

## ディレクトリ構造例
```
migrations/
├── 20250110_001_create_users_table.up.sql
├── 20250110_001_create_users_table.down.sql
├── 20250110_002_create_workouts_table.up.sql
├── 20250110_002_create_workouts_table.down.sql
├── 20250111_001_add_user_preferences.up.sql
└── 20250111_001_add_user_preferences.down.sql
```

## マイグレーションファイル例
```sql
-- 20250110_001_create_users_table.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- 20250110_001_create_users_table.down.sql
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

## マイグレーションツール
- [golang-migrate](https://github.com/golang-migrate/migrate)を使用
- `make migrate-up`でマイグレーション実行
- `make migrate-down`でロールバック

## 注意事項
- 各マイグレーションはアトミックに実行
- downマイグレーションは必ず用意
- 本番環境でのロールバック可能性を考慮