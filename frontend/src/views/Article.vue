<script setup lang="ts">
import { ArrowLeft, Heart, RefreshCw } from '@lucide/vue'
import { onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { api, ApiError, errorMessage, type Article } from '../api'
import { formatDate, formatNumber } from '../format'

const id = String(useRoute().params.id)
const lifetime = new AbortController()
const article = ref<Article>()
const likes = ref<number | null>(null)
const loading = ref(true)
const checkingLikes = ref(false)
const sendingLike = ref(false)
const articleError = ref('')
const countError = ref('')
const likeFeedback = ref('')

async function loadArticle() {
  loading.value = true
  articleError.value = ''
  try {
    article.value = await api.article(id, lifetime.signal)
    await loadLikes()
  } catch (cause) {
    if (!lifetime.signal.aborted) articleError.value = errorMessage(cause, '文章加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

async function loadLikes() {
  if (checkingLikes.value) return
  checkingLikes.value = true
  countError.value = ''
  try {
    const data = await api.likes(id, lifetime.signal)
    const count = Number(data.likes)
    if (!Number.isSafeInteger(count) || count < 0) throw new Error('Invalid like count')
    likes.value = count
  } catch (cause) {
    if (!lifetime.signal.aborted) {
      likes.value = null
      countError.value = errorMessage(cause, '暂时无法获取点赞数。')
    }
  } finally {
    checkingLikes.value = false
  }
}

async function like() {
  if (sendingLike.value || checkingLikes.value) return
  sendingLike.value = true
  likeFeedback.value = ''
  try {
    await api.like(id, lifetime.signal)
    likeFeedback.value = '已点赞。'
  } catch (cause) {
    if (lifetime.signal.aborted) return
    likeFeedback.value = cause instanceof ApiError && (cause.status === 0 || cause.status >= 500)
      ? '暂时无法确认点赞结果。'
      : errorMessage(cause, '点赞失败，请稍后重试。')
  } finally {
    if (!lifetime.signal.aborted) await loadLikes()
    sendingLike.value = false
  }
}

onMounted(loadArticle)
onUnmounted(() => lifetime.abort())
</script>

<template>
  <RouterLink class="back-link" to="/articles"><ArrowLeft :size="17" aria-hidden="true" />返回文章</RouterLink>
  <article v-if="article" class="panel article-detail">
    <time :datetime="article.CreatedAt">{{ formatDate(article.CreatedAt) }}</time>
    <h1>{{ article.Title }}</h1>
    <p class="article-preview">{{ article.Preview }}</p>
    <div class="article-body">{{ article.Content }}</div>
    <div class="article-footer">
      <button class="button" type="button" :disabled="sendingLike || checkingLikes" @click="like">
        <Heart :size="18" aria-hidden="true" />{{ sendingLike ? '正在点赞…' : '点赞' }}
        <span v-if="likes !== null">{{ formatNumber(likes) }}</span>
      </button>
      <button v-if="countError" class="button compact" type="button" @click="loadLikes"><RefreshCw :size="15" aria-hidden="true" />重试获取点赞数</button>
      <span class="muted" role="status">{{ countError || likeFeedback || (checkingLikes ? '正在获取点赞数…' : '') }}</span>
    </div>
  </article>
  <p v-else-if="loading" class="panel empty-state" role="status">正在载入文章…</p>
  <div v-else class="panel empty-state"><p role="alert">{{ articleError }}</p><button class="button" type="button" @click="loadArticle">重试</button></div>
</template>
