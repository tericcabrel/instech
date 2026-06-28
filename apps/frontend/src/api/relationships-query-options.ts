import { queryOptions } from '@tanstack/react-query';
import { getRelationshipsQuery } from './generated/instech';
import type { GetRelationshipsQueryParams } from './generated/model';
import {
  RelationshipListResponse,
  type RelationshipListResponseOutput,
} from './generated/model/RelationshipListResponse.zod.ts';

const normalizeRelationshipParams = (params?: GetRelationshipsQueryParams) => ({
  cursor: params?.cursor,
  kind: params?.kind,
  limit: params?.limit,
  toolId: params?.toolId,
});

export const relationshipKeys = {
  all: ['relationships'] as const,
  query: (params?: GetRelationshipsQueryParams) =>
    [...relationshipKeys.all, 'query', normalizeRelationshipParams(params)] as const,
};

export const relationshipsQueryOptions = (params?: GetRelationshipsQueryParams) =>
  queryOptions({
    queryFn: async (): Promise<RelationshipListResponseOutput> => {
      const response = await getRelationshipsQuery(params);

      return RelationshipListResponse.parse(response);
    },
    queryKey: relationshipKeys.query(params),
    retry: false,
  });
