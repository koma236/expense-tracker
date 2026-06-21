<script setup lang="ts">
import type { Budget } from '~/types'

const { current, label } = useMonth()
const api = useExpenseApi()

// 対象月の予算一覧。月セレクタの変更で再取得する。
const {
  data: budgets,
  pending,
  error,
  refresh,
} = await useAsyncData('budgets-manage', () => api.listBudgets(current.value), {
  watch: [current],
})

// 追加フォームの対象候補に使う支出カテゴリ。
const { data: categories } = await useAsyncData('budget-expense-categories', () =>
  api.listCategories('expense'),
)

function formatYen(n: number) {
  return `¥${n.toLocaleString('ja-JP')}`
}

// 既に予算が設定済みの対象（月全体 or カテゴリ）は候補から除外する。
const monthWideTaken = computed(() =>
  (budgets.value ?? []).some((b) => b.category === null),
)
const usedCategoryIds = computed(
  () => new Set((budgets.value ?? []).flatMap((b) => (b.category ? [b.category.id] : []))),
)
const availableCategories = computed(() =>
  (categories.value ?? []).filter((c) => !usedCategoryIds.value.has(c.id)),
)

// 追加フォーム。target は 'all'（月全体）またはカテゴリID（文字列）。
const target = ref<string>('all')
const newAmount = ref<number | null>(null)
const adding = ref(false)
const formError = ref('')

// 候補が変わったら選択をリセットする。
watchEffect(() => {
  if (target.value === 'all' && monthWideTaken.value) {
    target.value = availableCategories.value[0] ? String(availableCategories.value[0].id) : ''
  }
})

async function onAdd() {
  formError.value = ''
  if (!newAmount.value || newAmount.value < 1) {
    formError.value = '金額は1以上で入力してください'
    return
  }
  if (target.value === '') {
    formError.value = '対象を選択してください'
    return
  }
  adding.value = true
  try {
    await api.createBudget({
      category_id: target.value === 'all' ? null : Number(target.value),
      year_month: current.value,
      amount: newAmount.value,
    })
    newAmount.value = null
    target.value = monthWideTaken.value ? '' : 'all'
    await refresh()
  } catch (e) {
    formError.value = extractMessage(e)
  } finally {
    adding.value = false
  }
}

// 行内編集（金額のみ）
const editingId = ref<number | null>(null)
const editAmount = ref<number | null>(null)
const rowError = ref('')

function startEdit(b: Budget) {
  editingId.value = b.id
  editAmount.value = b.amount
  rowError.value = ''
}

function cancelEdit() {
  editingId.value = null
  editAmount.value = null
  rowError.value = ''
}

async function onSave(b: Budget) {
  rowError.value = ''
  if (!editAmount.value || editAmount.value < 1) {
    rowError.value = '金額は1以上で入力してください'
    return
  }
  try {
    await api.updateBudget(b.id, editAmount.value)
    cancelEdit()
    await refresh()
  } catch (e) {
    rowError.value = extractMessage(e)
  }
}

async function onDelete(b: Budget) {
  const name = b.category ? b.category.name : '月全体'
  if (!confirm(`「${name}」の予算を削除しますか？`)) return
  try {
    await api.deleteBudget(b.id)
    await refresh()
  } catch (e) {
    alert(extractMessage(e))
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-xl font-bold">予算設定</h1>
      <span class="text-sm text-gray-500">{{ label }}</span>
    </div>

    <!-- 追加フォーム -->
    <form class="flex flex-wrap items-start gap-2 mb-4" @submit.prevent="onAdd">
      <div>
        <select v-model="target" class="rounded border-gray-300 shadow-sm text-sm">
          <option v-if="!monthWideTaken" value="all">月全体</option>
          <option v-for="c in availableCategories" :key="c.id" :value="String(c.id)">
            {{ c.name }}
          </option>
        </select>
      </div>
      <div>
        <div class="flex items-center gap-1">
          <input
            v-model.number="newAmount"
            type="number"
            min="1"
            placeholder="金額"
            class="rounded border-gray-300 shadow-sm text-sm w-32"
          />
          <span class="text-sm text-gray-500">円</span>
        </div>
        <p v-if="formError" class="text-red-600 text-xs mt-1">{{ formError }}</p>
      </div>
      <button
        type="submit"
        :disabled="adding"
        class="rounded bg-indigo-600 px-3 py-2 text-sm text-white font-medium hover:bg-indigo-700 disabled:opacity-50"
      >
        ＋ 追加
      </button>
    </form>

    <p v-if="error" class="text-red-600 text-sm mb-4">
      予算の取得に失敗しました。バックエンド（:8080）が起動しているか確認してください。
    </p>

    <div class="bg-white rounded-lg border overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-gray-500 text-left">
          <tr>
            <th class="px-4 py-2 font-medium">対象</th>
            <th class="px-4 py-2 font-medium text-right">予算額</th>
            <th class="px-4 py-2 font-medium text-right">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          <tr v-if="pending">
            <td colspan="3" class="px-4 py-6 text-center text-gray-400">読み込み中...</td>
          </tr>
          <tr v-else-if="!budgets || budgets.length === 0">
            <td colspan="3" class="px-4 py-6 text-center text-gray-400">
              予算が設定されていません。上のフォームから追加できます。
            </td>
          </tr>
          <tr v-for="b in budgets" v-else :key="b.id" class="hover:bg-gray-50">
            <td class="px-4 py-2">
              <span v-if="!b.category" class="font-medium">月全体</span>
              <template v-else>{{ b.category.name }}</template>
            </td>
            <td class="px-4 py-2 text-right tabular-nums">
              <template v-if="editingId === b.id">
                <div class="flex items-center justify-end gap-1">
                  <input
                    v-model.number="editAmount"
                    type="number"
                    min="1"
                    class="rounded border-gray-300 shadow-sm text-sm w-28 text-right"
                    @keyup.enter="onSave(b)"
                  />
                  <span class="text-gray-500">円</span>
                </div>
                <p v-if="rowError" class="text-red-600 text-xs mt-1">{{ rowError }}</p>
              </template>
              <template v-else>{{ formatYen(b.amount) }}</template>
            </td>
            <td class="px-4 py-2 text-right whitespace-nowrap">
              <template v-if="editingId === b.id">
                <button class="text-indigo-600 hover:underline mr-3" @click="onSave(b)">保存</button>
                <button class="text-gray-500 hover:underline" @click="cancelEdit">キャンセル</button>
              </template>
              <template v-else>
                <button class="text-indigo-600 hover:underline mr-3" @click="startEdit(b)">編集</button>
                <button class="text-red-600 hover:underline" @click="onDelete(b)">削除</button>
              </template>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
