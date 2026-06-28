import { createFileRoute } from '@tanstack/react-router';
import { parseToolGraphSearch } from '@/features/tool-graph/shared/tool-graph-search';
import { ToolGraphHomeContainer } from '@/features/tool-graph/tool-graph-home.container';

const HomeRouteComponent = () => {
  const search = Route.useSearch();
  const navigate = Route.useNavigate();

  return (
    <ToolGraphHomeContainer
      onSearchChange={(patch) =>
        navigate({
          replace: true,
          search: (previous) => ({
            ...previous,
            ...patch,
          }),
        })
      }
      search={search}
    />
  );
};

export const Route = createFileRoute('/')({
  component: HomeRouteComponent,
  validateSearch: parseToolGraphSearch,
});
