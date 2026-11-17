<script setup lang="ts">
import { ArrowUpRight, BookOpen, PenLine } from '@lucide/vue'
import { onMounted, onUnmounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { api, errorMessage, type Article } from '../api'
import { formatDate } from '../format'

const articles = ref<Article[]>([])
const loading = ref(true)
const error = ref('')
let pending: AbortController | undefined

async function loadArticles() {
  pending?.abort()
  const request = new AbortController()
  pending = request
  loading.value = true
  error.value = ''
  try {
    const rows = await api.articles(request.signal)
    if (!request.signal.aborted) articles.value = [...(rows ?? [])].sort((a, b) =>
      Date.parse(b.CreatedAt) - Date.parse(a.CreatedAt) || b.ID - a.ID)
  } catch (cause) {
    if (!request.signal.aborted) error.value = errorMessage(cause, '文章加载失败，请稍后重试。')
  } finally {
    if (!request.signal.aborted) loading.value = false
  }
}

onMounted(loadArticles)
onUnmounted(() => pending?.abort())
</script>

<template>
  <section class="page-heading">
    <div><h1>文章</h1><p>阅读和分享有关外汇与货币的文章。</p></div>
    <RouterLink class="button primary" to="/articles/new"><PenLine :size="17" aria-hidden="true" />发布文章</RouterLink>
  </section>
  <p v-if="loading" class="panel empty-state" role="status">正在载入文章…</p>
  <div v-else-if="error" class="panel empty-state"><p role="alert">{{ error }}</p><button class="button" type="button" @click="loadArticles">重试</button></div>
  <div v-else-if="!articles.length" class="panel empty-state"><BookOpen :size="24" aria-hidden="true" /><p>暂无文章。</p></div>
  <div v-else class="article-grid">
    <article v-for="article in articles" :key="article.ID" class="panel article-card">
      <time :datetime="article.CreatedAt">{{ formatDate(article.CreatedAt) }}</time>
      <h2><RouterLink :to="{ name: 'article', params: { id: article.ID } }">{{ article.Title }}<ArrowUpRight :size="18" aria-hidden="true" /></RouterLink></h2>
      <p>{{ article.Preview }}</p>
    </article>
  </div>
</template>
