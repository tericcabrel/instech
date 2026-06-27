import { createFileRoute } from '@tanstack/react-router'

import { ToolGraphHomeContainer } from '@/features/tool-graph/tool-graph-home.container'
import { parseToolGraphSearch } from '@/features/tool-graph/shared/tool-graph-search'

export const Route = createFileRoute('/')({
  validateSearch: parseToolGraphSearch,
  component: HomeRouteComponent,
})

function HomeRouteComponent() {
  const search = Route.useSearch()
  const navigate = Route.useNavigate()

  return (
    <ToolGraphHomeContainer
      search={search}
      onSearchChange={(patch) =>
        navigate({
          search: (previous) => ({
            ...previous,
            ...patch,
          }),
          replace: true,
        })
      }
    />
  )
}
