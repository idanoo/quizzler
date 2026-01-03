import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  server: {
    host: '0.0.0.0',
    port: process.env.FRONTEND_PORT || 5172,
    proxy: {
      '/api': {
        target: 'http://app:8132',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
  },
});

