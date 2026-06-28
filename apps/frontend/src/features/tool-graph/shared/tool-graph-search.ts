import {
  DEFAULT_TOOL_GRAPH_SEARCH,
  TOOL_GRAPH_DEPTHS,
  TOOL_GRAPH_LAYOUT_MODES,
  TOOL_GRAPH_RELATIONSHIP_KINDS,
  type ToolGraphDepth,
  type ToolGraphLayoutMode,
  type ToolGraphRelationshipKind,
  type ToolGraphRouteSearch,
} from './tool-graph-types';

const isString = (value: unknown): value is string => typeof value === 'string';

const parseKinds = (value: unknown): ToolGraphRelationshipKind[] => {
  const cleanKinds = (kinds: ToolGraphRelationshipKind[]): ToolGraphRelationshipKind[] => {
    const unique = new Set(kinds);

    return TOOL_GRAPH_RELATIONSHIP_KINDS.filter((kind) => unique.has(kind));
  };

  if (Array.isArray(value)) {
    return cleanKinds(
      value.filter((item): item is ToolGraphRelationshipKind =>
        TOOL_GRAPH_RELATIONSHIP_KINDS.includes(item as ToolGraphRelationshipKind),
      ),
    );
  }

  if (isString(value)) {
    return cleanKinds(
      value
        .split(',')
        .map((part) => part.trim())
        .filter((part): part is ToolGraphRelationshipKind =>
          TOOL_GRAPH_RELATIONSHIP_KINDS.includes(part as ToolGraphRelationshipKind),
        ),
    );
  }

  return DEFAULT_TOOL_GRAPH_SEARCH.kinds;
};

const parseDepth = (value: unknown): ToolGraphDepth => {
  const parsed = Number(value);

  return TOOL_GRAPH_DEPTHS.includes(parsed as ToolGraphDepth)
    ? (parsed as ToolGraphDepth)
    : DEFAULT_TOOL_GRAPH_SEARCH.depth;
};

const parseLayoutMode = (value: unknown): ToolGraphLayoutMode => {
  if (isString(value) && TOOL_GRAPH_LAYOUT_MODES.includes(value as ToolGraphLayoutMode)) {
    return value as ToolGraphLayoutMode;
  }

  return DEFAULT_TOOL_GRAPH_SEARCH.layoutMode;
};

export const parseToolGraphSearch = (search: Record<string, unknown>): ToolGraphRouteSearch => ({
  depth: parseDepth(search.depth),
  kinds: parseKinds(search.kinds),
  layoutMode: parseLayoutMode(search.layoutMode),
  q: isString(search.q) ? search.q : DEFAULT_TOOL_GRAPH_SEARCH.q,
  tool: isString(search.tool) ? search.tool : DEFAULT_TOOL_GRAPH_SEARCH.tool,
});
