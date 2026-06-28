import { queryOptions } from '@tanstack/react-query';
import { getToolsId, getToolsIdAlternatives, getToolsIdGraph, getToolsSearch } from './generated/instech';
import type { GetToolsIdGraphParamsOutput } from './generated/model/GetToolsIdGraphParams.zod.ts';
import { Tool, type ToolOutput } from './generated/model/Tool.zod.ts';
import { ToolAlternative, type ToolAlternativeOutput } from './generated/model/ToolAlternative.zod.ts';
import { ToolGraphResponse, type ToolGraphResponseOutput } from './generated/model/ToolGraphResponse.zod.ts';
import { ToolSearchResultItem } from './generated/model/ToolSearchResultItem.zod.ts';

type ToolGraphQueryInput = {
  depth?: GetToolsIdGraphParamsOutput['depth'];
  kinds?: GetToolsIdGraphParamsOutput['kinds'];
  layoutMode?: GetToolsIdGraphParamsOutput['layoutMode'];
  slug: string;
};

const normalizeGraphParams = ({
  depth,
  kinds,
  layoutMode,
}: Omit<ToolGraphQueryInput, 'slug'>): GetToolsIdGraphParamsOutput => ({
  depth: depth ?? 1,
  kinds: kinds && kinds.length > 0 ? [...kinds].sort() : undefined,
  layoutMode,
});

export const toolKeys = {
  all: ['tools'] as const,
  alternatives: (slug: string) => [...toolKeys.detail(slug), 'alternatives'] as const,
  detail: (slug: string) => [...toolKeys.all, 'detail', slug] as const,
  graph: (slug: string, params: Pick<GetToolsIdGraphParamsOutput, 'depth' | 'kinds' | 'layoutMode'>) =>
    [...toolKeys.detail(slug), 'graph', params] as const,
  search: (keyword: string) => [...toolKeys.all, 'search', keyword.trim()] as const,
};

export const toolSearchQueryOptions = (q: string) =>
  queryOptions({
    enabled: q.trim().length > 0,
    queryFn: async (): Promise<ToolSearchResultItem[]> => {
      const response = await getToolsSearch({ q: q.trim() });

      return ToolSearchResultItem.array().parse(response);
    },
    queryKey: toolKeys.search(q),
    retry: false,
  });

export const toolQueryOptions = (slug: string) =>
  queryOptions({
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolOutput> => {
      const response = await getToolsId(slug);

      return Tool.parse(response);
    },
    queryKey: toolKeys.detail(slug),
    retry: false,
  });

export const toolAlternativesQueryOptions = (slug: string) =>
  queryOptions({
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolAlternativeOutput[]> => {
      const response = await getToolsIdAlternatives(slug);

      return ToolAlternative.array().parse(response);
    },
    queryKey: toolKeys.alternatives(slug),
    retry: false,
  });

export const toolGraphQueryOptions = ({ depth, kinds, layoutMode, slug }: ToolGraphQueryInput) => {
  const graphParams = normalizeGraphParams({ depth, kinds, layoutMode });

  return queryOptions({
    enabled: Boolean(slug),
    queryFn: async (): Promise<ToolGraphResponseOutput> => {
      const response = await getToolsIdGraph(slug, graphParams);

      return ToolGraphResponse.parse(response);
    },
    queryKey: toolKeys.graph(slug, graphParams),
    retry: false,
  });
};
