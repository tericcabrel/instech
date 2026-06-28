import tailwindcss from '@tailwindcss/vite';
import { devtools } from '@tanstack/devtools-vite';

import { tanstackStart } from '@tanstack/react-start/plugin/vite';

import viteReact from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

const API_PATH_REGEXP = /^\/api/;

const config = defineConfig({
  plugins: [devtools(), tailwindcss(), tanstackStart(), viteReact()],
  resolve: { tsconfigPaths: true },
  server: {
    proxy: {
      '/api': {
        changeOrigin: true,
        rewrite: (path) => path.replace(API_PATH_REGEXP, ''),
        target: 'http://localhost:8801',
      },
    },
  },
});

export default config;
