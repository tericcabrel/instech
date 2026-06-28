type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'OPTIONS' | 'HEAD';

type RequestConfig = {
  data?: unknown;
  headers?: HeadersInit;
  method: HttpMethod;
  params?: Record<string, unknown>;
  signal?: AbortSignal;
  url: string;
};

type OrvalRequestOptions = RequestInit & {
  params?: Record<string, unknown>;
};

type HttpError = {
  code: number;
  details: unknown;
  message: string;
};

const DEFAULT_BASE_URL = '/api';
const HTTP_NO_CONTENT_STATUS = 204;

const resolveBaseUrl = () => import.meta.env.VITE_API_BASE_URL ?? DEFAULT_BASE_URL;

const buildUrl = (url: string, params?: Record<string, unknown>) => {
  const normalizedPath = url.startsWith('/') ? url : `/${url}`;
  const requestUrl = new URL(`${resolveBaseUrl()}${normalizedPath}`, window.location.origin);

  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined || value === null) {
        continue;
      }

      if (Array.isArray(value)) {
        for (const item of value) {
          requestUrl.searchParams.append(key, String(item));
        }
        continue;
      }

      requestUrl.searchParams.set(key, String(value));
    }
  }

  return requestUrl.toString();
};

const toHttpError = (status: number, payload: unknown): HttpError => {
  if (
    typeof payload === 'object' &&
    payload !== null &&
    'code' in payload &&
    'message' in payload &&
    'details' in payload
  ) {
    return payload as HttpError;
  }

  return {
    code: status,
    details: payload,
    message: `Request failed with status ${status}`,
  };
};

class ApiError extends Error {
  readonly status: number;
  readonly body: HttpError;

  constructor(status: number, body: HttpError) {
    super(body.message);
    this.name = 'ApiError';
    this.status = status;
    this.body = body;
  }
}

const parseResponseBody = async (response: Response): Promise<unknown> => {
  const text = await response.text();

  if (!text) {
    return;
  }

  try {
    return JSON.parse(text) as unknown;
  } catch {
    return text;
  }
};

const buildBody = (data?: unknown) => {
  if (!data) {
    return;
  }

  if (typeof data === 'string' || data instanceof FormData || data instanceof Blob || data instanceof URLSearchParams) {
    return data;
  }

  return JSON.stringify(data);
};

const requestWithConfig = async <T>(config: RequestConfig): Promise<T> => {
  const headers = new Headers(config.headers);
  const hasBody = config.data !== undefined;
  const body = buildBody(config.data);

  if (hasBody && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  if (!headers.has('Accept')) {
    headers.set('Accept', 'application/json');
  }

  const response = await fetch(buildUrl(config.url, config.params), {
    body,
    headers,
    method: config.method,
    signal: config.signal,
  });

  if (response.status === HTTP_NO_CONTENT_STATUS) {
    return undefined as T;
  }

  const parsed = await parseResponseBody(response);

  if (!response.ok) {
    throw new ApiError(response.status, toHttpError(response.status, parsed));
  }

  return parsed as T;
};

export const customInstance = <T>(urlOrConfig: string | RequestConfig, options?: OrvalRequestOptions): Promise<T> => {
  if (typeof urlOrConfig === 'string') {
    const config: RequestConfig = {
      data: options?.body,
      headers: options?.headers,
      method: (options?.method as HttpMethod | undefined) ?? 'GET',
      params: options?.params,
      signal: options?.signal ?? undefined,
      url: urlOrConfig,
    };

    return requestWithConfig<T>(config);
  }

  return requestWithConfig<T>(urlOrConfig);
};
