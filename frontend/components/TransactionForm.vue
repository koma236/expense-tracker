<script setup lang="ts">
import type { Category, Transaction, TransactionInput, TxType } from '~/types'

const props = defineProps<{
  initial?: Transaction | null
  submitting?: boolean
  fieldErrors?: Record<string, string>
  submitLabel?: string
}>()

const emit = defineEmits<{
  (e: 'submit', value: TransactionInput): void
}>()

const api = useExpenseApi()

// フォーム状態（初期値は編集対象 or 既定）
const todayStr = new Date().toISOString().slice(0, 10)
const type = ref<TxType>(props.initial?.type ?? 'expense')
const occurredOn = ref<string>(props.initial?.occurred_on ?? todayStr)
const amount = ref<number | null>(props.initial?.amount ?? null)
const categoryId = ref<number | null>(props.initial?.category.id ?? null)
const memo = ref<string>(props.initial?.memo ?? '')

// 全カテゴリを取得し、選択中の種別で絞り込む
const { data: categories } = await useAsyncData('form-categories', () => api.listCategories())
const filteredCategories = computed<Category[]>(
  () => (categories.value ?? []).filter((c) => c.type === type.value),
)

// 種別変更時、現在のカテゴリが新しい種別に属さなければリセット
watch(type, () => {
  if (!filteredCategories.value.some((c) => c.id === categoryId.value)) {
    categoryId.value = null
  }
})

function onSubmit() {
  emit('submit', {
    occurred_on: occurredOn.value,
    amount: Number(amount.value ?? 0),
    type: type.value,
    category_id: Number(categoryId.value ?? 0),
    memo: memo.value.trim() === '' ? null : memo.value,
  })
}

function errorFor(field: string) {
  return props.fieldErrors?.[field]
}
</script>

<template>
  <form class="space-y-5 bg-white rounded-lg border p-6" @submit.prevent="onSubmit">
    <!-- 種別 -->
    <div>
      <span class="block text-sm font-medium mb-1">種別</span>
      <div class="flex gap-4">
        <label class="inline-flex items-center gap-2">
          <input v-model="type" type="radio" value="expense" class="text-indigo-600" />
          <span>支出</span>
        </label>
        <label class="inline-flex items-center gap-2">
          <input v-model="type" type="radio" value="income" class="text-indigo-600" />
          <span>収入</span>
        </label>
      </div>
    </div>

    <!-- 日付 -->
    <div>
      <label class="block text-sm font-medium mb-1" for="occurred_on">日付</label>
      <input
        id="occurred_on"
        v-model="occurredOn"
        type="date"
        class="w-full rounded border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
      />
      <p v-if="errorFor('occurred_on')" class="mt-1 text-sm text-red-600">{{ errorFor('occurred_on') }}</p>
    </div>

    <!-- 金額 -->
    <div>
      <label class="block text-sm font-medium mb-1" for="amount">金額（円）</label>
      <input
        id="amount"
        v-model.number="amount"
        type="number"
        min="1"
        placeholder="0"
        class="w-full rounded border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
      />
      <p v-if="errorFor('amount')" class="mt-1 text-sm text-red-600">{{ errorFor('amount') }}</p>
    </div>

    <!-- カテゴリ -->
    <div>
      <label class="block text-sm font-medium mb-1" for="category">カテゴリ</label>
      <select
        id="category"
        v-model.number="categoryId"
        class="w-full rounded border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
      >
        <option :value="null" disabled>選択してください</option>
        <option v-for="c in filteredCategories" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
      <p v-if="errorFor('category_id')" class="mt-1 text-sm text-red-600">{{ errorFor('category_id') }}</p>
    </div>

    <!-- メモ -->
    <div>
      <label class="block text-sm font-medium mb-1" for="memo">メモ（任意）</label>
      <input
        id="memo"
        v-model="memo"
        type="text"
        maxlength="255"
        class="w-full rounded border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
      />
      <p v-if="errorFor('memo')" class="mt-1 text-sm text-red-600">{{ errorFor('memo') }}</p>
    </div>

    <div class="flex items-center gap-3 pt-2">
      <button
        type="submit"
        :disabled="submitting"
        class="rounded bg-indigo-600 px-4 py-2 text-white font-medium hover:bg-indigo-700 disabled:opacity-50"
      >
        {{ submitLabel ?? '保存' }}
      </button>
      <NuxtLink to="/transactions" class="text-gray-500 hover:text-gray-700">キャンセル</NuxtLink>
    </div>
  </form>
</template>
