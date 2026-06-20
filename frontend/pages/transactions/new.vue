<script setup lang="ts">
import type { TransactionInput } from '~/types'

const api = useExpenseApi()
const submitting = ref(false)
const fieldErrors = ref<Record<string, string>>({})
const generalError = ref('')

async function onSubmit(input: TransactionInput) {
  submitting.value = true
  fieldErrors.value = {}
  generalError.value = ''
  try {
    await api.createTransaction(input)
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
    <h1 class="text-xl font-bold mb-4">取引の登録</h1>
    <p v-if="generalError" class="text-red-600 text-sm mb-3">{{ generalError }}</p>
    <TransactionForm
      :submitting="submitting"
      :field-errors="fieldErrors"
      submit-label="登録する"
      @submit="onSubmit"
    />
  </div>
</template>
