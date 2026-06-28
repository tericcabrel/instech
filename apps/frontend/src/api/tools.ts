import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { postTools } from './generated/instech';
import type { CreateToolRequest } from './generated/model';
import { toolKeys, toolSearchQueryOptions } from './tools-query-options';

export const useCreateTool = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: CreateToolRequest) => postTools(body),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: toolKeys.all }),
  });
};

export const useGetTools = () => useQuery(toolSearchQueryOptions('javascript'));

export const useSearchTools = (keyword: string) => useQuery(toolSearchQueryOptions(keyword));
