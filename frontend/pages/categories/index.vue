<script setup lang="ts">
import type { Category, TxType } from '~/types'

const api = useExpenseApi()

const typeTab = ref<TxType>('expense')

const {
  data: categories,
  pending,
  error,
  refresh,
} = await useAsyncData(
  'categories-manage',
  () => api.listCategories(typeTab.value),
  { watch: [typeTab] },
)

// 追加フォーム
const newName = ref('')
const adding = ref(false)
const formError = ref('')

async function onAdd() {
  formError.value = ''
  if (!newName.value.trim()) {
    formError.value = '名称を入力してください'
    return
  }
  adding.value = true
  try {
    await api.createCategory({ name: newName.value.trim(), type: typeTab.value })
    newName.value = ''
    await refresh()
  } catch (e) {
    formError.value = extractMessage(e)
  } finally {
    adding.value = false
  }
}

// 行内編集
const editingId = ref<number | null>(null)
const editName = ref('')
const rowError = ref('')

function startEdit(c: Category) {
  editingId.value = c.id
  editName.value = c.name
  rowError.value = ''
}

function cancelEdit() {
  editingId.value = null
  editName.value = ''
  rowError.value = ''
}

async function onSave(c: Category) {
  rowError.value = ''
  if (!editName.value.trim()) {
    rowError.value = '名称を入力してください'
    return
  }
  try {
    await api.updateCategory(c.id, { name: editName.value.trim(), type: c.type })
    cancelEdit()
    await refresh()
  } catch (e) {
    rowError.value = extractMessage(e)
  }
}

async function onDelete(c: Category) {
  if (!confirm(`カテゴリ「${c.name}」を削除しますか？`)) return
  try {
    await api.deleteCategory(c.id)
    await refresh()
  } catch (e) {
    alert(extractMessage(e))
  }
}
</script>

<template>
  <div>
    <h1 class="text-xl font-bold mb-4">カテゴリ管理</h1>

    <!-- 種別タブ -->
    <div class="flex gap-1 mb-4 border-b">
      <button
        v-for="tab in (['expense', 'income'] as TxType[])"
        :key="tab"
        class="px-4 py-2 text-sm border-b-2 -mb-px"
        :class="
          typeTab === tab
            ? 'border-indigo-600 text-indigo-600 font-medium'
            : 'border-transparent text-gray-500 hover:text-indigo-600'
        "
        @click="typeTab = tab"
      >
        {{ tab === 'expense' ? '支出' : '収入' }}
      </button>
    </div>

    <!-- 追加フォーム -->
    <form class="flex flex-wrap items-start gap-2 mb-4" @submit.prevent="onAdd">
      <div>
        <input
          v-model="newName"
          type="text"
          maxlength="50"
          placeholder="新しいカテゴリ名"
          class="rounded border-gray-300 shadow-sm text-sm"
        />
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
      カテゴリの取得に失敗しました。バックエンド（:8080）が起動しているか確認してください。
    </p>

    <div class="bg-white rounded-lg border overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-gray-500 text-left">
          <tr>
            <th class="px-4 py-2 font-medium">名称</th>
            <th class="px-4 py-2 font-medium text-right">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          <tr v-if="pending">
            <td colspan="2" class="px-4 py-6 text-center text-gray-400">読み込み中...</td>
          </tr>
          <tr v-else-if="!categories || categories.length === 0">
            <td colspan="2" class="px-4 py-6 text-center text-gray-400">
              カテゴリがありません。上のフォームから追加できます。
            </td>
          </tr>
          <tr v-for="c in categories" v-else :key="c.id" class="hover:bg-gray-50">
            <td class="px-4 py-2">
              <template v-if="editingId === c.id">
                <input
                  v-model="editName"
                  type="text"
                  maxlength="50"
                  class="rounded border-gray-300 shadow-sm text-sm w-full max-w-xs"
                  @keyup.enter="onSave(c)"
                />
                <p v-if="rowError" class="text-red-600 text-xs mt-1">{{ rowError }}</p>
              </template>
              <template v-else>{{ c.name }}</template>
            </td>
            <td class="px-4 py-2 text-right whitespace-nowrap">
              <template v-if="editingId === c.id">
                <button class="text-indigo-600 hover:underline mr-3" @click="onSave(c)">保存</button>
                <button class="text-gray-500 hover:underline" @click="cancelEdit">キャンセル</button>
              </template>
              <template v-else>
                <button class="text-indigo-600 hover:underline mr-3" @click="startEdit(c)">編集</button>
                <button class="text-red-600 hover:underline" @click="onDelete(c)">削除</button>
              </template>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
