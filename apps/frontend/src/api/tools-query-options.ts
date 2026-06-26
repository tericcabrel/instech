import { queryOptions } from '@tanstack/react-query'
import {
  getToolsSearch,
  getToolsId,
  getToolsIdAlternatives,
  getToolsIdGraph,
} from './generated/instech'
import { Tool, type ToolOutput } from './generated/model/Tool.zod.ts'
import {
  ToolAlternative,
  type ToolAlternativeOutput,
} from './generated/model/ToolAlternative.zod.ts'
import {
  ToolGraphResponse,
  type ToolGraphResponseOutput,
} from './generated/model/ToolGraphResponse.zod.ts'
import { ToolSearchResultItem } from './generated/model/ToolSearchResultItem.zod.ts'

export const toolKeys = {
  all: ['tools'] as const,
  search: (keyword: string) => [...toolKeys.all, 'search', keyword.trim()] as const,
  detail: (slug: string) => [...toolKeys.all, 'detail', slug] as const,
  alternatives: (slug: string) => [...toolKeys.detail(slug), 'alternatives'] as const,
  graph: (slug: string) => [...toolKeys.detail(slug), 'graph'] as const,
}

export const toolsQueryOptions = (keyword: string) =>
  queryOptions({
    queryKey: toolKeys.search(keyword),
    enabled: keyword.trim().length > 0,
    queryFn: async (): Promise<ToolSearchResultItem[]> => {
      const response = await getToolsSearch({ q: keyword.trim() })

      return ToolSearchResultItem.array().parse(response)
    },
    retry: false,
  })

export const toolDetailQueryOptions = (slug: string) =>
  queryOptions({
    queryKey: toolKeys.detail(slug),
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolOutput> => {
      const response = await getToolsId(slug)

      return Tool.parse(response)
    },
    retry: false,
  })

export const toolAlternativesQueryOptions = (slug: string) =>
  queryOptions({
    queryKey: toolKeys.alternatives(slug),
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolAlternativeOutput[]> => {
      const response = await getToolsIdAlternatives(slug)

      return ToolAlternative.array().parse(response)
    },
    retry: false,
  })

export const toolGraphQueryOptions = (slug: string) =>
  queryOptions({
    queryKey: toolKeys.graph(slug),
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolGraphResponseOutput> => {
      const response = await getToolsIdGraph(slug)

      return ToolGraphResponse.parse(response)
    },
    retry: false,
  })
