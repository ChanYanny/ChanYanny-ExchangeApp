import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import { session } from './api'
import Rates from './views/Rates.vue'
import './style.css'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior: () => ({ top: 0 }),
  routes: [
    { path: '/', name: 'rates', component: Rates },
    { path: '/rates/new', name: 'rate-new', component: () => import('./views/NewRate.vue'), meta: { auth: true } },
    { path: '/login', name: 'login', component: () => import('./views/Auth.vue') },
    { path: '/register', name: 'register', component: () => import('./views/Auth.vue') },
    { path: '/articles', name: 'articles', component: () => import('./views/Articles.vue'), meta: { auth: true } },
    { path: '/articles/new', name: 'article-new', component: () => import('./views/NewArticle.vue'), meta: { auth: true } },
    { path: '/articles/:id([1-9]\\d*)', name: 'article', component: () => import('./views/Article.vue'), meta: { auth: true } },
    { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('./views/NotFound.vue') },
  ],
})

router.beforeEach(to => {
  if (to.meta.auth && !session.value) return { name: 'login', query: { next: to.fullPath } }
})

createApp(App).use(router).mount('#app')
