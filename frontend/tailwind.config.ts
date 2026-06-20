import type { Config } from 'tailwindcss'

// @nuxtjs/tailwindcss が content を自動設定するが、明示しておく。
export default <Partial<Config>>{
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './app.vue',
  ],
}
