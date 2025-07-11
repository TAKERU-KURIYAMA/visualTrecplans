import { useState, useEffect } from 'react'

type Theme = 'light' | 'dark' | 'system'

export const useDarkMode = () => {
  const [theme, setTheme] = useState<Theme>(() => {
    // ローカルストレージからテーマを取得、なければsystemを使用
    const stored = localStorage.getItem('theme') as Theme
    return stored || 'system'
  })

  const [isDark, setIsDark] = useState(false)

  useEffect(() => {
    const applyTheme = (selectedTheme: Theme) => {
      const root = window.document.documentElement
      
      if (selectedTheme === 'system') {
        // システムの設定を使用
        const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        setIsDark(systemPrefersDark)
        
        if (systemPrefersDark) {
          root.classList.add('dark')
        } else {
          root.classList.remove('dark')
        }
      } else {
        // 明示的なテーマ設定
        const isDarkTheme = selectedTheme === 'dark'
        setIsDark(isDarkTheme)
        
        if (isDarkTheme) {
          root.classList.add('dark')
        } else {
          root.classList.remove('dark')
        }
      }
    }

    applyTheme(theme)

    // システムテーマの変更を監視
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    const handleChange = () => {
      if (theme === 'system') {
        applyTheme('system')
      }
    }

    mediaQuery.addEventListener('change', handleChange)
    return () => mediaQuery.removeEventListener('change', handleChange)
  }, [theme])

  const setThemeAndStore = (newTheme: Theme) => {
    setTheme(newTheme)
    localStorage.setItem('theme', newTheme)
  }

  const toggleTheme = () => {
    if (theme === 'light') {
      setThemeAndStore('dark')
    } else if (theme === 'dark') {
      setThemeAndStore('system')
    } else {
      setThemeAndStore('light')
    }
  }

  return {
    theme,
    isDark,
    setTheme: setThemeAndStore,
    toggleTheme,
  }
}