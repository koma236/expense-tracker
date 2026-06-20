<script setup lang="ts">
import type { TxType } from '~/types'

const { current } = useMonth()
const api = useExpenseApi()

const typeFilter = ref<TxType | ''>('')
const categoryFilter = ref<number | null>(null)

const { data: categories } = await useAsyncData('list-categories', () => api.listCategories())

const {
  data: transactions,
  pending,
  error,
  refresh,
} = await useAsyncData(
  'transactions',
  () =>
    api.listTransactions({
      yearMonth: current.value,
      type: typeFilter.value || undefined,
      categoryId: categoryFilter.value || undefined,
    }),
  { watch: [current, typeFilter, categoryFilter] },
)

const filteredCategories = computed(() =>
  (categories.value ?? []).filter((c) => !typeFilter.value || c.type === typeFilter.value),
)

function formatYen(n: number) {
  return `¥${n.toLocaleString('ja-JP')}`
}

async function onDelete(id: number) {
  if (!confirm('この取引を削除しますか？')) return
  try {
    await api.deleteTransaction(id)
    await refresh()
  } catch (e) {
    alert(extractMessage(e))
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-xl font-bold">取引一覧</h1>
      <NuxtLink
        to="/transactions/new"
        class="rounded bg-indigo-600 px-3 py-2 text-sm text-white font-medium hover:bg-indigo-700"
      >
        ＋ 新規
      </NuxtLink>
    </div>

    <!-- 絞り込み -->
    <div class="flex flex-wrap gap-3 mb-4 text-sm">
      <select v-model="typeFilter" class="rounded border-gray-300 shadow-sm">
        <option value="">すべての種別</option>
        <option value="expense">支出</option>
        <option value="income">収入</option>
      </select>
      <select v-model.number="categoryFilter" class="rounded border-gray-300 shadow-sm">
        <option :value="null">すべてのカテゴリ</option>
        <option v-for="c in filteredCategories" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
    </div>

    <p v-if="error" class="text-red-600 text-sm mb-4">
      取引の取得に失敗しました。バックエンド（:8080）が起動しているか確認してください。
    </p>

    <div class="bg-white rounded-lg border overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-gray-500 text-left">
          <tr>
            <th class="px-4 py-2 font-medium">日付</th>
            <th class="px-4 py-2 font-medium">カテゴリ</th>
            <th class="px-4 py-2 font-medium">メモ</th>
            <th class="px-4 py-2 font-medium text-right">金額</th>
            <th class="px-4 py-2 font-medium text-right">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          <tr v-if="pending">
            <td colspan="5" class="px-4 py-6 text-center text-gray-400">読み込み中...</td>
          </tr>
          <tr v-else-if="!transactions || transactions.length === 0">
            <td colspan="5" class="px-4 py-6 text-center text-gray-400">
              この月の取引はありません。「＋ 新規」から登録できます。
            </td>
          </tr>
          <tr v-for="t in transactions" v-else :key="t.id" class="hover:bg-gray-50">
            <td class="px-4 py-2 tabular-nums whitespace-nowrap">{{ t.occurred_on }}</td>
            <td class="px-4 py-2">{{ t.category.name }}</td>
            <td class="px-4 py-2 text-gray-500">{{ t.memo }}</td>
            <td
              class="px-4 py-2 text-right tabular-nums font-medium whitespace-nowrap"
              :class="t.type === 'income' ? 'text-emerald-600' : 'text-red-600'"
            >
              {{ t.type === 'income' ? '+' : '-' }}{{ formatYen(t.amount) }}
            </td>
            <td class="px-4 py-2 text-right whitespace-nowrap">
              <NuxtLink
                :to="`/transactions/${t.id}`"
                class="text-indigo-600 hover:underline mr-3"
              >
                編集
              </NuxtLink>
              <button class="text-red-600 hover:underline" @click="onDelete(t.id)">削除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
