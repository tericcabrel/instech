import type { ToolSearchResultItemOutput } from '@/api/generated/model/ToolSearchResultItem.zod'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'

import {
  TOOL_GRAPH_DEPTHS,
  TOOL_GRAPH_LAYOUT_MODES,
  TOOL_GRAPH_RELATIONSHIP_KINDS,
  type ToolGraphDepth,
  type ToolGraphLayoutMode,
  type ToolGraphRelationshipKind,
} from '../shared/tool-graph-types'

type ToolGraphControlsProps = {
  q: string
  selectedTool: string
  depth: ToolGraphDepth
  layoutMode: ToolGraphLayoutMode
  kinds: ToolGraphRelationshipKind[]
  isSearching: boolean
  suggestions: ToolSearchResultItemOutput[]
  onQueryChange: (value: string) => void
  onToolSelect: (slug: string) => void
  onDepthChange: (value: ToolGraphDepth) => void
  onLayoutModeChange: (value: ToolGraphLayoutMode) => void
  onKindToggle: (value: ToolGraphRelationshipKind) => void
  onReset: () => void
}

const toTitleCase = (value: string): string =>
  value
    .replaceAll('_', ' ')
    .split(' ')
    .map((part) => `${part[0]?.toUpperCase() ?? ''}${part.slice(1)}`)
    .join(' ')

export function ToolGraphControls({
  q,
  selectedTool,
  depth,
  layoutMode,
  kinds,
  isSearching,
  suggestions,
  onQueryChange,
  onToolSelect,
  onDepthChange,
  onLayoutModeChange,
  onKindToggle,
  onReset,
}: ToolGraphControlsProps) {
  return (
    <div className="feature-panel feature-toolbar">
      <div className="flex items-center justify-between gap-3">
        <h2 className="text-sm font-semibold tracking-wide uppercase">Graph filters</h2>
        <Button size="xs" variant="outline" onClick={onReset}>
          Reset
        </Button>
      </div>

      <div className="space-y-2">
        <label className="text-xs font-medium uppercase">Search tools</label>
        <Input
          placeholder="Type to find a tool..."
          value={q}
          onChange={(event) => onQueryChange(event.target.value)}
        />
        {isSearching ? (
          <p className="text-muted-foreground text-xs">Searching…</p>
        ) : null}
        {q.trim() && suggestions.length > 0 ? (
          <ul className="feature-suggestion-list">
            {suggestions.map((item) => (
              <li key={item.id}>
                <Button
                  size="xs"
                  variant={selectedTool === item.slug ? 'secondary' : 'ghost'}
                  className="w-full justify-between"
                  onClick={() => onToolSelect(item.slug)}
                >
                  <span className="truncate">{item.name}</span>
                  <span className="text-muted-foreground ml-2">{item.slug}</span>
                </Button>
              </li>
            ))}
          </ul>
        ) : null}
      </div>

      <Separator />

      <div className="space-y-2">
        <label className="text-xs font-medium uppercase">Depth</label>
        <div className="flex flex-wrap gap-2">
          {TOOL_GRAPH_DEPTHS.map((value) => (
            <Button
              key={value}
              size="xs"
              variant={depth === value ? 'default' : 'outline'}
              onClick={() => onDepthChange(value)}
            >
              Depth {value}
            </Button>
          ))}
        </div>
      </div>

      <div className="space-y-2">
        <label className="text-xs font-medium uppercase">Layout</label>
        <div className="flex flex-wrap gap-2">
          {TOOL_GRAPH_LAYOUT_MODES.map((value) => (
            <Button
              key={value}
              size="xs"
              variant={layoutMode === value ? 'default' : 'outline'}
              onClick={() => onLayoutModeChange(value)}
            >
              {toTitleCase(value)}
            </Button>
          ))}
        </div>
      </div>

      <div className="space-y-2">
        <label className="text-xs font-medium uppercase">Kinds</label>
        <div className="flex flex-wrap gap-1.5">
          {TOOL_GRAPH_RELATIONSHIP_KINDS.map((kind) => {
            const active = kinds.includes(kind)
            return (
              <button
                key={kind}
                type="button"
                className="rounded-none"
                onClick={() => onKindToggle(kind)}
              >
                <Badge variant={active ? 'default' : 'outline'}>{toTitleCase(kind)}</Badge>
              </button>
            )
          })}
        </div>
      </div>
    </div>
  )
}
