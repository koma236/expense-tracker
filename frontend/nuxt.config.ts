// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss'],
  // 本番は静的SPAとして nginx から配信する（Node ランタイム不要）。
  // データ取得はすべてブラウザ側で実行され、相対パス /api を nginx が Go API へプロキシする。
  ssr: false,
  runtimeConfig: {
    public: {
      // バックエンド API のベース URL。NUXT_PUBLIC_API_BASE で上書き可能。
      apiBase: 'http://localhost:8080/api',
    },
  },
})
