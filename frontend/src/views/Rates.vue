<script setup lang="ts">
import { ArrowDownRight, ArrowLeftRight, ArrowRight, ChevronDown, Clock3, Database, Globe, Plus, RefreshCw, Search } from '@lucide/vue'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { api, errorMessage, type ExchangeRate } from '../api'
import { currencyName, formatDate, formatNumber } from '../format'

const rates = ref<ExchangeRate[]>([])
const loading = ref(true)
const error = ref('')
const search = ref('')
const source = ref('')
const target = ref('')
const amount = ref<number | string>(100)
let pending: AbortController | undefined

const sorted = computed(() => [...rates.value].sort((a, b) => {
  const aDate = Date.parse(a.date) || 0
  const bDate = Date.parse(b.date) || 0
  return bDate - aDate || b._id - a._id
}))
const codes = computed(() => [...new Set(sorted.value.flatMap(row => [row.fromCurrency, row.toCurrency]))].sort())
const pairKey = (from: string, to: string) => JSON.stringify([from, to])
const latest = computed(() => {
  const pairs = new Map<string, ExchangeRate>()
  for (const row of sorted.value) {
    const key = pairKey(row.fromCurrency, row.toCurrency)
    if (Number.isFinite(row.rate) && row.rate > 0 && !pairs.has(key)) pairs.set(key, row)
  }
  return pairs
})
const visible = computed(() => {
  const query = search.value.trim().toLowerCase()
  return sorted.value.filter(row => [row.fromCurrency, row.toCurrency, currencyName(row.fromCurrency), currencyName(row.toCurrency), `${row.fromCurrency}/${row.toCurrency}`]
    .some(value => value.toLowerCase().includes(query)))
})
const conversion = computed(() => {
  if (!source.value || !target.value) return null
  if (source.value === target.value) return { rate: 1, date: null }
  const direct = latest.value.get(pairKey(source.value, target.value))
  if (direct) return { rate: direct.rate, date: direct.date }
  const inverse = latest.value.get(pairKey(target.value, source.value))
  return inverse ? { rate: 1 / inverse.rate, date: inverse.date } : null
})
const result = computed(() => {
  if (typeof amount.value !== 'number' || amount.value < 0 || !conversion.value) return null
  const value = amount.value * conversion.value.rate
  return Number.isFinite(value) ? value : null
})

function swap() {
  const previous = source.value
  source.value = target.value
  target.value = previous
}

async function loadRates() {
  pending?.abort()
  const request = new AbortController()
  pending = request
  loading.value = true
  error.value = ''
  try {
    const rows = await api.rates(request.signal)
    if (request.signal.aborted) return
    rates.value = rows ?? []
    const first = latest.value.values().next().value
    if (!codes.value.includes(source.value)) source.value = first?.fromCurrency ?? codes.value[0] ?? ''
    if (!codes.value.includes(target.value)) target.value = first?.toCurrency ?? codes.value[0] ?? ''
  } catch (cause) {
    if (!request.signal.aborted) error.value = errorMessage(cause, '暂时无法载入汇率，请确认后端服务已启动。')
  } finally {
    if (!request.signal.aborted) loading.value = false
  }
}

onMounted(loadRates)
onUnmounted(() => pending?.abort())
</script>

<template>
  <section class="page-heading">
    <div>
      <h1>汇率与换算</h1>
      <p>浏览已记录的货币汇率，按最新记录换算金额。</p>
    </div>
    <RouterLink class="button primary" to="/rates/new"><Plus :size="18" aria-hidden="true" />新增汇率</RouterLink>
  </section>

  <div class="dashboard-top">
    <section class="panel converter" aria-labelledby="converter-title">
      <div class="panel-heading">
        <div class="section-icon"><ArrowLeftRight :size="20" aria-hidden="true" /></div>
        <div><h2 id="converter-title">货币换算</h2><p>选择币种并输入金额</p></div>
      </div>
      <div class="convert-fields">
        <label class="field">原币种
          <span class="select-field"><select v-model="source" :disabled="!codes.length">
            <option v-if="!codes.length" value="">暂无可用币种</option>
            <option v-for="code in codes" :key="code" :value="code">{{ code }} · {{ currencyName(code) }}</option>
          </select><ChevronDown :size="16" aria-hidden="true" /></span>
        </label>
        <button class="button icon-button swap" type="button" aria-label="交换币种" :disabled="!codes.length" @click="swap"><ArrowLeftRight :size="18" aria-hidden="true" /></button>
        <label class="field">目标币种
          <span class="select-field"><select v-model="target" :disabled="!codes.length">
            <option v-if="!codes.length" value="">暂无可用币种</option>
            <option v-for="code in codes" :key="code" :value="code">{{ code }} · {{ currencyName(code) }}</option>
          </select><ChevronDown :size="16" aria-hidden="true" /></span>
        </label>
      </div>
      <label class="field">金额
        <span class="amount-field"><input v-model.number="amount" type="number" min="0" step="any" inputmode="decimal" /><span>{{ source }}</span></span>
      </label>
      <div class="conversion-result" aria-live="polite">
        <span>换算结果</span>
        <strong v-if="result !== null">{{ formatNumber(result) }} <small>{{ target }}</small></strong>
        <strong v-else>—</strong>
        <span v-if="conversion && Number.isFinite(conversion.rate)">1 {{ source }} = {{ formatNumber(conversion.rate) }} {{ target }}<template v-if="conversion.date"> · {{ formatDate(conversion.date) }}</template></span>
        <span v-else-if="source && target">暂无这组币种的汇率记录。</span>
      </div>
    </section>

    <aside class="panel overview" aria-labelledby="overview-title">
      <div class="panel-heading">
        <div class="section-icon"><Globe :size="20" aria-hidden="true" /></div>
        <div><h2 id="overview-title">数据概览</h2><p>来自汇率数据库</p></div>
      </div>
      <div class="overview-line"><Database :size="17" aria-hidden="true" /><span>汇率记录</span><strong>{{ loading || error ? '—' : formatNumber(rates.length) }}</strong></div>
      <div class="overview-line"><ArrowDownRight :size="17" aria-hidden="true" /><span>货币对</span><strong>{{ loading || error ? '—' : formatNumber(latest.size) }}</strong></div>
      <div class="overview-line"><Clock3 :size="17" aria-hidden="true" /><span>最近记录</span><strong class="overview-date">{{ sorted[0] && !error ? formatDate(sorted[0].date) : '—' }}</strong></div>
    </aside>
  </div>

  <section class="panel rates-panel" aria-labelledby="rates-title">
    <div class="table-top">
      <div><h2 id="rates-title">汇率记录</h2><p v-if="!loading && !error">共 {{ visible.length }} 条</p></div>
      <div class="table-actions">
        <label class="search-field"><Search :size="17" aria-hidden="true" /><input v-model="search" type="search" aria-label="筛选币种" placeholder="搜索币种或货币对" /></label>
        <button class="button icon-button" type="button" aria-label="刷新汇率" :disabled="loading" @click="loadRates"><RefreshCw :size="17" aria-hidden="true" /></button>
      </div>
    </div>
    <p v-if="loading" class="empty-state" role="status">正在载入汇率…</p>
    <div v-else-if="error" class="empty-state"><p role="alert">{{ error }}</p><button class="button" type="button" @click="loadRates">重试</button></div>
    <p v-else-if="!visible.length" class="empty-state">{{ rates.length ? '没有符合条件的汇率。' : '暂无汇率记录。' }}</p>
    <div v-else class="table-scroll">
      <table>
        <thead><tr><th scope="col">货币对</th><th scope="col">汇率</th><th scope="col">记录时间</th></tr></thead>
        <tbody>
          <tr v-for="row in visible" :key="row._id">
            <td><span class="pair"><span>{{ row.fromCurrency }}</span><ArrowRight :size="14" aria-hidden="true" /><span>{{ row.toCurrency }}</span></span></td>
            <td class="numeric">{{ formatNumber(row.rate) }}</td>
            <td class="muted">{{ formatDate(row.date) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
