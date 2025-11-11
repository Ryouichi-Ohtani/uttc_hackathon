import { useTranslation } from '@/i18n/useTranslation'
import { LanguageIcon } from '@heroicons/react/24/outline'

export const LanguageSwitcher = () => {
  const { language, setLanguage } = useTranslation()

  const toggleLanguage = () => {
    const nextLang = language === 'ja' ? 'en' : language === 'en' ? 'zh' : 'ja'
    setLanguage(nextLang)
  }

  const languageLabels = {
    ja: '日本語',
    en: 'English',
    zh: '中文'
  }

  return (
    <button
      onClick={toggleLanguage}
      className="flex items-center gap-2 px-3 py-2 rounded-lg text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
      aria-label="Switch language"
      title={`Current: ${languageLabels[language]}`}
    >
      <LanguageIcon className="w-5 h-5" />
      <span className="text-sm font-medium">{language.toUpperCase()}</span>
    </button>
  )
}
