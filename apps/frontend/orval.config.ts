import { defineConfig } from 'orval';

export default defineConfig({
  instech: {
    input: '../backend/openapi.yaml',
    output: {
      client: 'fetch',
      mock: false,
      mode: 'split',
      namingConvention: 'PascalCase',
      override: {
        mutator: {
          name: 'customInstance',
          path: './src/api/mutator.ts',
        },
        zod: {
          generate: {
            body: true,
            param: true,
            query: true,
            response: true,
          },
        },
      },
      schemas: {
        path: './src/api/generated/model',
        type: 'zod',
      },
      target: './src/api/generated/instech.ts',
    },
  },
});
