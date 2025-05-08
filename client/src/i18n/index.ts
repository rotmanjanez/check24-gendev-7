import { createI18n } from 'vue-i18n'
import de from './locales/de'
import en from './locales/en'

const i18n = createI18n({
    locale: 'de',
    fallbackLocale: 'de',
    missingWarn: true,
    fallbackWarn: true,
    globalInjection: true,
    messages: { de, en, }
  })

export default i18n