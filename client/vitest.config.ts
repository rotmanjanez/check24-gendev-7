import { defineConfig } from 'vitest/config';
import { configDefaults } from 'vitest/config';
import vue from '@vitejs/plugin-vue';
import path from 'node:path';
import tailwindcss from '@tailwindcss/vite';
import vueI18n from '@intlify/unplugin-vue-i18n/vite';

export default defineConfig({
  plugins: [vue(), tailwindcss(), vueI18n({
    runtimeOnly: false,
    include: [path.resolve(__dirname, './src/i18n/locales/**')],
  })],
  base: "/check24-gendev-7/",
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://gendev-server:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('src/api')) {
            return 'openapi-client';
          }
          if (id.includes('node_modules')) {
            return 'vendor';
          }
        },
      },
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html', 'json-summary'],
      reportsDirectory: './coverage',
    },
    exclude: [...configDefaults.exclude, 'e2e/**'],
  },
});
