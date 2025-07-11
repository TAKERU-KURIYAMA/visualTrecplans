export const Footer = () => {
  return (
    <footer className="bg-gray-50 border-t">
      <div className="container py-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              VisualTrecplans
            </h3>
            <p className="text-sm text-gray-600">
              ビジュアル重視のトレーニング・サプリ管理アプリケーション
            </p>
          </div>
          
          <div>
            <h4 className="text-sm font-semibold text-gray-900 mb-4">
              機能
            </h4>
            <ul className="space-y-2 text-sm text-gray-600">
              <li>トレーニング記録</li>
              <li>サプリメント管理</li>
              <li>進捗可視化</li>
              <li>人体図UI</li>
            </ul>
          </div>
          
          <div>
            <h4 className="text-sm font-semibold text-gray-900 mb-4">
              サポート
            </h4>
            <ul className="space-y-2 text-sm text-gray-600">
              <li>ヘルプセンター</li>
              <li>お問い合わせ</li>
              <li>利用規約</li>
              <li>プライバシーポリシー</li>
            </ul>
          </div>
          
          <div>
            <h4 className="text-sm font-semibold text-gray-900 mb-4">
              言語
            </h4>
            <ul className="space-y-2 text-sm text-gray-600">
              <li>日本語</li>
              <li>English</li>
              <li>Español</li>
              <li>Français</li>
            </ul>
          </div>
        </div>
        
        <div className="mt-8 pt-8 border-t border-gray-200">
          <p className="text-center text-sm text-gray-500">
            © 2025 VisualTrecplans. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  )
}