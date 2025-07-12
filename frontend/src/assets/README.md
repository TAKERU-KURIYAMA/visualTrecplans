# assets ディレクトリ

## 概要
静的ファイル（画像、フォント、アイコン等）を配置するディレクトリです。
アプリケーションで使用するメディアファイルを管理します。

## ディレクトリ構造例
```
assets/
├── images/         # 画像ファイル
│   ├── logo.svg
│   ├── body/      # 人体図関連の画像
│   └── icons/     # アイコン画像
├── fonts/          # カスタムフォント
└── data/           # 静的JSONデータ
    └── exercises.json  # エクササイズマスタデータ
```

## 命名規則
- ファイル名は kebab-case を使用
- 画像は適切な形式を選択（SVG推奨）
- アイコンは統一されたサイズと命名規則を維持

## 画像最適化
- SVGファイルを優先的に使用
- ラスター画像は適切な圧縮を実施
- レスポンシブ対応が必要な場合は複数サイズを準備

## 例
```
assets/
├── images/
│   ├── logo.svg
│   ├── body/
│   │   ├── body-front.svg
│   │   ├── body-back.svg
│   │   └── muscle-groups/
│   │       ├── chest.svg
│   │       ├── back.svg
│   │       └── legs.svg
│   └── icons/
│       ├── dumbbell.svg
│       ├── calendar.svg
│       └── chart.svg
```