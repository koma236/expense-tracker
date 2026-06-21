<script setup lang="ts">
import { ArcElement, Chart, Legend, Tooltip } from 'chart.js'
import { Doughnut } from 'vue-chartjs'
import type { ExpenseByCategory } from '~/types'

Chart.register(ArcElement, Tooltip, Legend)

const props = defineProps<{ items: ExpenseByCategory[] }>()

// カテゴリ用の配色（足りなければ循環）。
const palette = [
  '#6366f1', '#ef4444', '#10b981', '#f59e0b', '#3b82f6',
  '#ec4899', '#14b8a6', '#8b5cf6', '#f97316', '#84cc16',
]

const chartData = computed(() => ({
  labels: props.items.map((i) => i.category.name),
  datasets: [
    {
      data: props.items.map((i) => i.total),
      backgroundColor: props.items.map((_, idx) => palette[idx % palette.length]),
      borderWidth: 1,
      borderColor: '#fff',
    },
  ],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'right' as const },
    tooltip: {
      callbacks: {
        label: (ctx: { label: string; parsed: number }) =>
          `${ctx.label}: ¥${ctx.parsed.toLocaleString('ja-JP')}`,
      },
    },
  },
}
</script>

<template>
  <div class="bg-white rounded-lg border p-4">
    <h2 class="text-sm font-medium text-gray-600 mb-2">カテゴリ別支出</h2>
    <div v-if="items.length === 0" class="h-64 flex items-center justify-center text-gray-400 text-sm">
      この月の支出はありません。
    </div>
    <div v-else class="h-64">
      <Doughnut :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>
