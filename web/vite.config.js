import { svelte } from '@sveltejs/vite-plugin-svelte';
import { defineConfig } from 'vite';
import { resolve } from 'node:path';

export default defineConfig({
  // The Go server serves built assets under /static/; the Go shell references
  // the entry script at /static/assets/<hash>.js (see internal/render/assets.go).
  // A relative base keeps chunk loader URLs relative to that script, so
  // dynamic imports resolve to /static/assets/... at the origin root AND to
  // <mount>/static/assets/... when a reverse proxy mounts the app under a URL
  // prefix (X-Forwarded-Prefix) — no rebuild per mount point.
  base: './',
  plugins: [svelte()],
  build: {
    manifest: true,
    outDir: 'dist',
    emptyOutDir: true,
    chunkSizeWarningLimit: 1200,
    rollupOptions: {
      input: {
        app: resolve(__dirname, 'src/main.js'),
      },
      output: {
        manualChunks(id) {
          if (id.includes('highlight.js')) return 'hljs';
          // NOTE: do NOT force @pierre/diffs / shiki into a single manual chunk.
          // The dynamic import('@pierre/diffs') already splits it off, and shiki
          // lazy-loads individual language grammars as separate chunks. Grouping
          // them all here inlines every grammar into one ~10 MB blob that stalls
          // the diff modal on slow/remote links.
        },
      },
    },
  },
});
