import path from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  base: "/check24-gendev-7/",
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
})
