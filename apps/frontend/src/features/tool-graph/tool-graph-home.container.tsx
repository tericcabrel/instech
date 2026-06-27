import { useQuery } from '@tanstack/react-query'
import { useEffect, useMemo, useRef, useState } from 'react'

import { toolGraphQueryOptions, toolSearchQueryOptions } from '@/api/tools-query-options'
import { Alert } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

import { ToolGraphCanvas, ZOOM_STEP, zoomViewBox } from './tool-graph-home/tool-graph-canvas'
import { ToolGraphControls } from './tool-graph-home/tool-graph-controls'
import { ToolGraphSidePanel } from './tool-graph-home/tool-graph-side-panel'
import { computeNodePositions, computeViewBox } from './shared/tool-graph-layout'
import {
  DEFAULT_TOOL_GRAPH_SEARCH,
  type ToolGraphRouteSearch,
  type ToolGraphRelationshipKind,
} from './shared/tool-graph-types'

type ToolGraphHomeContainerProps = {
  search: ToolGraphRouteSearch
  onSearchChange: (patch: Partial<ToolGraphRouteSearch>) => void
}

const useDebouncedValue = (value: string, delayMs: number): string => {
  const [debouncedValue, setDebouncedValue] = useState(value)

  useEffect(() => {
    const timeout = window.setTimeout(() => {
      setDebouncedValue(value)
    }, delayMs)

    return () => window.clearTimeout(timeout)
  }, [delayMs, value])

  return debouncedValue
}

export function ToolGraphHomeContainer({ search, onSearchChange }: ToolGraphHomeContainerProps) {
  const [selectedNodeId, setSelectedNodeId] = useState('')
  const [hoveredNodeId, setHoveredNodeId] = useState('')
  const [viewBox, setViewBox] = useState('0 0 760 520')
  const hasFitByTool = useRef<string>('')
  const debouncedQuery = useDebouncedValue(search.q.trim(), 250)

  const searchQuery = useQuery(toolSearchQueryOptions(debouncedQuery))
  const graphQuery = useQuery(
    toolGraphQueryOptions({
      slug: search.tool,
      depth: search.depth,
      kinds: search.kinds.length > 0 ? search.kinds : undefined,
      layoutMode: search.layoutMode,
    }),
  )

  const graph = graphQuery.data
  const positionPoints = useMemo(() => {
    if (!graph) {
      return []
    }

    return Array.from(computeNodePositions(graph, search.layoutMode).values())
  }, [graph, search.layoutMode])

  useEffect(() => {
    if (!graph) {
      return
    }

    if (!selectedNodeId || !graph.nodes.some((node) => node.id === selectedNodeId)) {
      setSelectedNodeId(graph.focusNodeId)
    }

    if (!hoveredNodeId || !graph.nodes.some((node) => node.id === hoveredNodeId)) {
      setHoveredNodeId(graph.focusNodeId)
    }

    if (hasFitByTool.current !== search.tool) {
      setViewBox(computeViewBox(positionPoints))
      hasFitByTool.current = search.tool
    }
  }, [graph, hoveredNodeId, positionPoints, search.tool, selectedNodeId])

  const handleKindToggle = (value: ToolGraphRelationshipKind) => {
    const nextKinds = search.kinds.includes(value)
      ? search.kinds.filter((kind) => kind !== value)
      : [...search.kinds, value]

    onSearchChange({ kinds: nextKinds })
  }

  const handleReset = () => {
    onSearchChange(DEFAULT_TOOL_GRAPH_SEARCH)
    setSelectedNodeId('')
    setHoveredNodeId('')
    setViewBox('0 0 760 520')
    hasFitByTool.current = ''
  }

  const handleViewBoxReset = () => {
    setViewBox(computeViewBox(positionPoints))
  }

  return (
    <main className="feature-page-shell">
      <header className="feature-page-header">
        <div>
          <p className="island-kicker mb-1">Instech atlas</p>
          <h1 className="text-2xl font-semibold">Homepage Graph Playground</h1>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="secondary">/{search.tool ? `?tool=${search.tool}` : ''}</Badge>
          <Button size="xs" variant="outline" onClick={() => setViewBox(zoomViewBox(viewBox, ZOOM_STEP))}>
            Zoom in
          </Button>
          <Button
            size="xs"
            variant="outline"
            onClick={() => setViewBox(zoomViewBox(viewBox, 1 / ZOOM_STEP))}
          >
            Zoom out
          </Button>
        </div>
      </header>

      <section className="feature-page-grid">
        <ToolGraphControls
          q={search.q}
          selectedTool={search.tool}
          depth={search.depth}
          layoutMode={search.layoutMode}
          kinds={search.kinds}
          isSearching={searchQuery.isFetching}
          suggestions={searchQuery.data ?? []}
          onQueryChange={(value) => onSearchChange({ q: value })}
          onToolSelect={(slug) => onSearchChange({ tool: slug })}
          onDepthChange={(value) => onSearchChange({ depth: value })}
          onLayoutModeChange={(value) => onSearchChange({ layoutMode: value })}
          onKindToggle={handleKindToggle}
          onReset={handleReset}
        />

        {!search.tool ? (
          <section className="feature-panel feature-empty-state">
            <h2 className="mb-2 text-lg font-semibold">Pick a tool to start exploring</h2>
            <p className="text-muted-foreground text-sm">
              Search from the left panel, pick a tool, then inspect graph links by kind, depth, and layout.
            </p>
          </section>
        ) : null}

        {search.tool && graphQuery.isLoading ? (
          <section className="feature-panel feature-empty-state">
            <p className="text-muted-foreground text-sm">Loading graph…</p>
          </section>
        ) : null}

        {search.tool && graphQuery.isError ? (
          <section className="feature-panel feature-empty-state">
            <Alert variant="destructive">
              Failed to load graph for <strong>{search.tool}</strong>. Try selecting another tool.
            </Alert>
          </section>
        ) : null}

        {search.tool && graph ? (
          <>
            <ToolGraphCanvas
              graph={graph}
              layoutMode={search.layoutMode}
              selectedNodeId={selectedNodeId}
              hoveredNodeId={hoveredNodeId}
              viewBox={viewBox}
              onViewBoxReset={handleViewBoxReset}
              onNodeHover={setHoveredNodeId}
              onNodeSelect={setSelectedNodeId}
            />
            <ToolGraphSidePanel graph={graph} selectedNodeId={selectedNodeId} />
          </>
        ) : null}
      </section>

      <footer className="feature-status-strip">
        <span>Query: {search.q.trim() || 'none'}</span>
        <span>Tool: {search.tool || 'none'}</span>
        <span>Depth: {search.depth}</span>
        <span>Layout: {search.layoutMode}</span>
        <span>Kinds: {search.kinds.length ? search.kinds.join(', ') : 'all'}</span>
      </footer>
    </main>
  )
}
