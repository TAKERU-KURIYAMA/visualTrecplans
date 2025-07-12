# components ディレクトリ

## 概要
共通UIコンポーネントを配置するディレクトリです。
再利用可能なコンポーネントをAtomic Designの原則に基づいて整理します。

## ディレクトリ構造
```
components/
├── common/      # 汎用コンポーネント（Button、Input、Modal等）
├── layout/      # レイアウト関連（Header、Footer、Sidebar等）
├── forms/       # フォーム関連（FormField、Select、Checkbox等）
└── charts/      # グラフ・チャート関連（ProgressChart、BodyDiagram等）
```

## 命名規則
- コンポーネント名は PascalCase を使用
- ファイル名は コンポーネント名.tsx
- スタイルファイルは コンポーネント名.module.css（CSS Modules使用時）

## 例
```typescript
// components/common/Button.tsx
export const Button: React.FC<ButtonProps> = ({ children, onClick }) => {
  return <button onClick={onClick}>{children}</button>
}
```