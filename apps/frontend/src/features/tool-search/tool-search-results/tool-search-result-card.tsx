import { Link } from '@tanstack/react-router';

import type { ToolSearchResultItemOutput } from '@/api/generated/model/ToolSearchResultItem.zod';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { DEFAULT_TOOL_GRAPH_SEARCH } from '@/features/tool-graph/shared/tool-graph-types';

type ToolSearchResultCardProps = {
  item: ToolSearchResultItemOutput;
};

const toLabel = (value: string): string => `${value.charAt(0).toUpperCase()}${value.slice(1).replaceAll('_', ' ')}`;

export const ToolSearchResultCard = ({ item }: ToolSearchResultCardProps) => (
  <Card className="feature-panel p-0">
    <CardHeader className="space-y-2">
      <CardTitle>{item.name}</CardTitle>
      <p className="text-muted-foreground text-xs">{item.slug}</p>
      <div className="flex flex-wrap gap-1.5">
        <Badge variant="secondary">{toLabel(item.category)}</Badge>
        <Badge variant="outline">{toLabel(item.subType)}</Badge>
      </div>
    </CardHeader>
    <CardContent className="pb-2" />
    <CardFooter className="flex flex-wrap justify-end gap-2">
      <Button asChild size="xs" variant="outline">
        <Link search={{ ...DEFAULT_TOOL_GRAPH_SEARCH, tool: item.slug }} to="/">
          View graph
        </Link>
      </Button>
      <Button asChild size="xs">
        <Link params={{ slug: item.slug }} to="/tools/$slug">
          View details
        </Link>
      </Button>
    </CardFooter>
  </Card>
);
