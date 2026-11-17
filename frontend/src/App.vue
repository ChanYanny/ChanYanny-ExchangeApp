<script setup lang="ts">
import { ArrowLeftRight, BookOpen, LogOut } from '@lucide/vue'
import { watch } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { session, setSession } from './api'

const route = useRoute()
const router = useRouter()

watch(session, value => {
  if (!value && route.meta.auth) {
    void router.replace({ name: 'login', query: { next: route.fullPath, expired: '1' } })
  }
})

async function logout() {
  if (route.meta.auth) await router.push({ name: 'rates' })
  setSession(null)
}
</script>

<template>
  <a class="skip-link" href="#main">跳至内容</a>
  <header class="site-header">
    <div class="header-inner">
      <RouterLink class="brand" to="/" aria-label="ExchangeApp 首页">
        <span class="brand-icon"><ArrowLeftRight :size="20" aria-hidden="true" /></span>
        <span>Exchange<span class="brand-accent">App</span></span>
      </RouterLink>
      <nav class="navigation" aria-label="主导航">
        <RouterLink to="/" :class="{ current: route.name === 'rate-new' }"><ArrowLeftRight :size="16" aria-hidden="true" />汇率</RouterLink>
        <RouterLink to="/articles" :class="{ current: route.path.startsWith('/articles/') }"><BookOpen :size="16" aria-hidden="true" />文章</RouterLink>
      </nav>
      <div class="account">
        <template v-if="session">
          <span class="username" :title="session.username">{{ session.username }}</span>
          <button class="button compact" type="button" @click="logout"><LogOut :size="15" aria-hidden="true" />退出</button>
        </template>
        <RouterLink v-else class="button compact" to="/login">登录</RouterLink>
      </div>
    </div>
  </header>
  <main id="main" class="page-shell" tabindex="-1">
    <RouterView v-slot="{ Component, route: activeRoute }">
      <component :is="Component" :key="activeRoute.path" />
    </RouterView>
  </main>
</template>
