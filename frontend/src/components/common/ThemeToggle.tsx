import {
  SunIcon,
  MoonIcon,
  ComputerDesktopIcon,
} from '@heroicons/react/24/outline'
import { useDarkMode } from '@/hooks/useDarkMode'

export const ThemeToggle = () => {
  const { theme, toggleTheme } = useDarkMode()

  const getIcon = () => {
    switch (theme) {
      case 'light':
        return <SunIcon className="h-5 w-5" />
      case 'dark':
        return <MoonIcon className="h-5 w-5" />
      case 'system':
        return <ComputerDesktopIcon className="h-5 w-5" />
      default:
        return <SunIcon className="h-5 w-5" />
    }
  }

  const getTitle = () => {
    switch (theme) {
      case 'light':
        return 'ライトモード（ダークモードに切り替え）'
      case 'dark':
        return 'ダークモード（システム設定に切り替え）'
      case 'system':
        return 'システム設定（ライトモードに切り替え）'
      default:
        return 'テーマを切り替え'
    }
  }

  return (
    <button
      onClick={toggleTheme}
      className="btn btn-ghost p-2 h-auto w-auto"
      title={getTitle()}
      aria-label="テーマを切り替え"
    >
      {getIcon()}
    </button>
  )
}
