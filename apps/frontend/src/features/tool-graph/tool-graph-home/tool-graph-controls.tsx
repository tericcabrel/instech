import type { ToolSearchResultItemOutput } from '@/api/generated/model/ToolSearchResultItem.zod';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import {
  TOOL_GRAPH_DEPTHS,
  TOOL_GRAPH_LAYOUT_MODES,
  TOOL_GRAPH_RELATIONSHIP_KINDS,
  type ToolGraphDepth,
  type ToolGraphLayoutMode,
  type ToolGraphRelationshipKind,
} from '../shared/tool-graph-types';

type ToolGraphControlsProps = {
  depth: ToolGraphDepth;
  isSearching: boolean;
  kinds: ToolGraphRelationshipKind[];
  layoutMode: ToolGraphLayoutMode;
  onDepthChange: (value: ToolGraphDepth) => void;
  onKindToggle: (value: ToolGraphRelationshipKind) => void;
  onLayoutModeChange: (value: ToolGraphLayoutMode) => void;
  onQueryChange: (value: string) => void;
  onReset: () => void;
  onToolSelect: (slug: string) => void;
  q: string;
  selectedTool: string;
  suggestions: ToolSearchResultItemOutput[];
};

const toTitleCase = (value: string): string =>
  value
    .replaceAll('_', ' ')
    .split(' ')
    .map((part) => `${part[0]?.toUpperCase() ?? ''}${part.slice(1)}`)
    .join(' ');

export const ToolGraphControls = ({
  depth,
  isSearching,
  kinds,
  layoutMode,
  onDepthChange,
  onKindToggle,
  onLayoutModeChange,
  onQueryChange,
  onReset,
  onToolSelect,
  q,
  selectedTool,
  suggestions,
}: ToolGraphControlsProps) => (
  <div className="feature-panel feature-toolbar">
    <div className="flex items-center justify-between gap-3">
      <h2 className="text-sm font-semibold tracking-wide uppercase">Graph filters</h2>
      <Button onClick={onReset} size="xs" variant="outline">
        Reset
      </Button>
    </div>

    <div className="space-y-2">
      <Label className="text-xs font-medium uppercase">Search tools</Label>
      <Input onChange={(event) => onQueryChange(event.target.value)} placeholder="Type to find a tool..." value={q} />
      {isSearching ? <p className="text-muted-foreground text-xs">Searching…</p> : null}
      {q.trim() && suggestions.length > 0 ? (
        <ul className="feature-suggestion-list">
          {suggestions.map((item) => (
            <li key={item.id}>
              <Button
                className="w-full justify-between"
                onClick={() => onToolSelect(item.slug)}
                size="xs"
                variant={selectedTool === item.slug ? 'secondary' : 'ghost'}>
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
      <Label className="text-xs font-medium uppercase">Depth</Label>
      <div className="flex flex-wrap gap-2">
        {TOOL_GRAPH_DEPTHS.map((value) => (
          <Button
            key={value}
            onClick={() => onDepthChange(value)}
            size="xs"
            variant={depth === value ? 'default' : 'outline'}>
            Depth {value}
          </Button>
        ))}
      </div>
    </div>

    <div className="space-y-2">
      <Label className="text-xs font-medium uppercase">Layout</Label>
      <div className="flex flex-wrap gap-2">
        {TOOL_GRAPH_LAYOUT_MODES.map((value) => (
          <Button
            key={value}
            onClick={() => onLayoutModeChange(value)}
            size="xs"
            variant={layoutMode === value ? 'default' : 'outline'}>
            {toTitleCase(value)}
          </Button>
        ))}
      </div>
    </div>

    <div className="space-y-2">
      <Label className="text-xs font-medium uppercase">Kinds</Label>
      <div className="flex flex-wrap gap-1.5">
        {TOOL_GRAPH_RELATIONSHIP_KINDS.map((kind) => {
          const active = kinds.includes(kind);

          return (
            <button className="rounded-none" key={kind} onClick={() => onKindToggle(kind)} type="button">
              <Badge variant={active ? 'default' : 'outline'}>{toTitleCase(kind)}</Badge>
            </button>
          );
        })}
      </div>
    </div>
  </div>
);
