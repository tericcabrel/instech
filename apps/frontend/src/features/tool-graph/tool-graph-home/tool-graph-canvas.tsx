import type { ToolGraphResponseOutput } from '@/api/generated/model/ToolGraphResponse.zod';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';

import { computeNodePositions, kindColor, kindLabel } from '../shared/tool-graph-layout';
import type { ToolGraphLayoutMode } from '../shared/tool-graph-types';

type ToolGraphCanvasProps = {
  graph: ToolGraphResponseOutput;
  hoveredNodeId: string;
  layoutMode: ToolGraphLayoutMode;
  onNodeHover: (nodeId: string) => void;
  onNodeSelect: (nodeId: string) => void;
  onViewBoxReset: () => void;
  selectedNodeId: string;
  viewBox: string;
};

export const ZOOM_STEP = 0.88;

export const zoomViewBox = (viewBox: string, factor: number): string => {
  const [x, y, width, height] = viewBox.split(' ').map(Number);
  const nextWidth = width * factor;
  const nextHeight = height * factor;
  const nextX = x + (width - nextWidth) / 2;
  const nextY = y + (height - nextHeight) / 2;

  return `${nextX} ${nextY} ${nextWidth} ${nextHeight}`;
};

export const ToolGraphCanvas = ({
  graph,
  hoveredNodeId,
  layoutMode,
  onNodeHover,
  onNodeSelect,
  onViewBoxReset,
  selectedNodeId,
  viewBox,
}: ToolGraphCanvasProps) => {
  const positions = computeNodePositions(graph, layoutMode);

  return (
    <section className="feature-panel feature-graph-shell">
      <header className="mb-2 flex items-center justify-between gap-3">
        <div>
          <h2 className="text-base font-semibold">Relationship graph</h2>
          <p className="text-muted-foreground text-xs">Click a node to pin details in side panel.</p>
        </div>
        <div className="flex items-center gap-1">
          <Button onClick={onViewBoxReset} size="xs" variant="outline">
            Fit
          </Button>
          <Button onClick={onViewBoxReset} size="xs" variant="outline">
            Reset camera
          </Button>
        </div>
      </header>

      <div className="feature-graph-canvas">
        <svg aria-label="Tool relationship graph" role="img" viewBox={viewBox}>
          <g>
            {graph.links.map((link) => {
              const source = positions.get(link.source);
              const target = positions.get(link.target);

              if (!source || !target) {
                return null;
              }

              return (
                <line
                  key={link.id}
                  stroke={kindColor(link.kind)}
                  strokeOpacity={0.72}
                  strokeWidth={2.2}
                  x1={source.x}
                  x2={target.x}
                  y1={source.y}
                  y2={target.y}
                />
              );
            })}
          </g>

          <g>
            {graph.nodes.map((node) => {
              const point = positions.get(node.id);

              if (!point) {
                return null;
              }

              const isSelected = selectedNodeId === node.id;
              const isHovered = hoveredNodeId === node.id;
              const valueWhenIsHovered = isHovered ? 19 : 17;
              const radius = isSelected ? 21 : valueWhenIsHovered;

              return (
                <g
                  key={node.id}
                  onClick={() => onNodeSelect(node.id)}
                  onFocus={() => onNodeHover(node.id)}
                  onMouseEnter={() => onNodeHover(node.id)}
                  role="button">
                  <circle
                    cx={point.x}
                    cy={point.y}
                    fill={node.isFocus ? '#4fb8b2' : '#e7f0e8'}
                    r={radius}
                    stroke={isSelected ? '#173a40' : '#4fb8b2'}
                    strokeWidth={isSelected ? 3.5 : 2}
                  />
                  <text
                    fill="var(--sea-ink)"
                    fontSize={13}
                    fontWeight={isSelected ? 700 : 500}
                    textAnchor="middle"
                    x={point.x}
                    y={point.y - 28}>
                    {node.name}
                  </text>
                </g>
              );
            })}
          </g>
        </svg>
      </div>

      <footer className="mt-3 flex flex-wrap gap-1.5">
        {graph.meta.kindsApplied.map((kind) => (
          <Badge key={kind} style={{ borderColor: kindColor(kind) }} variant="outline">
            {kindLabel(kind)}
          </Badge>
        ))}
      </footer>
    </section>
  );
};
