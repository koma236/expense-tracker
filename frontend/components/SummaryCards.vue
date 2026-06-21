<script setup lang="ts">
import type { SummaryTotals } from '~/types'

const props = defineProps<{ totals: SummaryTotals }>()

function formatYen(n: number) {
  return `¥${n.toLocaleString('ja-JP')}`
}

// 収支差は符号付きで表示する。
const balanceLabel = computed(() => {
  const b = props.totals.balance
  const sign = b > 0 ? '+' : b < 0 ? '-' : ''
  return `${sign}¥${Math.abs(b).toLocaleString('ja-JP')}`
})
</script>

<template>
  <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
    <div class="bg-white rounded-lg border p-4">
      <p class="text-xs text-gray-500">収入</p>
      <p class="mt-1 text-lg font-bold tabular-nums text-emerald-600">
        {{ formatYen(totals.income) }}
      </p>
    </div>
    <div class="bg-white rounded-lg border p-4">
      <p class="text-xs text-gray-500">支出</p>
      <p class="mt-1 text-lg font-bold tabular-nums text-red-600">
        {{ formatYen(totals.expense) }}
      </p>
    </div>
    <div class="bg-white rounded-lg border p-4">
      <p class="text-xs text-gray-500">収支</p>
      <p
        class="mt-1 text-lg font-bold tabular-nums"
        :class="totals.balance >= 0 ? 'text-emerald-600' : 'text-red-600'"
      >
        {{ balanceLabel }}
      </p>
    </div>
  </div>
</template>
