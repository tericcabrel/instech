type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'OPTIONS' | 'HEAD'

type RequestConfig = {
  url: string
  method: HttpMethod
  params?: Record<string, unknown>
  data?: unknown
  headers?: HeadersInit
  signal?: AbortSignal
}

export type HTTPError = {
  code: number
  message: string
  details: unknown
}

const defaultBaseUrl = '/api'

const resolveBaseUrl = () => import.meta.env.VITE_API_BASE_URL ?? defaultBaseUrl

const buildUrl = (url: string, params?: Record<string, unknown>) => {
  const normalizedPath = url.startsWith('/') ? url : `/${url}`
  const requestUrl = new URL(`${resolveBaseUrl()}${normalizedPath}`, window.location.origin)

  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined || value === null) {
        continue
      }

      if (Array.isArray(value)) {
        for (const item of value) {
          requestUrl.searchParams.append(key, String(item))
        }
        continue
      }

      requestUrl.searchParams.set(key, String(value))
    }
  }

  return requestUrl.toString()
}

const toHttpError = (status: number, payload: unknown): HTTPError => {
  if (
    typeof payload === 'object' &&
    payload !== null &&
    'code' in payload &&
    'message' in payload &&
    'details' in payload
  ) {
    return payload as HTTPError
  }

  return {
    code: status,
    message: `Request failed with status ${status}`,
    details: payload,
  }
}

export class ApiError extends Error {
  public readonly status: number
  public readonly body: HTTPError

  constructor(status: number, body: HTTPError) {
    super(body.message)
    this.name = 'ApiError'
    this.status = status
    this.body = body
  }
}

export const customInstance = async <T>(config: RequestConfig): Promise<T> => {
  const headers = new Headers(config.headers)
  const hasBody = config.data !== undefined

  if (hasBody && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }

  if (!headers.has('Accept')) {
    headers.set('Accept', 'application/json')
  }

  const response = await fetch(buildUrl(config.url, config.params), {
    method: config.method,
    headers,
    body: hasBody ? JSON.stringify(config.data) : undefined,
    signal: config.signal,
  })

  if (response.status === 204) {
    return undefined as T
  }

  const text = await response.text()
  const parsed = text ? (JSON.parse(text) as unknown) : undefined

  if (!response.ok) {
    throw new ApiError(response.status, toHttpError(response.status, parsed))
  }

  return parsed as T
}
