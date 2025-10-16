import en from './en.json'
import ja from './ja.json'

export type Language = 'en' | 'ja'

export interface I18nConfig {
  currentLanguage: Language
  translations: Record<string, any>
}

class I18n {
  private currentLanguage: Language = 'ja' // Default to Japanese
  private translations: Record<Language, any> = {
    en,
    ja
  }

  setLanguage(lang: Language) {
    this.currentLanguage = lang
    localStorage.setItem('language', lang)
  }

  getLanguage(): Language {
    const saved = localStorage.getItem('language') as Language
    return saved || this.currentLanguage
  }

  t(key: string): string {
    const keys = key.split('.')
    let value: any = this.translations[this.getLanguage()]

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

// Initialize language from localStorage
if (typeof window !== 'undefined') {
  const savedLang = localStorage.getItem('language') as Language
  if (savedLang) {
    i18n.setLanguage(savedLang)
  }
}
