import { queryOptions } from '@tanstack/react-query'
import { getRelationshipsQuery } from './generated/instech'
import { RelationshipListResponse as RelationshipListResponseSchema } from './generated/model'
import type {
  GetRelationshipsQueryParams,
  RelationshipListResponseOutput,
} from './generated/model'

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
    queryFn: async (): Promise<RelationshipListResponseOutput> => {
      const response = await getRelationshipsQuery(params)

      return RelationshipListResponseSchema.parse(response)
    },
    retry: false,
  })
