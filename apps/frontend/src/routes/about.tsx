import { createFileRoute } from '@tanstack/react-router';

const About = () => (
  <main className="mx-auto w-full max-w-3xl px-6 py-12">
    <h1 className="mb-3 text-3xl font-semibold">About</h1>
    <p className="text-sm text-neutral-600 dark:text-neutral-300">
      This is a minimal TanStack Start app generated from TanStack CLI and configured with TanStack Router, TanStack
      Query, and TanStack Form.
    </p>
  </main>
);

export const Route = createFileRoute('/about')({
  component: About,
});
