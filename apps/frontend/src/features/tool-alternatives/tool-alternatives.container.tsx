import { useQuery } from '@tanstack/react-query';
import { Link } from '@tanstack/react-router';

import { toolAlternativesQueryOptions, toolQueryOptions } from '@/api/tools-query-options';
import { Alert } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';

import { ToolAlternativeCard } from './tool-alternatives/tool-alternative-card';

type ToolAlternativesContainerProps = {
  slug: string;
};

const isNotFoundError = (error: Error): boolean => {
  const normalizedMessage = error.message.toLowerCase();

  return normalizedMessage.includes('404') || normalizedMessage.includes('not found');
};

export const ToolAlternativesContainer = ({ slug }: ToolAlternativesContainerProps) => {
  const normalizedSlug = slug.trim();
  const alternativesQuery = useQuery(toolAlternativesQueryOptions(normalizedSlug));
  const detailQuery = useQuery(toolQueryOptions(normalizedSlug));

  const headerTitle = detailQuery.data ? `Alternatives for ${detailQuery.data.name}` : 'Tool alternatives';

  return (
    <main className="feature-page-shell">
      <header className="feature-page-header">
        <div>
          <p className="island-kicker mb-1">Instech atlas</p>
          <h1 className="text-2xl font-semibold">{headerTitle}</h1>
          {normalizedSlug ? <p className="text-muted-foreground mt-1 text-sm">{normalizedSlug}</p> : null}
        </div>
        {normalizedSlug ? (
          <Button asChild size="xs" variant="outline">
            <Link params={{ slug }} to="/tools/$slug">
              Back to tool details
            </Link>
          </Button>
        ) : null}
      </header>

      {normalizedSlug ? null : (
        <section className="feature-panel">
          <p className="text-muted-foreground text-sm">No tool slug was provided.</p>
        </section>
      )}

      {normalizedSlug && alternativesQuery.isLoading ? (
        <section className="feature-panel">
          <p className="text-muted-foreground text-sm">Loading alternatives...</p>
        </section>
      ) : null}

      {normalizedSlug && alternativesQuery.isError ? (
        <section>
          <Alert variant="destructive">
            {isNotFoundError(alternativesQuery.error) ? (
              <p>Tool not found for slug: {normalizedSlug}</p>
            ) : (
              <p>Failed to load alternatives. Please try again.</p>
            )}
          </Alert>
        </section>
      ) : null}

      {normalizedSlug &&
      !alternativesQuery.isLoading &&
      !alternativesQuery.isError &&
      (alternativesQuery.data?.length ?? 0) === 0 ? (
        <section className="feature-panel">
          <p className="text-muted-foreground text-sm">No alternative tools are listed for this tool yet.</p>
        </section>
      ) : null}

      {(alternativesQuery.data?.length ?? 0) > 0 ? (
        <section className="grid gap-3">
          {alternativesQuery.data?.map((item) => (
            <ToolAlternativeCard item={item} key={item.id} />
          ))}
        </section>
      ) : null}
    </main>
  );
};
