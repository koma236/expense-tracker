<script setup lang="ts">
import type { TransactionInput } from '~/types'

const route = useRoute()
const id = Number(route.params.id)
const api = useExpenseApi()

const { data: transaction, error: loadError } = await useAsyncData(`transaction-${id}`, () =>
  api.getTransaction(id),
)

const submitting = ref(false)
const fieldErrors = ref<Record<string, string>>({})
const generalError = ref('')

async function onSubmit(input: TransactionInput) {
  submitting.value = true
  fieldErrors.value = {}
  generalError.value = ''
  try {
    await api.updateTransaction(id, input)
    await navigateTo('/transactions')
  } catch (e) {
    fieldErrors.value = extractFieldErrors(e)
    if (Object.keys(fieldErrors.value).length === 0) generalError.value = extractMessage(e)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="max-w-lg">
    <h1 class="text-xl font-bold mb-4">取引の編集</h1>
    <p v-if="loadError" class="text-red-600 text-sm">取引が見つかりませんでした。</p>
    <template v-else>
      <p v-if="generalError" class="text-red-600 text-sm mb-3">{{ generalError }}</p>
      <TransactionForm
        :initial="transaction"
        :submitting="submitting"
        :field-errors="fieldErrors"
        submit-label="更新する"
        @submit="onSubmit"
      />
    </template>
  </div>
</template>
