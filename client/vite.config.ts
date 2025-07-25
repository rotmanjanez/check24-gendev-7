import path from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import vueI18n from '@intlify/unplugin-vue-i18n/vite'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

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
        target: `http://${process.env.VITE_BACKEND_HOST || 'localhost'}:8080`, // Backend server
        changeOrigin: true, 
        rewrite: (path) => path.replace(/^\/api/, ''), // Remove '/api' prefix
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
          if (id.includes('node_modules/vue')) {
            return 'vue';
          }
          if (id.includes('node_modules/vue-router')) {
            return 'vue-router';
          }
          if (id.includes('node_modules/@intlify')) {
            return 'vue-i18n';
          }
          if (id.includes('node_modules/@tailwindcss')) {
            return 'tailwindcss';
          }
          if (id.includes('node_modules/posthog-js')) {
            return 'vendor-posthog';
          }
          if (id.includes('node_modules/reka-ui')) {
            return 'vendor-reka-ui';
          }
          if (id.includes('node_modules')) {
            return 'vendor';
          }
        },
      },
    },
  },
})
