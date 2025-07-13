# シーケンス図

## 概要
VisualTrecplansアプリケーションの主要な機能フローをシーケンス図で表現します。

## 1. ユーザー登録フロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Frontend as フロントエンド
    participant Backend as バックエンドAPI
    participant DB as データベース
    participant Email as メールサービス

    User->>Frontend: 登録フォーム入力
    Frontend->>Frontend: バリデーション（Zod）
    
    alt バリデーションエラー
        Frontend->>User: エラーメッセージ表示
    else バリデーション成功
        Frontend->>Backend: POST /api/v1/auth/register
        Backend->>Backend: パスワード強度チェック
        Backend->>Backend: ユーザー重複チェック
        Backend->>Backend: パスワードハッシュ化（bcrypt）
        Backend->>DB: ユーザー情報保存
        DB-->>Backend: 保存完了
        Backend->>Backend: 監査ログ記録
        Backend-->>Frontend: 201 Created + ユーザー情報
        Frontend->>Frontend: 認証ストア更新
        Frontend->>User: 登録完了メッセージ
        Frontend->>Frontend: ダッシュボードにリダイレクト
    end
```

## 2. ログインフロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Frontend as フロントエンド
    participant Backend as バックエンドAPI
    participant DB as データベース
    participant Cache as キャッシュ

    User->>Frontend: ログインフォーム入力
    Frontend->>Frontend: バリデーション（Zod）
    Frontend->>Backend: POST /api/v1/auth/login
    
    Backend->>Cache: ブルートフォース攻撃チェック
    alt レート制限超過
        Backend-->>Frontend: 429 Too Many Requests
        Frontend->>User: エラーメッセージ表示
    else 制限内
        Backend->>DB: ユーザー検索（email）
        DB-->>Backend: ユーザー情報
        Backend->>Backend: パスワード検証（bcrypt）
        
        alt 認証失敗
            Backend->>Cache: 失敗回数カウント
            Backend->>Backend: 監査ログ記録
            Backend-->>Frontend: 401 Unauthorized
            Frontend->>User: エラーメッセージ表示
        else 認証成功
            Backend->>Backend: JWT生成（アクセス+リフレッシュ）
            Backend->>Cache: リフレッシュトークン保存
            Backend->>Backend: 監査ログ記録
            Backend-->>Frontend: 200 OK + トークン + ユーザー情報
            Frontend->>Frontend: 認証ストア更新
            Frontend->>Frontend: LocalStorage保存
            Frontend->>User: ダッシュボードにリダイレクト
        end
    end
```

## 3. ワークアウト記録フロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Frontend as フロントエンド
    participant Backend as バックエンドAPI
    participant DB as データベース

    User->>Frontend: ワークアウト記録ページアクセス
    
    par マスタデータ取得
        Frontend->>Backend: GET /api/v1/muscle-groups
        Backend->>DB: 筋肉部位取得
        DB-->>Backend: 筋肉部位一覧
        Backend-->>Frontend: 200 OK + 筋肉部位データ
    and
        Frontend->>Backend: GET /api/v1/exercises
        Backend->>DB: エクササイズ取得
        DB-->>Backend: エクササイズ一覧
        Backend-->>Frontend: 200 OK + エクササイズデータ
    end
    
    Frontend->>Frontend: ドロップダウン初期化
    Frontend->>User: フォーム表示
    
    User->>Frontend: 筋肉部位選択
    Frontend->>Frontend: エクササイズフィルタリング
    
    User->>Frontend: エクササイズ選択
    User->>Frontend: 重量・回数・セット数入力
    User->>Frontend: メモ入力（任意）
    User->>Frontend: 実施日時選択
    User->>Frontend: 保存ボタンクリック
    
    Frontend->>Frontend: バリデーション（Zod）
    
    alt バリデーションエラー
        Frontend->>User: エラーメッセージ表示
    else バリデーション成功
        Frontend->>Backend: POST /api/v1/workouts + Authorization
        Backend->>Backend: JWT検証
        Backend->>Backend: 筋肉部位バリデーション
        Backend->>DB: ワークアウト保存
        DB-->>Backend: 保存完了
        Backend-->>Frontend: 201 Created + ワークアウトデータ
        Frontend->>Frontend: ワークアウトストア更新
        Frontend->>User: 成功メッセージ表示
        Frontend->>Frontend: 履歴ページにリダイレクト
    end
```

## 4. ワークアウト履歴表示フロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Frontend as フロントエンド
    participant Backend as バックエンドAPI
    participant DB as データベース

    User->>Frontend: 履歴ページアクセス
    Frontend->>Frontend: 初期フィルター設定（default）
    
    Frontend->>Backend: GET /api/v1/workouts?page=1&per_page=20 + Authorization
    Backend->>Backend: JWT検証
    Backend->>DB: ワークアウト検索（ユーザー別）
    DB-->>Backend: ワークアウト一覧 + 総件数
    Backend-->>Frontend: 200 OK + ページネーション情報
    
    Frontend->>Frontend: 日付別グループ化
    Frontend->>Frontend: コンポーネント描画
    Frontend->>User: 履歴一覧表示
    
    opt フィルター変更
        User->>Frontend: フィルター設定変更
        Frontend->>Backend: GET /api/v1/workouts?muscle_group=chest&start_date=2024-01-01
        Backend->>DB: フィルター付き検索
        DB-->>Backend: フィルター結果
        Backend-->>Frontend: 200 OK + フィルター結果
        Frontend->>Frontend: 一覧更新
        Frontend->>User: フィルター結果表示
    end
    
    opt ページング
        User->>Frontend: 次ページクリック
        Frontend->>Backend: GET /api/v1/workouts?page=2&per_page=20
        Backend->>DB: 次ページデータ取得
        DB-->>Backend: 次ページデータ
        Backend-->>Frontend: 200 OK + 次ページデータ
        Frontend->>Frontend: 一覧更新
        Frontend->>User: 次ページ表示
    end
```

## 5. ワークアウト編集フロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Frontend as フロントエンド
    participant Backend as バックエンドAPI
    participant DB as データベース

    User->>Frontend: 編集ボタンクリック
    Frontend->>Backend: GET /api/v1/workouts/:id + Authorization
    Backend->>Backend: JWT検証
    Backend->>DB: ワークアウト取得
    DB-->>Backend: ワークアウトデータ
    Backend->>Backend: 所有者権限チェック
    
    alt 権限なし
        Backend-->>Frontend: 403 Forbidden
        Frontend->>User: エラーメッセージ表示
    else 権限あり
        Backend-->>Frontend: 200 OK + ワークアウトデータ
        Frontend->>Frontend: フォーム初期化
        Frontend->>User: 編集フォーム表示
        
        User->>Frontend: データ修正
        User->>Frontend: 保存ボタンクリック
        
        Frontend->>Frontend: バリデーション（Zod）
        Frontend->>Backend: PUT /api/v1/workouts/:id + Authorization
        Backend->>Backend: JWT検証
        Backend->>Backend: 所有者権限チェック
        Backend->>Backend: バリデーション
        Backend->>DB: ワークアウト更新
        DB-->>Backend: 更新完了
        Backend-->>Frontend: 200 OK + 更新データ
        Frontend->>Frontend: ワークアウトストア更新
        Frontend->>User: 成功メッセージ表示
        Frontend->>Frontend: 履歴ページにリダイレクト
    end
```

## 6. ワークアウト削除フロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Frontend as フロントエンド
    participant Backend as バックエンドAPI
    participant DB as データベース

    User->>Frontend: 削除ボタンクリック
    Frontend->>Frontend: 確認ダイアログ表示
    
    alt キャンセル
        User->>Frontend: キャンセルクリック
        Frontend->>Frontend: ダイアログクローズ
    else 削除実行
        User->>Frontend: 削除確認
        Frontend->>Backend: DELETE /api/v1/workouts/:id + Authorization
        Backend->>Backend: JWT検証
        Backend->>DB: ワークアウト取得
        DB-->>Backend: ワークアウトデータ
        Backend->>Backend: 所有者権限チェック
        
        alt 権限なし
            Backend-->>Frontend: 403 Forbidden
            Frontend->>User: エラーメッセージ表示
        else 権限あり
            Backend->>DB: ソフトデリート実行
            DB-->>Backend: 削除完了
            Backend-->>Frontend: 204 No Content
            Frontend->>Frontend: ワークアウトストア更新
            Frontend->>User: 成功メッセージ表示
            Frontend->>Frontend: 一覧から該当項目削除
        end
    end
```

## 7. 統計データ取得フロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Frontend as フロントエンド
    participant Backend as バックエンドAPI
    participant DB as データベース

    User->>Frontend: ダッシュボードアクセス
    Frontend->>Backend: GET /api/v1/workouts/stats?period=month + Authorization
    Backend->>Backend: JWT検証
    Backend->>DB: 統計クエリ実行
    
    par 基本統計
        DB->>DB: 総ワークアウト数集計
        DB->>DB: 総セット数・回数集計
        DB->>DB: 総重量集計
    and 筋肉部位別統計
        DB->>DB: 筋肉部位別ワークアウト数集計
    and よく使用されるエクササイズ
        DB->>DB: エクササイズ使用頻度集計
    and 週次進捗
        DB->>DB: 週別ワークアウト統計集計
    end
    
    DB-->>Backend: 統計データ
    Backend->>Backend: データ整形・計算
    Backend-->>Frontend: 200 OK + 統計データ
    Frontend->>Frontend: グラフ・チャート描画
    Frontend->>User: ダッシュボード表示
```

## 8. トークンリフレッシュフロー

```mermaid
sequenceDiagram
    participant Frontend as フロントエンド
    participant Backend as バックエンドAPI
    participant Cache as キャッシュ

    Frontend->>Backend: API request + Authorization
    Backend->>Backend: JWT検証
    
    alt トークン有効
        Backend-->>Frontend: 200 OK + レスポンス
    else トークン期限切れ
        Backend-->>Frontend: 401 Unauthorized + token_expired
        Frontend->>Frontend: リフレッシュトークン取得
        
        alt リフレッシュトークンなし
            Frontend->>Frontend: ログイン画面にリダイレクト
        else リフレッシュトークンあり
            Frontend->>Backend: POST /api/v1/auth/refresh
            Backend->>Cache: リフレッシュトークン検証
            
            alt 無効/期限切れ
                Backend-->>Frontend: 401 Unauthorized
                Frontend->>Frontend: ログイン画面にリダイレクト
            else 有効
                Backend->>Backend: 新しいJWT生成
                Backend->>Cache: 新しいリフレッシュトークン保存
                Backend-->>Frontend: 200 OK + 新しいトークン
                Frontend->>Frontend: 認証ストア更新
                Frontend->>Backend: 元のAPI request再実行
                Backend-->>Frontend: 200 OK + レスポンス
            end
        end
    end
```

## 9. エラーハンドリングフロー

```mermaid
sequenceDiagram
    participant Frontend as フロントエンド
    participant Backend as バックエンドAPI
    participant ErrorBoundary as エラーバウンダリー
    participant Toast as トースト通知

    Frontend->>Backend: API request
    
    alt サーバーエラー (5xx)
        Backend-->>Frontend: 500 Internal Server Error
        Frontend->>ErrorBoundary: エラー捕捉
        ErrorBoundary->>Toast: エラートースト表示
        ErrorBoundary->>Frontend: フォールバックUI表示
    else クライアントエラー (4xx)
        Backend-->>Frontend: 400/401/403/404
        Frontend->>Frontend: エラー種別判定
        
        alt 認証エラー (401)
            Frontend->>Frontend: トークンリフレッシュ試行
        else バリデーションエラー (400/422)
            Frontend->>Toast: バリデーションエラー表示
        else 権限エラー (403)
            Frontend->>Toast: 権限エラー表示
        else リソース未存在 (404)
            Frontend->>Toast: 404エラー表示
        end
    else ネットワークエラー
        Frontend->>Frontend: ネットワークエラー判定
        Frontend->>Toast: 接続エラー表示
        Frontend->>Frontend: リトライ機能提供
    end
```