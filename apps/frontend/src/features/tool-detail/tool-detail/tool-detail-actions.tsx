import { Link } from '@tanstack/react-router';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { DEFAULT_TOOL_GRAPH_SEARCH } from '@/features/tool-graph/shared/tool-graph-types';

type ToolDetailActionsProps = {
  slug: string;
};

export const ToolDetailActions = ({ slug }: ToolDetailActionsProps) => (
  <Card className="feature-panel p-0">
    <CardHeader>
      <CardTitle>Explore next</CardTitle>
      <CardDescription>Open this tool in the graph view or inspect alternative tools.</CardDescription>
    </CardHeader>
    <CardContent />
    <CardFooter className="flex flex-wrap justify-end gap-2">
      <Button asChild size="xs" variant="outline">
        <Link search={{ ...DEFAULT_TOOL_GRAPH_SEARCH, tool: slug }} to="/">
          View graph
        </Link>
      </Button>
      <Button asChild size="xs">
        <Link params={{ slug }} to="/alternatives/$slug">
          View alternatives
        </Link>
      </Button>
    </CardFooter>
  </Card>
);
