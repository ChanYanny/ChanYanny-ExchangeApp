import { shallowRef } from 'vue'

export interface ExchangeRateInput {
  fromCurrency: string
  toCurrency: string
  rate: number
}

export interface ExchangeRate extends ExchangeRateInput {
  _id: number
  date: string
}

export interface ArticleInput {
  Title: string
  Preview: string
  Content: string
}

export interface Article extends ArticleInput {
  ID: number
  CreatedAt: string
}

interface Credentials {
  username: string
  password: string
}

interface Session {
  username: string
  token: string
}

const sessionKey = 'exchange-app.session'
const apiBase = (import.meta.env.VITE_API_BASE_URL || '/api').replace(/\/+$/, '')

function restoreSession(): Session | null {
  try {
    const saved = JSON.parse(sessionStorage.getItem(sessionKey) ?? 'null')
    if (typeof saved?.username === 'string' && typeof saved?.token === 'string' && saved.token) {
      return { username: saved.username, token: saved.token }
    }
  } catch {
    // Browsers can block storage. Authentication still works for the current page.
  }
  return null
}

export const session = shallowRef<Session | null>(restoreSession())

export function setSession(value: Session | null) {
  session.value = value
  try {
    if (value) sessionStorage.setItem(sessionKey, JSON.stringify(value))
    else sessionStorage.removeItem(sessionKey)
  } catch {
    // Keep the in-memory session if storage is unavailable.
  }
}

export class ApiError extends Error {
  constructor(readonly status: number, message: string) {
    super(message)
  }
}

export function errorMessage(error: unknown, fallback = '请求未完成，请稍后重试。') {
  if (!(error instanceof ApiError)) return fallback
  switch (error.status) {
    case 0: return '无法连接后端，请确认服务已启动后重试。'
    case 400: return '提交的数据有误，请检查输入。'
    case 401: return '登录已失效，请重新登录。'
    case 404: return '未找到对应的数据。'
    case 502:
    case 503:
    case 504: return '后端服务暂时不可用，请稍后重试。'
    default: return fallback
  }
}

interface RequestOptions {
  body?: object
  authenticated?: boolean
  signal?: AbortSignal
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const currentSession = options.authenticated ? session.value : null
  if (options.authenticated && !currentSession) throw new ApiError(401, 'Missing session')
  const headers = new Headers({ Accept: 'application/json' })
  if (options.body) headers.set('Content-Type', 'application/json')
  if (currentSession) headers.set('Authorization', currentSession.token)

  let response: Response
  try {
    response = await fetch(apiBase + path, {
      method: options.body ? 'POST' : 'GET',
      body: options.body ? JSON.stringify(options.body) : undefined,
      headers,
      signal: options.signal,
    })
  } catch (error) {
    if (options.signal?.aborted) throw error
    throw new ApiError(0, 'Network request failed')
  }

  const data = await response.json().catch(() => undefined)
  if (!response.ok) {
    if (response.status === 401 && currentSession && session.value === currentSession) {
      setSession(null)
    }
    throw new ApiError(response.status, typeof data?.error === 'string' ? data.error : response.statusText)
  }
  if (data === undefined) throw new ApiError(502, 'Expected a JSON response')
  return data as T
}

export const api = {
  async authenticate(mode: 'login' | 'register', credentials: Credentials, signal: AbortSignal) {
    const { token } = await request<{ token: string }>(`/auth/${mode}`, { body: credentials, signal })
    if (signal.aborted) return
    if (typeof token !== 'string' || !token) throw new ApiError(502, 'Missing token')
    setSession({ username: credentials.username, token })
  },
  rates: (signal?: AbortSignal) => request<ExchangeRate[] | null>('/exchangeRates', { signal }),
  createRate: (body: ExchangeRateInput) =>
    request<ExchangeRate>('/exchangeRates', { body, authenticated: true }),
  articles: (signal?: AbortSignal) =>
    request<Article[] | null>('/articles', { authenticated: true, signal }),
  article: (id: string, signal?: AbortSignal) =>
    request<Article>(`/articles/${id}`, { authenticated: true, signal }),
  createArticle: (body: ArticleInput) =>
    request<Article>('/articles', { body, authenticated: true }),
  likes: (id: string, signal?: AbortSignal) =>
    request<{ likes: string }>(`/articles/${id}/likes`, { authenticated: true, signal }),
  like: (id: string, signal?: AbortSignal) =>
    request<{ message: string }>(`/articles/${id}/likes`, { body: {}, authenticated: true, signal }),
}
