// 画面間で共有する「対象月」の状態。ヘッダーの月セレクタと取引一覧が参照する。
export function useMonth() {
  const current = useState<string>('year_month', () => {
    const d = new Date()
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
  })

  function shift(delta: number) {
    const [y, m] = current.value.split('-').map(Number)
    const d = new Date(y, m - 1 + delta, 1)
    current.value = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
  }

  // 表示用ラベル（例: 2026年6月）
  const label = computed(() => {
    const [y, m] = current.value.split('-')
    return `${y}年${Number(m)}月`
  })

  return { current, label, prev: () => shift(-1), next: () => shift(1) }
}
