import { useQuery } from '@tanstack/react-query';

import { toolQueryOptions } from '@/api/tools-query-options';
import { Alert } from '@/components/ui/alert';
import { ToolDetailActions } from './tool-detail/tool-detail-actions';
import { ToolDetailOverviewCard } from './tool-detail/tool-detail-overview-card';

type ToolDetailContainerProps = {
  slug: string;
};

const toLabel = (value: string): string => `${value.charAt(0).toUpperCase()}${value.slice(1).replaceAll('_', ' ')}`;

const isNotFoundError = (error: Error): boolean => {
  const normalizedMessage = error.message.toLowerCase();

  return normalizedMessage.includes('404') || normalizedMessage.includes('not found');
};

export const ToolDetailContainer = ({ slug }: ToolDetailContainerProps) => {
  const normalizedSlug = slug.trim();
  const detailQuery = useQuery(toolQueryOptions(normalizedSlug));

  return (
    <main className="feature-page-shell">
      <header className="feature-page-header">
        <div>
          <p className="island-kicker mb-1">Instech atlas</p>
          <h1 className="text-2xl font-semibold">Tool details</h1>
        </div>
      </header>

      {normalizedSlug ? null : (
        <section className="feature-panel">
          <p className="text-muted-foreground text-sm">No tool slug was provided.</p>
        </section>
      )}

      {normalizedSlug && detailQuery.isLoading ? (
        <section className="feature-panel">
          <p className="text-muted-foreground text-sm">Loading tool details...</p>
        </section>
      ) : null}

      {normalizedSlug && detailQuery.isError ? (
        <section>
          <Alert variant="destructive">
            {isNotFoundError(detailQuery.error) ? (
              <p>Tool not found for slug: {normalizedSlug}</p>
            ) : (
              <p>Failed to load tool details. Please try again.</p>
            )}
          </Alert>
        </section>
      ) : null}

      {detailQuery.data ? (
        <section className="grid gap-3">
          <ToolDetailOverviewCard
            categoryLabel={toLabel(detailQuery.data.category)}
            details={detailQuery.data.details}
            devStatusLabel={toLabel(detailQuery.data.devStatus)}
            github={detailQuery.data.github}
            name={detailQuery.data.name}
            prolang={detailQuery.data.prolang}
            releaseYear={detailQuery.data.releaseYear}
            slug={detailQuery.data.slug}
            subTypeLabel={toLabel(detailQuery.data.subType)}
            tags={detailQuery.data.tags}
            useCases={detailQuery.data.useCases}
            website={detailQuery.data.website}
          />
          <ToolDetailActions slug={detailQuery.data.slug} />
        </section>
      ) : null}
    </main>
  );
};
