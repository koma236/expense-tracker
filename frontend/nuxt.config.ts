// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },
  runtimeConfig: {
    public: {
      // バックエンド API のベース URL。NUXT_PUBLIC_API_BASE で上書き可能。
      apiBase: 'http://localhost:8080/api',
    },
  },
})
