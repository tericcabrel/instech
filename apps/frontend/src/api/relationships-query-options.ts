import { queryOptions } from '@tanstack/react-query'
import { getRelationshipsQuery } from './generated/instech'
import type { GetRelationshipsQueryParams, RelationshipListResponse } from './generated/model'

const normalizeRelationshipParams = (params?: GetRelationshipsQueryParams) => ({
  toolId: params?.toolId,
  kind: params?.kind,
  cursor: params?.cursor,
  limit: params?.limit,
})

export const relationshipKeys = {
  all: ['relationships'] as const,
  query: (params?: GetRelationshipsQueryParams) =>
    [...relationshipKeys.all, 'query', normalizeRelationshipParams(params)] as const,
}

export const relationshipsQueryOptions = (params?: GetRelationshipsQueryParams) =>
  queryOptions({
    queryKey: relationshipKeys.query(params),
    queryFn: async (): Promise<RelationshipListResponse> =>
      (await getRelationshipsQuery(params)) as unknown as RelationshipListResponse,
    retry: false,
  })
