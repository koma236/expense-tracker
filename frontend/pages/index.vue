<script setup lang="ts">
// ルートは取引一覧へリダイレクト（プロトタイプの入口）。
await navigateTo('/transactions', { redirectCode: 302 })
</script>

<template>
  <div />
</template>
