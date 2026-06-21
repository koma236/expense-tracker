<script setup lang="ts">
const { current, label } = useMonth()
const api = useExpenseApi()

const {
  data: summary,
  pending: summaryPending,
  error: summaryError,
} = await useAsyncData('summary', () => api.getSummary(current.value), { watch: [current] })

const { data: trend, error: trendError } = await useAsyncData(
  'summary-trend',
  () => api.getTrend(current.value, 6),
  { watch: [current] },
)

const hasError = computed(() => Boolean(summaryError.value || trendError.value))
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-xl font-bold">ダッシュボード</h1>
      <span class="text-sm text-gray-500">{{ label }}</span>
    </div>

    <p v-if="hasError" class="text-red-600 text-sm mb-4">
      集計の取得に失敗しました。バックエンド（:8080）が起動しているか確認してください。
    </p>

    <p v-else-if="summaryPending" class="text-gray-400 text-sm mb-4">読み込み中...</p>

    <template v-else>
      <!-- サマリ -->
      <SummaryCards v-if="summary" :totals="summary.totals" class="mb-4" />

      <!-- グラフ（Chart.js は client-only） -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <ClientOnly>
          <CategoryPieChart :items="summary?.expense_by_category ?? []" />
          <template #fallback>
            <div class="bg-white rounded-lg border p-4 h-72" />
          </template>
        </ClientOnly>
        <ClientOnly>
          <MonthlyTrendChart :points="trend ?? []" />
          <template #fallback>
            <div class="bg-white rounded-lg border p-4 h-72" />
          </template>
        </ClientOnly>
      </div>

      <!-- 予算進捗（F-4 で実装予定） -->
      <div class="mt-4 bg-white rounded-lg border p-4 text-sm text-gray-400">
        予算進捗は今後実装予定です（F-4）。
      </div>
    </template>
  </div>
</template>
