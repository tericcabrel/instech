import { queryOptions } from '@tanstack/react-query'
import { getRelationshipsQuery } from './generated/instech'
import type { GetRelationshipsQueryParams } from './generated/model'
import {
  RelationshipListResponse,
  type RelationshipListResponseOutput,
} from './generated/model/relationshipListResponse.zod.ts'

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

      return RelationshipListResponse.parse(response)
    },
    retry: false,
  })
