<script setup lang="ts">
import {
  BarElement,
  CategoryScale,
  Chart,
  Legend,
  LinearScale,
  Tooltip,
} from 'chart.js'
import { Bar } from 'vue-chartjs'
import type { TrendPoint } from '~/types'

Chart.register(BarElement, CategoryScale, LinearScale, Tooltip, Legend)

const props = defineProps<{ points: TrendPoint[] }>()

// ラベルは「6月」のように月だけを短く表示する。
function shortLabel(ym: string) {
  const m = Number(ym.split('-')[1])
  return `${m}月`
}

const chartData = computed(() => ({
  labels: props.points.map((p) => shortLabel(p.year_month)),
  datasets: [
    {
      label: '収入',
      data: props.points.map((p) => p.income),
      backgroundColor: '#10b981',
    },
    {
      label: '支出',
      data: props.points.map((p) => p.expense),
      backgroundColor: '#ef4444',
    },
  ],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'top' as const },
    tooltip: {
      callbacks: {
        label: (ctx: { dataset: { label?: string }; parsed: { y: number } }) =>
          `${ctx.dataset.label}: ¥${ctx.parsed.y.toLocaleString('ja-JP')}`,
      },
    },
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        callback: (value: number | string) => `¥${Number(value).toLocaleString('ja-JP')}`,
      },
    },
  },
}
</script>

<template>
  <div class="bg-white rounded-lg border p-4">
    <h2 class="text-sm font-medium text-gray-600 mb-2">月別収支推移</h2>
    <div class="h-64">
      <Bar :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>
