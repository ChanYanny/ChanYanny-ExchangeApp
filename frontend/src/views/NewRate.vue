<script setup lang="ts">
import { ArrowLeft, ArrowRight, Plus } from '@lucide/vue'
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { api, errorMessage } from '../api'

const router = useRouter()
const source = ref('')
const target = ref('')
const rate = ref<number | string>('')
const busy = ref(false)
const error = ref('')

async function submit() {
  if (busy.value) return
  const fromCurrency = source.value.trim().toUpperCase()
  const toCurrency = target.value.trim().toUpperCase()
  const value = rate.value
  if (!/^[A-Z]{3}$/.test(fromCurrency) || !/^[A-Z]{3}$/.test(toCurrency) || fromCurrency === toCurrency) {
    error.value = '请输入两个不同的三字母货币代码。'
    return
  }
  if (typeof value !== 'number' || !Number.isFinite(value) || value <= 0) {
    error.value = '汇率必须是大于零的数字。'
    return
  }
  busy.value = true
  error.value = ''
  try {
    await api.createRate({ fromCurrency, toCurrency, rate: value })
    if (router.currentRoute.value.name === 'rate-new') await router.push({ name: 'rates' })
  } catch (cause) {
    error.value = errorMessage(cause, '汇率保存失败，请稍后重试。')
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <RouterLink class="back-link" to="/"><ArrowLeft :size="17" aria-hidden="true" />返回汇率</RouterLink>
  <section class="panel form-panel">
    <div class="section-icon"><Plus :size="21" aria-hidden="true" /></div>
    <h1>新增汇率</h1>
    <p class="muted">记录 1 单位原币种可兑换的目标币种金额。</p>
    <form class="form-stack" @submit.prevent="submit">
      <div class="two-fields">
        <label class="field">原币种<input v-model="source" class="currency-input" placeholder="USD" maxlength="3" required /></label>
        <label class="field">目标币种<input v-model="target" class="currency-input" placeholder="CNY" maxlength="3" required /></label>
      </div>
      <label class="field">汇率<input v-model.number="rate" type="number" step="any" inputmode="decimal" placeholder="例如 7.12" required /></label>
      <p v-if="error" class="notice error" role="alert">{{ error }}</p>
      <div class="form-actions"><button class="button primary" type="submit" :disabled="busy">{{ busy ? '正在保存…' : '保存汇率' }}<ArrowRight :size="17" aria-hidden="true" /></button></div>
    </form>
  </section>
</template>
