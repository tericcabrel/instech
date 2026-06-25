import { defineConfig } from 'orval'

export default defineConfig({
  instech: {
    input: '../backend/openapi.yaml',
    output: {
      mode: 'split',
      target: './src/api/generated/instech.ts',
      schemas: {
        path: './src/api/generated/model',
        type: 'zod',
      },
      client: 'fetch',
      mock: false,
      override: {
        mutator: {
          path: './src/api/mutator.ts',
          name: 'customInstance',
        },
      },
    },
  },
})
