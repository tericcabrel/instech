import { createFileRoute } from '@tanstack/react-router';

import { toolAlternativesQueryOptions, toolQueryOptions } from '@/api/tools-query-options';
import { ToolAlternativesContainer } from '@/features/tool-alternatives/tool-alternatives.container';

const ToolAlternativesRouteComponent = () => {
  const { slug } = Route.useParams();

  return <ToolAlternativesContainer slug={slug} />;
};

export const Route = createFileRoute('/alternatives/$slug')({
  component: ToolAlternativesRouteComponent,
  loader: async ({ context, params }) => {
    const slug = params.slug.trim();

    // Orval client calls expect browser globals; skip prefetch during SSR.
    if (!slug || typeof window === 'undefined') {
      return;
    }

    await Promise.all([
      context.queryClient.ensureQueryData(toolAlternativesQueryOptions(slug)),
      context.queryClient.ensureQueryData(toolQueryOptions(slug)),
    ]);
  },
});
