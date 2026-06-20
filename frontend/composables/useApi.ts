// バックエンド API への fetch ラッパ。
// 各画面・composable はこれを経由して API を呼び出す。
export function useApi() {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase

  return {
    get: <T>(path: string) => $fetch<T>(path, { baseURL }),
    post: <T>(path: string, body: unknown) =>
      $fetch<T>(path, { baseURL, method: 'POST', body }),
    put: <T>(path: string, body: unknown) =>
      $fetch<T>(path, { baseURL, method: 'PUT', body }),
    delete: <T>(path: string) =>
      $fetch<T>(path, { baseURL, method: 'DELETE' }),
  }
}
