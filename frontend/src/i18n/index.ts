import { translations, Language } from './translations'

class I18n {
  private currentLanguage: Language = 'ja'

  constructor() {
    // Load saved language from localStorage
    const savedLang = localStorage.getItem('language')
    if (savedLang && (savedLang === 'ja' || savedLang === 'en')) {
      this.currentLanguage = savedLang
    }
  }

  setLanguage(lang: Language) {
    this.currentLanguage = lang
    localStorage.setItem('language', lang)
    // Dispatch event for components to react to language change
    window.dispatchEvent(new CustomEvent('languagechange', { detail: lang }))
  }

  getLanguage(): Language {
    return this.currentLanguage
  }

  t(key: string): string {
    const keys = key.split('.')
    let value: any = translations[this.currentLanguage]

    for (const k of keys) {
      if (value && typeof value === 'object') {
        value = value[k]
      } else {
        return key // Return key if translation not found
      }
    }

    return typeof value === 'string' ? value : key
  }
}

export const i18n = new I18n()

// Hook for React components
import { useState, useEffect } from 'react'

export const useTranslation = () => {
  const [language, setLanguage] = useState(i18n.getLanguage())

  useEffect(() => {
    const handleLanguageChange = (e: Event) => {
      const customEvent = e as CustomEvent
      setLanguage(customEvent.detail)
    }

    window.addEventListener('languagechange', handleLanguageChange)
    return () => window.removeEventListener('languagechange', handleLanguageChange)
  }, [])

  return {
    t: (key: string) => i18n.t(key),
    language,
    setLanguage: (lang: Language) => i18n.setLanguage(lang),
  }
}
