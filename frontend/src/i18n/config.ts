import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import zhTranslations from './locales/zh.json';
import enTranslations from './locales/en.json';

// Detect user's preferred language based on IP geolocation
const detectLanguage = async (): Promise<string> => {
  try {
    const apiKey = process.env.REACT_APP_IPGEO_API_KEY;
    const baseUrl = process.env.REACT_APP_IPGEO_BASE_URL;
    
    if (!apiKey || !baseUrl) {
      return 'zh'; // Default to Chinese
    }

    const response = await fetch(`${baseUrl}/ipgeo?apiKey=${apiKey}`);
    const data = await response.json();
    
    // Check if country is in Chinese-speaking regions
    const chineseRegions = ['CN', 'HK', 'MO', 'TW', 'SG'];
    if (chineseRegions.includes(data.country_code2)) {
      return 'zh';
    }
    
    return 'en';
  } catch (error) {
    console.error('Failed to detect language:', error);
    return 'zh'; // Default to Chinese on error
  }
};

// Initialize i18n
const initI18n = async () => {
  const savedLanguage = localStorage.getItem('language');
  const detectedLanguage = savedLanguage || await detectLanguage();

  i18n
    .use(initReactI18next)
    .init({
      resources: {
        zh: {
          translation: zhTranslations,
        },
        en: {
          translation: enTranslations,
        },
      },
      lng: detectedLanguage,
      fallbackLng: 'zh',
      interpolation: {
        escapeValue: false,
      },
    });

  // Save language preference
  i18n.on('languageChanged', (lng) => {
    localStorage.setItem('language', lng);
  });
};

initI18n();

export default i18n;
