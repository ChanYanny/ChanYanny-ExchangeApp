<script setup lang="ts">
import { ArrowRight, UserRound } from '@lucide/vue'
import { computed, onUnmounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { api, ApiError, errorMessage } from '../api'

const route = useRoute()
const router = useRouter()
const registering = computed(() => route.name === 'register')
const username = ref('')
const password = ref('')
const busy = ref(false)
const error = ref('')
const lifetime = new AbortController()

onUnmounted(() => lifetime.abort())

async function submit() {
  if (busy.value) return
  if (!username.value.trim() || !password.value) {
    error.value = '请输入用户名和密码。'
    return
  }
  busy.value = true
  error.value = ''
  try {
    await api.authenticate(registering.value ? 'register' : 'login', {
      username: username.value.trim(),
      password: password.value,
    }, lifetime.signal)
    if (lifetime.signal.aborted) return
    const next = typeof route.query.next === 'string' ? router.resolve(route.query.next) : null
    await router.replace(next?.meta.auth ? next.fullPath : '/')
  } catch (cause) {
    if (lifetime.signal.aborted) return
    error.value = cause instanceof ApiError && cause.status === 401
      ? '用户名或密码错误。'
      : errorMessage(cause, registering.value ? '注册失败，请检查用户名或稍后重试。' : '登录失败，请稍后重试。')
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <section class="panel form-panel auth-panel">
    <div class="section-icon"><UserRound :size="21" aria-hidden="true" /></div>
    <h1>{{ registering ? '创建账户' : '登录' }}</h1>
    <p v-if="!registering && route.query.expired === '1'" class="notice" role="status">登录已失效，请重新登录。</p>
    <p v-else-if="route.query.next" class="muted">{{ registering ? '注册后继续访问。' : '登录后继续访问。' }}</p>
    <form class="form-stack" @submit.prevent="submit">
      <label class="field">用户名<input v-model="username" autocomplete="username" required /></label>
      <label class="field">密码<input v-model="password" type="password" :autocomplete="registering ? 'new-password' : 'current-password'" required /></label>
      <p v-if="error" class="notice error" role="alert">{{ error }}</p>
      <button class="button primary" type="submit" :disabled="busy">
        {{ busy ? '正在提交…' : registering ? '注册并登录' : '登录' }}
        <ArrowRight :size="17" aria-hidden="true" />
      </button>
    </form>
    <div class="auth-switch">
      {{ registering ? '已有账户？' : '还没有账户？' }}
      <RouterLink :to="{ name: registering ? 'login' : 'register', query: route.query }">{{ registering ? '登录' : '注册' }}</RouterLink>
    </div>
  </section>
</template>
