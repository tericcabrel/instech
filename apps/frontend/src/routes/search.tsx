import { createFileRoute } from '@tanstack/react-router'

import { toolSearchQueryOptions } from '@/api/tools-query-options'
import {
  parseToolSearch,
  type ToolSearchRouteSearch,
} from '@/features/tool-search/shared/parse-tool-search'
import { ToolSearchResultsContainer } from '@/features/tool-search/tool-search-results.container'

export const Route = createFileRoute('/search')({
  validateSearch: parseToolSearch,
  loaderDeps: ({ search }) => ({ q: search.q.trim() }),
  loader: async ({ context, deps }) => {
    const query = deps.q

    // Orval client calls expect browser globals; skip prefetch during SSR.
    if (!query || typeof window === 'undefined') {
      return
    }

    await context.queryClient.ensureQueryData(toolSearchQueryOptions(query))
  },
  component: SearchRouteComponent,
})

function SearchRouteComponent() {
  const search = Route.useSearch()
  const navigate = Route.useNavigate()

  return (
    <ToolSearchResultsContainer
      q={search.q}
      onQueryChange={(q) =>
        navigate({
          search: (): ToolSearchRouteSearch => ({ q }),
          replace: true,
        })
      }
    />
  )
}
