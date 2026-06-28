import { useQuery } from '@tanstack/react-query';
import { useEffect, useMemo, useRef, useState } from 'react';

import { toolGraphQueryOptions, toolSearchQueryOptions } from '@/api/tools-query-options';
import { Alert } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { computeNodePositions, computeViewBox } from './shared/tool-graph-layout';
import {
  DEFAULT_TOOL_GRAPH_SEARCH,
  type ToolGraphRelationshipKind,
  type ToolGraphRouteSearch,
} from './shared/tool-graph-types';
import { ToolGraphCanvas, ZOOM_STEP, zoomViewBox } from './tool-graph-home/tool-graph-canvas';
import { ToolGraphControls } from './tool-graph-home/tool-graph-controls';
import { ToolGraphSidePanel } from './tool-graph-home/tool-graph-side-panel';

type ToolGraphHomeContainerProps = {
  onSearchChange: (patch: Partial<ToolGraphRouteSearch>) => void;
  search: ToolGraphRouteSearch;
};

const useDebouncedValue = (value: string, delayMs: number): string => {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const timeout = window.setTimeout(() => {
      setDebouncedValue(value);
    }, delayMs);

    return () => window.clearTimeout(timeout);
  }, [delayMs, value]);

  return debouncedValue;
};

export const ToolGraphHomeContainer = ({ onSearchChange, search }: ToolGraphHomeContainerProps) => {
  const [selectedNodeId, setSelectedNodeId] = useState('');
  const [hoveredNodeId, setHoveredNodeId] = useState('');
  const [viewBox, setViewBox] = useState('0 0 760 520');
  const hasFitByTool = useRef<string>('');
  const debouncedQuery = useDebouncedValue(search.q.trim(), 250);

  const searchQuery = useQuery(toolSearchQueryOptions(debouncedQuery));
  const graphQuery = useQuery(
    toolGraphQueryOptions({
      depth: search.depth,
      kinds: search.kinds.length > 0 ? search.kinds : undefined,
      layoutMode: search.layoutMode,
      slug: search.tool,
    }),
  );

  const graph = graphQuery.data;
  const positionPoints = useMemo(() => {
    if (!graph) {
      return [];
    }

    return Array.from(computeNodePositions(graph, search.layoutMode).values());
  }, [graph, search.layoutMode]);

  useEffect(() => {
    if (!graph) {
      return;
    }

    if (!selectedNodeId || !graph.nodes.some((node) => node.id === selectedNodeId)) {
      setSelectedNodeId(graph.focusNodeId);
    }

    if (!hoveredNodeId || !graph.nodes.some((node) => node.id === hoveredNodeId)) {
      setHoveredNodeId(graph.focusNodeId);
    }

    if (hasFitByTool.current !== search.tool) {
      setViewBox(computeViewBox(positionPoints));
      hasFitByTool.current = search.tool;
    }
  }, [graph, hoveredNodeId, positionPoints, search.tool, selectedNodeId]);

  const handleKindToggle = (value: ToolGraphRelationshipKind) => {
    const nextKinds = search.kinds.includes(value)
      ? search.kinds.filter((kind) => kind !== value)
      : [...search.kinds, value];

    onSearchChange({ kinds: nextKinds });
  };

  const handleReset = () => {
    onSearchChange(DEFAULT_TOOL_GRAPH_SEARCH);
    setSelectedNodeId('');
    setHoveredNodeId('');
    setViewBox('0 0 760 520');
    hasFitByTool.current = '';
  };

  const handleViewBoxReset = () => {
    setViewBox(computeViewBox(positionPoints));
  };

  return (
    <main className="feature-page-shell">
      <header className="feature-page-header">
        <div>
          <p className="island-kicker mb-1">Instech atlas</p>
          <h1 className="text-2xl font-semibold">Homepage Graph Playground</h1>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="secondary">/{search.tool ? `?tool=${search.tool}` : ''}</Badge>
          <Button onClick={() => setViewBox(zoomViewBox(viewBox, ZOOM_STEP))} size="xs" variant="outline">
            Zoom in
          </Button>
          <Button onClick={() => setViewBox(zoomViewBox(viewBox, 1 / ZOOM_STEP))} size="xs" variant="outline">
            Zoom out
          </Button>
        </div>
      </header>

      <section className="feature-page-grid">
        <ToolGraphControls
          depth={search.depth}
          isSearching={searchQuery.isFetching}
          kinds={search.kinds}
          layoutMode={search.layoutMode}
          onDepthChange={(value) => onSearchChange({ depth: value })}
          onKindToggle={handleKindToggle}
          onLayoutModeChange={(value) => onSearchChange({ layoutMode: value })}
          onQueryChange={(value) => onSearchChange({ q: value })}
          onReset={handleReset}
          onToolSelect={(slug) => onSearchChange({ tool: slug })}
          q={search.q}
          selectedTool={search.tool}
          suggestions={searchQuery.data ?? []}
        />

        {search.tool ? null : (
          <section className="feature-panel feature-empty-state">
            <h2 className="mb-2 text-lg font-semibold">Pick a tool to start exploring</h2>
            <p className="text-muted-foreground text-sm">
              Search from the left panel, pick a tool, then inspect graph links by kind, depth, and layout.
            </p>
          </section>
        )}

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
              hoveredNodeId={hoveredNodeId}
              layoutMode={search.layoutMode}
              onNodeHover={setHoveredNodeId}
              onNodeSelect={setSelectedNodeId}
              onViewBoxReset={handleViewBoxReset}
              selectedNodeId={selectedNodeId}
              viewBox={viewBox}
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
        <span>Kinds: {search.kinds.length > 0 ? search.kinds.join(', ') : 'all'}</span>
      </footer>
    </main>
  );
};
