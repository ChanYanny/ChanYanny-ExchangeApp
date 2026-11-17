const numberFormat = new Intl.NumberFormat('zh-CN', { maximumSignificantDigits: 12 })
const dateFormat = new Intl.DateTimeFormat('zh-CN', {
  year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit',
})
const currencyNames = new Intl.DisplayNames(['zh-CN'], { type: 'currency', fallback: 'none' })

export const formatNumber = (value: number) => numberFormat.format(value)

export function formatDate(value: string) {
  const time = Date.parse(value)
  return Number.isFinite(time) ? dateFormat.format(time) : '—'
}

export function currencyName(code: string) {
  return /^[A-Z]{3}$/.test(code) ? currencyNames.of(code) ?? code : code
}
