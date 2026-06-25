import { queryOptions } from '@tanstack/react-query'
import {
  getToolsId,
  getToolsIdAlternatives,
  getToolsIdGraph,
} from './generated/instech'
import {
  Tool as ToolSchema,
  ToolAlternative as ToolAlternativeSchema,
  ToolGraphResponse as ToolGraphResponseSchema,
} from './generated/model'
import type {
  ToolAlternativeOutput,
  ToolGraphResponseOutput,
  ToolOutput,
} from './generated/model'

export const toolKeys = {
  all: ['tools'] as const,
  detail: (slug: string) => [...toolKeys.all, 'detail', slug] as const,
  alternatives: (slug: string) => [...toolKeys.detail(slug), 'alternatives'] as const,
  graph: (slug: string) => [...toolKeys.detail(slug), 'graph'] as const,
}

export const toolDetailQueryOptions = (slug: string) =>
  queryOptions({
    queryKey: toolKeys.detail(slug),
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolOutput> => {
      const response = await getToolsId(slug)

      return ToolSchema.parse(response)
    },
    retry: false,
  })

export const toolAlternativesQueryOptions = (slug: string) =>
  queryOptions({
    queryKey: toolKeys.alternatives(slug),
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolAlternativeOutput[]> => {
      const response = await getToolsIdAlternatives(slug)

      return ToolAlternativeSchema.array().parse(response)
    },
    retry: false,
  })

export const toolGraphQueryOptions = (slug: string) =>
  queryOptions({
    queryKey: toolKeys.graph(slug),
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolGraphResponseOutput> => {
      const response = await getToolsIdGraph(slug)

      return ToolGraphResponseSchema.parse(response)
    },
    retry: false,
  })
