import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { postTools } from "./generated/instech"
import type { CreateToolRequest } from "./generated/model"
import { toolKeys, toolsQueryOptions } from "./tools-query-options"

export const useCreateTool = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (body: CreateToolRequest) => postTools(body),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: toolKeys.all }),
  })
}

export const useGetTools = () => {
  return useQuery(toolsQueryOptions('javascript'))
}

export const useSearchTools = (keyword: string) => {
  return useQuery(toolsQueryOptions(keyword))
}