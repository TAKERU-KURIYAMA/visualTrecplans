/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string
  readonly VITE_APP_NAME: string
  readonly VITE_APP_VERSION: string
  readonly VITE_APP_TITLE: string
  readonly VITE_APP_DESCRIPTION: string
  
  // 機能フラグ
  readonly VITE_FEATURE_SOCIAL: string
  readonly VITE_FEATURE_EXPORT: string
  readonly VITE_FEATURE_ANALYTICS: string
  
  // 外部API設定
  readonly VITE_GOOGLE_ANALYTICS_ID?: string
  readonly VITE_SENTRY_DSN?: string
  
  // デバッグ設定
  readonly VITE_DEBUG: string
  readonly VITE_LOG_LEVEL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}