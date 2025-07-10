-- データベース初期化スクリプト
-- このファイルはPostgreSQLコンテナ起動時に自動実行されます

-- 拡張機能の有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 基本的なデータベース設定
ALTER DATABASE trecplans_dev SET timezone TO 'UTC';

-- 初期データの投入準備
-- テーブル作成はマイグレーションで行うため、ここでは基本設定のみ

-- 開発用ユーザーの作成（必要に応じて）
-- CREATE USER dev_user WITH PASSWORD 'dev_password';
-- GRANT ALL PRIVILEGES ON DATABASE trecplans_dev TO dev_user;

-- ログ設定
-- SET log_statement = 'all';
-- SET log_min_duration_statement = 0;

-- 開発用の設定
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;

-- 作業メモリの設定
-- SET work_mem = '16MB';

-- 初期化完了のログ
SELECT 'Database initialized successfully' AS status;