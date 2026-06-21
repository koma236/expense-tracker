<script setup lang="ts">
const { label, prev, next } = useMonth()
</script>

<template>
  <div class="min-h-screen bg-gray-50 text-gray-800">
    <header class="bg-white border-b shadow-sm">
      <div class="mx-auto max-w-4xl px-4 py-3 flex items-center justify-between">
        <NuxtLink to="/transactions" class="text-lg font-bold text-indigo-600">
          家計簿
        </NuxtLink>

        <!-- 月セレクタ -->
        <div class="flex items-center gap-2">
          <button
            class="px-2 py-1 rounded hover:bg-gray-100 text-gray-600"
            aria-label="前の月"
            @click="prev"
          >
            ←
          </button>
          <span class="w-24 text-center font-medium tabular-nums">{{ label }}</span>
          <button
            class="px-2 py-1 rounded hover:bg-gray-100 text-gray-600"
            aria-label="次の月"
            @click="next"
          >
            →
          </button>
        </div>
      </div>

      <!-- ナビ -->
      <nav class="mx-auto max-w-4xl px-4">
        <ul class="flex gap-1 text-sm">
          <li>
            <NuxtLink
              to="/transactions"
              class="inline-block px-3 py-2 border-b-2 border-transparent hover:text-indigo-600"
              active-class="!border-indigo-600 text-indigo-600 font-medium"
            >
              取引一覧
            </NuxtLink>
          </li>
          <li>
            <NuxtLink
              to="/categories"
              class="inline-block px-3 py-2 border-b-2 border-transparent hover:text-indigo-600"
              active-class="!border-indigo-600 text-indigo-600 font-medium"
            >
              カテゴリ
            </NuxtLink>
          </li>
          <li>
            <NuxtLink
              to="/dashboard"
              class="inline-block px-3 py-2 border-b-2 border-transparent hover:text-indigo-600"
              active-class="!border-indigo-600 text-indigo-600 font-medium"
            >
              ダッシュボード
            </NuxtLink>
          </li>
          <li>
            <span class="inline-block px-3 py-2 text-gray-300 cursor-not-allowed" title="今後実装予定">
              予算
            </span>
          </li>
        </ul>
      </nav>
    </header>

    <main class="mx-auto max-w-4xl px-4 py-6">
      <slot />
    </main>
  </div>
</template>
