import { defineConfig } from 'orval'

export default defineConfig({
  instech: {
    input: '../backend/openapi.yaml',
    output: {
      mode: 'split',
      target: './src/api/generated/instech.ts',
      schemas: './src/api/generated/model',
      client: 'fetch',
      mock: false,
      override: {
        mutator: {
          path: './src/api/mutator.ts',
          name: 'customInstance',
        },
        zod: {
          generate: {
            body: true,
            response: true,
            query: true,
            param: true,
          },
        },
      },
    },
  },
})
