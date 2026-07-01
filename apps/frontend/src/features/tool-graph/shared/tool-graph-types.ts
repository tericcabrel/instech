export const TOOL_GRAPH_RELATIONSHIP_KINDS = [
  'built_on',
  'inspired_by',
  'alternative_to',
  'replaced_by',
  'used_with',
] as const;

export const TOOL_GRAPH_LAYOUT_MODES = ['chronological', 'force'] as const;

export const TOOL_GRAPH_DEPTHS = [1, 2] as const;

export const DEFAULT_TOOL_GRAPH_SEARCH = {
  depth: 1 as const,
  kinds: [] as ToolGraphRelationshipKind[],
  layoutMode: 'force' as const,
  q: '',
  tool: '',
};

export type ToolGraphRelationshipKind = (typeof TOOL_GRAPH_RELATIONSHIP_KINDS)[number];

export type ToolGraphLayoutMode = (typeof TOOL_GRAPH_LAYOUT_MODES)[number];

export type ToolGraphDepth = (typeof TOOL_GRAPH_DEPTHS)[number];

export type ToolGraphRouteSearch = {
  depth: ToolGraphDepth;
  kinds: ToolGraphRelationshipKind[];
  layoutMode: ToolGraphLayoutMode;
  q: string;
  tool: string;
};
