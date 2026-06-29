import { createFileRoute } from '@tanstack/react-router';

import { toolQueryOptions } from '@/api/tools-query-options';
import { ToolDetailContainer } from '@/features/tool-detail/tool-detail.container';

const ToolDetailRouteComponent = () => {
  const { slug } = Route.useParams();

  return <ToolDetailContainer slug={slug} />;
};

export const Route = createFileRoute('/tools/$slug')({
  component: ToolDetailRouteComponent,
  loader: async ({ context, params }) => {
    const slug = params.slug.trim();

    // Orval client calls expect browser globals; skip prefetch during SSR.
    if (!slug || typeof window === 'undefined') {
      return;
    }

    await context.queryClient.ensureQueryData(toolQueryOptions(slug));
  },
});
