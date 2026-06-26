import { queryOptions } from '@tanstack/react-query'
import {
  getToolsSearch,
  getToolsId,
  getToolsIdAlternatives,
  getToolsIdGraph,
} from './generated/instech'
import type {
  GetToolsIdGraphParams,
  Tool,
  ToolAlternative,
  ToolGraphResponse,
  ToolSearchResultItem,
} from './generated/model'

export const toolKeys = {
  all: ['tools'] as const,
  search: (keyword: string) => [...toolKeys.all, 'search', keyword.trim()] as const,
  detail: (slug: string) => [...toolKeys.all, 'detail', slug] as const,
  alternatives: (slug: string) => [...toolKeys.detail(slug), 'alternatives'] as const,
  graph: (
    slug: string,
    params?: {
      depth?: 1 | 2
      kinds?: string[]
      layoutMode?: 'chronological' | 'force'
    }
  ) => [...toolKeys.detail(slug), 'graph', params ?? {}] as const,
}

export const toolsQueryOptions = (keyword: string) =>
  queryOptions({
    queryKey: toolKeys.search(keyword),
    enabled: keyword.trim().length > 0,
    queryFn: async (): Promise<ToolSearchResultItem[]> =>
      (await getToolsSearch({ q: keyword.trim() })) as unknown as ToolSearchResultItem[],
    retry: false,
  })

export const toolDetailQueryOptions = (slug: string) =>
  queryOptions({
    queryKey: toolKeys.detail(slug),
    enabled: Boolean(slug),
    queryFn: async (): Promise<Tool> => (await getToolsId(slug)) as unknown as Tool,
    retry: false,
  })

export const toolAlternativesQueryOptions = (slug: string) =>
  queryOptions({
    queryKey: toolKeys.alternatives(slug),
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolAlternative[]> =>
      (await getToolsIdAlternatives(slug)) as unknown as ToolAlternative[],
    retry: false,
  })

export const toolGraphQueryOptions = (
  slug: string,
  params?: GetToolsIdGraphParams
) =>
  queryOptions({
    queryKey: toolKeys.graph(slug, params),
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolGraphResponse> =>
      (await getToolsIdGraph(slug, params)) as unknown as ToolGraphResponse,
    retry: false,
  })
