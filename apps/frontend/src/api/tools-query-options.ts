import { queryOptions } from '@tanstack/react-query'
import {
  getToolsSearch,
  getToolsId,
  getToolsIdAlternatives,
  getToolsIdGraph,
} from './generated/instech'
import type { GetToolsIdGraphParamsOutput } from './generated/model/GetToolsIdGraphParams.zod.ts'
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
  graph: (
    slug: string,
    params: Pick<GetToolsIdGraphParamsOutput, 'depth' | 'kinds' | 'layoutMode'>,
  ) => [...toolKeys.detail(slug), 'graph', params] as const,
}

type ToolGraphQueryInput = {
  slug: string
  depth?: GetToolsIdGraphParamsOutput['depth']
  kinds?: GetToolsIdGraphParamsOutput['kinds']
  layoutMode?: GetToolsIdGraphParamsOutput['layoutMode']
}

const normalizeGraphParams = ({
  depth,
  kinds,
  layoutMode,
}: Omit<ToolGraphQueryInput, 'slug'>): GetToolsIdGraphParamsOutput => ({
  depth: depth ?? 1,
  kinds: kinds?.length ? [...kinds].sort() : undefined,
  layoutMode,
})

export const toolSearchQueryOptions = (q: string) =>
  queryOptions({
    queryKey: toolKeys.search(q),
    enabled: q.trim().length > 0,
    queryFn: async (): Promise<ToolSearchResultItem[]> => {
      const response = await getToolsSearch({ q: q.trim() })

      return ToolSearchResultItem.array().parse(response)
    },
    retry: false,
  })

export const toolQueryOptions = (slug: string) =>
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

export const toolGraphQueryOptions = ({
  slug,
  depth,
  kinds,
  layoutMode,
}: ToolGraphQueryInput) => {
  const graphParams = normalizeGraphParams({ depth, kinds, layoutMode })

  return queryOptions({
    queryKey: toolKeys.graph(slug, graphParams),
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolGraphResponseOutput> => {
      const response = await getToolsIdGraph(slug, graphParams)

      return ToolGraphResponse.parse(response)
    },
    retry: false,
  })
}
