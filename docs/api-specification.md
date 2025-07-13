# API仕様書

## 概要
VisualTrecplans APIは、フィットネス・トレーニング管理アプリケーションのバックエンドAPIです。
認証機能、ワークアウト記録機能、マスタデータ管理機能を提供します。

## ベースURL
- 開発環境: `http://localhost:8080/api/v1`
- 本番環境: `https://api.trecplans.com/api/v1`

## 認証
すべての保護されたエンドポイントには、Authorization ヘッダーにJWTトークンが必要です。

```
Authorization: Bearer <JWT_TOKEN>
```

## 共通レスポンス形式

### 成功レスポンス
```json
{
  "data": {},
  "message": "success"
}
```

### エラーレスポンス
```json
{
  "error": "エラーメッセージ",
  "code": "ERROR_CODE",
  "details": {}
}
```

## エンドポイント一覧

### 認証 (Authentication)

#### POST /auth/register
ユーザー登録

**リクエスト**
```json
{
  "username": "string (3-50文字)",
  "email": "string (有効なメールアドレス)",
  "password": "string (8文字以上、英数字記号を含む)"
}
```

**レスポンス (201)**
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "created_at": "datetime"
}
```

#### POST /auth/login
ログイン

**リクエスト**
```json
{
  "email": "string",
  "password": "string"
}
```

**レスポンス (200)**
```json
{
  "access_token": "string",
  "refresh_token": "string",
  "user": {
    "id": "uuid",
    "username": "string",
    "email": "string"
  }
}
```

#### POST /auth/refresh
トークンリフレッシュ

**リクエスト**
```json
{
  "refresh_token": "string"
}
```

**レスポンス (200)**
```json
{
  "access_token": "string",
  "refresh_token": "string"
}
```

#### POST /auth/logout
ログアウト

**ヘッダー**: `Authorization: Bearer <token>`

**レスポンス (200)**
```json
{
  "message": "Logged out successfully"
}
```

#### GET /auth/profile
プロフィール取得

**ヘッダー**: `Authorization: Bearer <token>`

**レスポンス (200)**
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

#### PUT /auth/profile
プロフィール更新

**ヘッダー**: `Authorization: Bearer <token>`

**リクエスト**
```json
{
  "username": "string (optional)",
  "email": "string (optional)"
}
```

**レスポンス (200)**
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "updated_at": "datetime"
}
```

### ワークアウト (Workouts)

#### POST /workouts
ワークアウト作成

**ヘッダー**: `Authorization: Bearer <token>`

**リクエスト**
```json
{
  "muscle_group": "string (chest|back|shoulders|arms|core|legs|glutes|full_body)",
  "exercise_name": "string (1-100文字)",
  "exercise_icon": "string (optional, max 50文字)",
  "weight_kg": "number (optional, 0-999.99)",
  "reps": "integer (optional, 1-999)",
  "sets": "integer (optional, 1-99)",
  "notes": "string (optional, max 500文字)",
  "performed_at": "datetime (ISO8601)"
}
```

**レスポンス (201)**
```json
{
  "id": "uuid",
  "muscle_group": "string",
  "exercise_name": "string",
  "exercise_icon": "string",
  "weight_kg": "number",
  "reps": "integer",
  "sets": "integer",
  "notes": "string",
  "performed_at": "datetime",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

#### GET /workouts
ワークアウト一覧取得

**ヘッダー**: `Authorization: Bearer <token>`

**クエリパラメータ**
- `muscle_group`: string (optional) - 筋肉部位フィルター
- `start_date`: string (optional) - 開始日フィルター (YYYY-MM-DD)
- `end_date`: string (optional) - 終了日フィルター (YYYY-MM-DD)
- `exercise_name`: string (optional) - エクササイズ名検索
- `page`: integer (default: 1) - ページ番号
- `per_page`: integer (default: 20, max: 100) - 1ページあたりの件数
- `order_by`: string (default: performed_at) - ソート項目
- `order`: string (default: desc) - ソート順 (asc|desc)

**レスポンス (200)**
```json
{
  "workouts": [
    {
      "id": "uuid",
      "muscle_group": "string",
      "exercise_name": "string",
      "exercise_icon": "string",
      "weight_kg": "number",
      "reps": "integer",
      "sets": "integer",
      "notes": "string",
      "performed_at": "datetime",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  ],
  "total": "integer",
  "page": "integer",
  "per_page": "integer",
  "total_pages": "integer"
}
```

#### GET /workouts/:id
ワークアウト詳細取得

**ヘッダー**: `Authorization: Bearer <token>`

**レスポンス (200)**
```json
{
  "id": "uuid",
  "muscle_group": "string",
  "exercise_name": "string",
  "exercise_icon": "string",
  "weight_kg": "number",
  "reps": "integer",
  "sets": "integer",
  "notes": "string",
  "performed_at": "datetime",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

#### PUT /workouts/:id
ワークアウト更新

**ヘッダー**: `Authorization: Bearer <token>`

**リクエスト**
```json
{
  "muscle_group": "string (optional)",
  "exercise_name": "string (optional)",
  "exercise_icon": "string (optional)",
  "weight_kg": "number (optional)",
  "reps": "integer (optional)",
  "sets": "integer (optional)",
  "notes": "string (optional)",
  "performed_at": "datetime (optional)"
}
```

**レスポンス (200)**
```json
{
  "id": "uuid",
  "muscle_group": "string",
  "exercise_name": "string",
  "exercise_icon": "string",
  "weight_kg": "number",
  "reps": "integer",
  "sets": "integer",
  "notes": "string",
  "performed_at": "datetime",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

#### DELETE /workouts/:id
ワークアウト削除

**ヘッダー**: `Authorization: Bearer <token>`

**レスポンス (204)**
レスポンスボディなし

#### GET /workouts/stats
ワークアウト統計取得

**ヘッダー**: `Authorization: Bearer <token>`

**クエリパラメータ**
- `period`: string (default: month) - 集計期間 (week|month|year)

**レスポンス (200)**
```json
{
  "total_workouts": "integer",
  "total_sets": "integer", 
  "total_reps": "integer",
  "total_weight_lifted": "number",
  "workouts_by_muscle": {
    "chest": "integer",
    "back": "integer",
    "legs": "integer"
  },
  "most_used_exercises": [
    {
      "exercise_name": "string",
      "count": "integer"
    }
  ],
  "recent_workouts": [
    {
      "id": "uuid",
      "exercise_name": "string",
      "performed_at": "datetime"
    }
  ],
  "weekly_progress": [
    {
      "week": "string",
      "workout_count": "integer",
      "total_sets": "integer",
      "total_reps": "integer",
      "total_weight": "number"
    }
  ]
}
```

### マスタデータ (Master Data)

#### GET /muscle-groups
筋肉部位一覧取得

**クエリパラメータ**
- `lang`: string (default: ja) - 言語 (ja|en)
- `category`: string (optional) - カテゴリフィルター (upper|lower|core|full_body)

**レスポンス (200)**
```json
{
  "data": [
    {
      "code": "string",
      "name_ja": "string",
      "name_en": "string", 
      "category": "string",
      "color_code": "string",
      "sort_order": "integer"
    }
  ]
}
```

#### GET /exercises
エクササイズ一覧取得

**クエリパラメータ**
- `muscle_group`: string (optional) - 筋肉部位フィルター
- `lang`: string (default: ja) - 言語 (ja|en)

**レスポンス (200)**
```json
{
  "data": [
    {
      "id": "integer",
      "muscle_group_code": "string",
      "name_ja": "string",
      "name_en": "string",
      "icon_name": "string",
      "is_custom": "boolean",
      "sort_order": "integer"
    }
  ]
}
```

#### GET /exercises/custom
カスタムエクササイズ一覧取得

**ヘッダー**: `Authorization: Bearer <token>`

**レスポンス (200)**
```json
{
  "data": [
    {
      "id": "integer",
      "muscle_group_code": "string",
      "name_ja": "string", 
      "name_en": "string",
      "icon_name": "string",
      "is_custom": true,
      "sort_order": "integer"
    }
  ]
}
```

#### POST /exercises/custom
カスタムエクササイズ作成

**ヘッダー**: `Authorization: Bearer <token>`

**リクエスト**
```json
{
  "name": "string (1-100文字)",
  "muscle_group_code": "string (chest|back|shoulders|arms|core|legs|glutes|full_body)",
  "icon_name": "string (optional, max 50文字)"
}
```

**レスポンス (201)**
```json
{
  "id": "integer",
  "muscle_group_code": "string",
  "name_ja": "string",
  "name_en": "string", 
  "icon_name": "string",
  "is_custom": true,
  "sort_order": "integer"
}
```

#### GET /exercise-icons
エクササイズアイコン一覧取得

**クエリパラメータ**
- `category`: string (optional) - カテゴリフィルター

**レスポンス (200)**
```json
{
  "data": [
    {
      "name": "string",
      "svg_path": "string",
      "category": "string"
    }
  ]
}
```

## HTTPステータスコード

| ステータス | 説明 |
|----------|------|
| 200 | 成功 |
| 201 | 作成成功 |
| 204 | 削除成功 |
| 400 | 不正なリクエスト |
| 401 | 認証が必要 |
| 403 | アクセス権限なし |
| 404 | リソースが見つからない |
| 409 | 競合 |
| 422 | バリデーションエラー |
| 429 | レート制限超過 |
| 500 | サーバーエラー |

## エラーコード

| コード | 説明 |
|--------|------|
| AUTH_REQUIRED | 認証が必要 |
| INVALID_TOKEN | 無効なトークン |
| TOKEN_EXPIRED | トークンの有効期限切れ |
| INVALID_REQUEST | 不正なリクエスト |
| VALIDATION_ERROR | バリデーションエラー |
| INVALID_CREDENTIALS | 認証情報が不正 |
| USER_EXISTS | ユーザーが既に存在 |
| USER_NOT_FOUND | ユーザーが見つからない |
| WORKOUT_NOT_FOUND | ワークアウトが見つからない |
| INVALID_MUSCLE_GROUP | 無効な筋肉部位 |
| ACCESS_DENIED | アクセス拒否 |
| RATE_LIMIT_EXCEEDED | レート制限超過 |
| INTERNAL_ERROR | サーバー内部エラー |

## レート制限

| エンドポイント | 制限 |
|---------------|------|
| 一般API | 100リクエスト/分 |
| 認証API | 20リクエスト/分 |
| ログイン | 5リクエスト/分 |

## キャッシュ

| エンドポイント | キャッシュ時間 |
|---------------|--------------|
| /muscle-groups | 1時間 |
| /exercises | 1時間 |
| /exercise-icons | 24時間 |