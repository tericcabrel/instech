import type { ToolGraphResponseOutput } from '@/api/generated/model/ToolGraphResponse.zod'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

import {
  computeNodePositions,
  kindColor,
  kindLabel,
} from '../shared/tool-graph-layout'
import type { ToolGraphLayoutMode } from '../shared/tool-graph-types'

type ToolGraphCanvasProps = {
  graph: ToolGraphResponseOutput
  layoutMode: ToolGraphLayoutMode
  selectedNodeId: string
  hoveredNodeId: string
  viewBox: string
  onViewBoxReset: () => void
  onNodeHover: (nodeId: string) => void
  onNodeSelect: (nodeId: string) => void
}

const ZOOM_STEP = 0.88

function zoomViewBox(viewBox: string, factor: number): string {
  const [x, y, width, height] = viewBox.split(' ').map(Number)
  const nextWidth = width * factor
  const nextHeight = height * factor
  const nextX = x + (width - nextWidth) / 2
  const nextY = y + (height - nextHeight) / 2

  return `${nextX} ${nextY} ${nextWidth} ${nextHeight}`
}

export function ToolGraphCanvas({
  graph,
  layoutMode,
  selectedNodeId,
  hoveredNodeId,
  viewBox,
  onViewBoxReset,
  onNodeHover,
  onNodeSelect,
}: ToolGraphCanvasProps) {
  const positions = computeNodePositions(graph, layoutMode)

  return (
    <section className="feature-panel feature-graph-shell">
      <header className="mb-2 flex items-center justify-between gap-3">
        <div>
          <h2 className="text-base font-semibold">Relationship graph</h2>
          <p className="text-muted-foreground text-xs">Click a node to pin details in side panel.</p>
        </div>
        <div className="flex items-center gap-1">
          <Button size="xs" variant="outline" onClick={onViewBoxReset}>
            Fit
          </Button>
          <Button size="xs" variant="outline" onClick={onViewBoxReset}>
            Reset camera
          </Button>
        </div>
      </header>

      <div className="feature-graph-canvas">
        <svg viewBox={viewBox} role="img" aria-label="Tool relationship graph">
          <g>
            {graph.links.map((link) => {
              const source = positions.get(link.source)
              const target = positions.get(link.target)

              if (!source || !target) {
                return null
              }

              return (
                <line
                  key={link.id}
                  x1={source.x}
                  y1={source.y}
                  x2={target.x}
                  y2={target.y}
                  stroke={kindColor(link.kind)}
                  strokeWidth={2.2}
                  strokeOpacity={0.72}
                />
              )
            })}
          </g>

          <g>
            {graph.nodes.map((node) => {
              const point = positions.get(node.id)

              if (!point) {
                return null
              }

              const isSelected = selectedNodeId === node.id
              const isHovered = hoveredNodeId === node.id
              const radius = isSelected ? 21 : isHovered ? 19 : 17

              return (
                <g
                  key={node.id}
                  onMouseEnter={() => onNodeHover(node.id)}
                  onFocus={() => onNodeHover(node.id)}
                  onClick={() => onNodeSelect(node.id)}
                >
                  <circle
                    cx={point.x}
                    cy={point.y}
                    r={radius}
                    fill={node.isFocus ? '#4fb8b2' : '#e7f0e8'}
                    stroke={isSelected ? '#173a40' : '#4fb8b2'}
                    strokeWidth={isSelected ? 3.5 : 2}
                  />
                  <text
                    x={point.x}
                    y={point.y - 28}
                    textAnchor="middle"
                    fill="var(--sea-ink)"
                    fontSize={13}
                    fontWeight={isSelected ? 700 : 500}
                  >
                    {node.name}
                  </text>
                </g>
              )
            })}
          </g>
        </svg>
      </div>

      <footer className="mt-3 flex flex-wrap gap-1.5">
        {graph.meta.kindsApplied.map((kind) => (
          <Badge key={kind} variant="outline" style={{ borderColor: kindColor(kind) }}>
            {kindLabel(kind)}
          </Badge>
        ))}
      </footer>
    </section>
  )
}

export { zoomViewBox, ZOOM_STEP }
