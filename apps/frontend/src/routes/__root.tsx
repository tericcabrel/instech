import { TanStackDevtools } from '@tanstack/react-devtools';
import type { QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtoolsPanel } from '@tanstack/react-query-devtools';
import { createRootRouteWithContext, HeadContent, Outlet, Scripts } from '@tanstack/react-router';
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools';
import appCss from '../styles.css?url';

type MyRouterContext = {
  queryClient: QueryClient;
};

const RootDocument = ({ children }: { children: React.ReactNode }) => (
  <html lang="en">
    <head>
      <HeadContent />
    </head>
    <body className="font-sans antialiased">
      {children ?? <Outlet />}
      <TanStackDevtools
        config={{
          position: 'bottom-right',
        }}
        plugins={[
          {
            name: 'Tanstack Router',
            render: <TanStackRouterDevtoolsPanel />,
          },
          {
            name: 'Tanstack Query',
            render: <ReactQueryDevtoolsPanel />,
          },
        ]}
      />
      <Scripts />
    </body>
  </html>
);

const RootNotFound = () => (
  <div className="mx-auto flex min-h-screen w-full max-w-2xl flex-col items-start justify-center gap-3 p-6">
    <h1 className="text-xl font-semibold">Page not found</h1>
    <p className="text-sm text-muted-foreground">The page you are looking for does not exist or has moved.</p>
    <a className="inline-flex h-9 items-center rounded border px-3 text-sm" href="/">
      Go to home
    </a>
  </div>
);

const RootError = ({ error, reset }: { error: Error; reset: () => void }) => (
  <div className="mx-auto flex min-h-screen w-full max-w-2xl flex-col items-start justify-center gap-3 p-6">
    <h1 className="text-xl font-semibold">Something went wrong</h1>
    <p className="text-sm text-muted-foreground">We could not render this page because of an unexpected error.</p>
    <pre className="w-full overflow-auto rounded border bg-muted p-3 text-xs">{error.message}</pre>
    <button className="inline-flex h-9 items-center rounded border px-3 text-sm" onClick={() => reset()} type="button">
      Try again
    </button>
  </div>
);

export const Route = createRootRouteWithContext<MyRouterContext>()({
  errorComponent: RootError,
  head: () => ({
    links: [
      {
        href: appCss,
        rel: 'stylesheet',
      },
    ],
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        content: 'width=device-width, initial-scale=1',
        name: 'viewport',
      },
      {
        title: 'Instech Web',
      },
    ],
  }),
  notFoundComponent: RootNotFound,
  shellComponent: RootDocument,
});
