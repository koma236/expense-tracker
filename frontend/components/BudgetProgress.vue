<script setup lang="ts">
import type { BudgetProgress } from '~/types'

defineProps<{ items: BudgetProgress[] }>()

function formatYen(n: number) {
  return `¥${n.toLocaleString('ja-JP')}`
}

function label(item: BudgetProgress) {
  return item.category ? item.category.name : '月全体'
}

// 消化率（%）。バーの幅は 100% で頭打ちにする。
function percent(item: BudgetProgress) {
  return Math.round(item.ratio * 100)
}
function barWidth(item: BudgetProgress) {
  return Math.min(percent(item), 100)
}

// 状態ごとの配色。
const barColor: Record<BudgetProgress['status'], string> = {
  under: 'bg-emerald-500',
  near: 'bg-amber-500',
  over: 'bg-red-500',
}
const statusLabel: Record<BudgetProgress['status'], string> = {
  under: '',
  near: '接近',
  over: '超過',
}
const statusBadge: Record<BudgetProgress['status'], string> = {
  under: '',
  near: 'bg-amber-100 text-amber-700',
  over: 'bg-red-100 text-red-700',
}
</script>

<template>
  <div class="bg-white rounded-lg border p-4">
    <h2 class="text-sm font-medium text-gray-700 mb-3">予算進捗</h2>

    <p v-if="items.length === 0" class="text-sm text-gray-400">
      予算が設定されていません。「予算」画面から設定できます。
    </p>

    <ul v-else class="space-y-3">
      <li v-for="item in items" :key="item.category ? item.category.id : 'all'">
        <div class="flex items-center justify-between text-sm mb-1">
          <span class="flex items-center gap-2">
            <span class="font-medium text-gray-700">{{ label(item) }}</span>
            <span
              v-if="item.status !== 'under'"
              class="text-xs px-1.5 py-0.5 rounded"
              :class="statusBadge[item.status]"
            >
              {{ statusLabel[item.status] }}
            </span>
          </span>
          <span class="tabular-nums text-gray-500">
            {{ formatYen(item.spent) }} / {{ formatYen(item.budget) }}
            <span class="ml-1 text-gray-400">({{ percent(item) }}%)</span>
          </span>
        </div>
        <div class="h-2 rounded bg-gray-100 overflow-hidden">
          <div
            class="h-full rounded transition-all"
            :class="barColor[item.status]"
            :style="{ width: `${barWidth(item)}%` }"
          />
        </div>
      </li>
    </ul>
  </div>
</template>
