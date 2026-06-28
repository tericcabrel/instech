import { createFileRoute } from '@tanstack/react-router';

import { toolSearchQueryOptions } from '@/api/tools-query-options';
import { parseToolSearch, type ToolSearchRouteSearch } from '@/features/tool-search/shared/parse-tool-search';
import { ToolSearchResultsContainer } from '@/features/tool-search/tool-search-results.container';

const SearchRouteComponent = () => {
  const search = Route.useSearch();
  const navigate = Route.useNavigate();

  return (
    <ToolSearchResultsContainer
      onQueryChange={(q) =>
        navigate({
          replace: true,
          search: (): ToolSearchRouteSearch => ({ q }),
        })
      }
      q={search.q}
    />
  );
};

export const Route = createFileRoute('/search')({
  component: SearchRouteComponent,
  loader: async ({ context, deps }) => {
    const query = deps.q;

    // Orval client calls expect browser globals; skip prefetch during SSR.
    if (!query || typeof window === 'undefined') {
      return;
    }

    await context.queryClient.ensureQueryData(toolSearchQueryOptions(query));
  },
  loaderDeps: ({ search }) => ({ q: search.q.trim() }),
  validateSearch: parseToolSearch,
});
