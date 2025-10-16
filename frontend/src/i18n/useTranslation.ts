import { useState } from 'react';
import { translations, Language } from './translations';

const getStoredLanguage = (): Language => {
  const stored = localStorage.getItem('language');
  if (stored && (stored === 'ja' || stored === 'en' || stored === 'zh')) {
    return stored as Language;
  }

  // Detect browser language
  const browserLang = navigator.language.toLowerCase();
  if (browserLang.startsWith('ja')) return 'ja';
  if (browserLang.startsWith('zh')) return 'zh';
  return 'en';
};

export const useTranslation = () => {
  const [language, setLanguageState] = useState<Language>(getStoredLanguage());

  const setLanguage = (lang: Language) => {
    setLanguageState(lang);
    localStorage.setItem('language', lang);
  };

  const t = (key: string): string => {
    const keys = key.split('.');
    let value: any = translations[language];

    for (const k of keys) {
      value = value?.[k];
    }

    return value || key;
  };

  return { language, setLanguage, t };
};
