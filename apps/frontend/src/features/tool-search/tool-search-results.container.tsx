import { useQuery } from '@tanstack/react-query';
import { useEffect, useState } from 'react';

import { toolSearchQueryOptions } from '@/api/tools-query-options';
import { Alert } from '@/components/ui/alert';
import { Input } from '@/components/ui/input';

import { ToolSearchResultCard } from './tool-search-results/tool-search-result-card';

type ToolSearchResultsContainerProps = {
  onQueryChange: (q: string) => void;
  q: string;
};

const DEBOUNCE_MS = 300;

const useDebouncedValue = (value: string, delayMs: number): string => {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const timeout = window.setTimeout(() => {
      setDebouncedValue(value);
    }, delayMs);

    return () => window.clearTimeout(timeout);
  }, [delayMs, value]);

  return debouncedValue;
};

export const ToolSearchResultsContainer = ({ onQueryChange, q }: ToolSearchResultsContainerProps) => {
  const [draft, setDraft] = useState(q);
  const debouncedDraft = useDebouncedValue(draft, DEBOUNCE_MS);
  const normalizedDebouncedDraft = debouncedDraft.trim();
  const normalizedQuery = q.trim();

  const resultsQuery = useQuery(toolSearchQueryOptions(normalizedQuery));

  useEffect(() => {
    setDraft(q);
  }, [q]);

  useEffect(() => {
    if (normalizedDebouncedDraft === normalizedQuery) {
      return;
    }

    onQueryChange(normalizedDebouncedDraft);
  }, [normalizedDebouncedDraft, normalizedQuery, onQueryChange]);

  return (
    <main className="feature-page-shell">
      <header className="feature-page-header">
        <div>
          <p className="island-kicker mb-1">Instech atlas</p>
          <h1 className="text-2xl font-semibold">Search tools</h1>
        </div>
      </header>

      <section className="feature-panel space-y-3">
        <label className="text-xs font-medium uppercase" htmlFor="tool-search-input">
          Search tools
        </label>
        <Input
          id="tool-search-input"
          onChange={(event) => setDraft(event.target.value)}
          placeholder="Type a keyword..."
          value={draft}
        />
        <p className="text-muted-foreground text-xs">Search updates after {DEBOUNCE_MS}ms to avoid request churn.</p>
      </section>

      {normalizedQuery.length === 0 ? (
        <section className="feature-panel mt-3">
          <p className="text-muted-foreground text-sm">Enter a keyword to search tools.</p>
        </section>
      ) : null}

      {normalizedQuery.length > 0 && resultsQuery.isLoading ? (
        <section className="feature-panel mt-3">
          <p className="text-muted-foreground text-sm">Searching...</p>
        </section>
      ) : null}

      {normalizedQuery.length > 0 && resultsQuery.isError ? (
        <section className="mt-3">
          <Alert variant="destructive">Failed to load search results. Try another keyword.</Alert>
        </section>
      ) : null}

      {normalizedQuery.length > 0 &&
      !resultsQuery.isLoading &&
      !resultsQuery.isError &&
      (resultsQuery.data?.length ?? 0) === 0 ? (
        <section className="feature-panel mt-3">
          <p className="text-muted-foreground text-sm">
            No tools matched <strong>{normalizedQuery}</strong>.
          </p>
        </section>
      ) : null}

      {normalizedQuery.length > 0 && resultsQuery.isFetching ? (
        <p className="text-muted-foreground mt-3 text-xs">Refreshing results...</p>
      ) : null}

      {(resultsQuery.data?.length ?? 0) > 0 ? (
        <section className="mt-3 grid gap-3">
          {resultsQuery.data?.map((item) => (
            <ToolSearchResultCard item={item} key={item.id} />
          ))}
        </section>
      ) : null}
    </main>
  );
};
