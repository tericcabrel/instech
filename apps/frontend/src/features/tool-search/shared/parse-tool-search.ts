const isString = (value: unknown): value is string => typeof value === 'string';

export type ToolSearchRouteSearch = {
  q: string;
};

export const parseToolSearch = (search: Record<string, unknown>): ToolSearchRouteSearch => ({
  q: isString(search.q) ? search.q.trim() : '',
});
