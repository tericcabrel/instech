import { createFileRoute } from '@tanstack/react-router'
import { useGetTools } from '@/api/tools'

export const Route = createFileRoute('/demo/tanstack-query')({
  component: TanStackQueryDemo,
})

function TanStackQueryDemo() {
  const { data, error, isError, isLoading } = useGetTools()
  const tools = data ?? []

  return (
    <main className="demo-page demo-center">
      <section className="demo-panel w-full max-w-2xl">
        <p className="island-kicker mb-2">TanStack Query</p>
        <h1 className="demo-title mb-6">Available tools</h1>

        {isLoading ? (
          <p className="text-sm text-neutral-500">Loading tools...</p>
        ) : null}

        {isError ? (
          <p className="text-sm text-red-600">
            Failed to load tools: {error.message}
          </p>
        ) : null}

        {!isLoading && !isError && tools.length === 0 ? (
          <p className="text-sm text-neutral-500">No tools available.</p>
        ) : null}

        {!isLoading && !isError && tools.length > 0 ? (
          <ul className="mb-4 space-y-2">
            {tools.map((tool) => (
              <li key={tool.id} className="demo-list-item">
                <span className="text-base font-medium">{tool.name}</span>
                <span className="ml-2 text-sm text-neutral-500">
                  ({tool.id})
                </span>
              </li>
            ))}
          </ul>
        ) : null}
      </section>
    </main>
  )
}
