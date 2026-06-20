import type { ApiError } from '~/types'

// $fetch が投げる FetchError の .data から共通エラー形式を取り出す。

export function extractFieldErrors(err: unknown): Record<string, string> {
  const data = (err as { data?: ApiError })?.data
  const out: Record<string, string> = {}
  for (const d of data?.error?.details ?? []) {
    out[d.field] = d.message
  }
  return out
}

export function extractMessage(err: unknown): string {
  const data = (err as { data?: ApiError })?.data
  return data?.error?.message ?? '通信エラーが発生しました'
}
