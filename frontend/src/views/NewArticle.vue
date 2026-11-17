<script setup lang="ts">
import { ArrowLeft, ArrowRight, PenLine } from '@lucide/vue'
import { reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { api, ApiError, errorMessage } from '../api'

const router = useRouter()
const draft = reactive({ Title: '', Preview: '', Content: '' })
const busy = ref(false)
const error = ref('')

async function submit() {
  if (busy.value) return
  const body = { Title: draft.Title.trim(), Preview: draft.Preview.trim(), Content: draft.Content.trim() }
  if (!body.Title || !body.Preview || !body.Content) {
    error.value = '请填写标题、摘要和正文。'
    return
  }
  busy.value = true
  error.value = ''
  try {
    const article = await api.createArticle(body)
    if (router.currentRoute.value.name === 'article-new') {
      await router.push({ name: 'article', params: { id: article.ID } })
    }
  } catch (cause) {
    error.value = cause instanceof ApiError && (cause.status === 0 || cause.status >= 500)
      ? '暂时无法确认发布结果。请先查看文章列表，再决定是否重试。'
      : errorMessage(cause)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <RouterLink class="back-link" to="/articles"><ArrowLeft :size="17" aria-hidden="true" />返回文章</RouterLink>
  <section class="panel form-panel editor-panel">
    <div class="section-icon"><PenLine :size="21" aria-hidden="true" /></div>
    <h1>发布文章</h1>
    <form class="form-stack" @submit.prevent="submit">
      <label class="field">标题<input v-model="draft.Title" required /></label>
      <label class="field">摘要<textarea v-model="draft.Preview" rows="3" required /></label>
      <label class="field">正文<textarea v-model="draft.Content" rows="10" required /></label>
      <p v-if="error" class="notice error" role="alert">{{ error }}</p>
      <div class="form-actions"><button class="button primary" type="submit" :disabled="busy">{{ busy ? '正在发布…' : '发布文章' }}<ArrowRight :size="17" aria-hidden="true" /></button></div>
    </form>
  </section>
</template>
