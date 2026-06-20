<script setup lang="ts">
type Health = { status: string; db: string }

const api = useApi()
const { data, error, status } = await useAsyncData('health', () =>
  api.get<Health>('/health'),
)
</script>

<template>
  <main class="container">
    <h1>ExpenseTracker</h1>
    <p>家計簿アプリ（学習用課題）の土台です。</p>

    <section class="card">
      <h2>バックエンド疎通確認</h2>
      <p v-if="status === 'pending'">確認中...</p>
      <p v-else-if="error" class="ng">
        API に接続できませんでした。バックエンド（:8080）と MySQL が起動しているか確認してください。
      </p>
      <ul v-else>
        <li>status: <strong>{{ data?.status }}</strong></li>
        <li>db: <strong>{{ data?.db }}</strong></li>
      </ul>
    </section>
  </main>
</template>

<style scoped>
.container {
  font-family: system-ui, sans-serif;
  max-width: 640px;
  margin: 40px auto;
  padding: 0 16px;
}
.card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 16px 20px;
  margin-top: 24px;
}
.ng {
  color: #c0392b;
}
</style>
