import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    proxy: {
      '/rest': {
        target: 'https://gonic.ekskog.net',
        changeOrigin: true,
        secure: true
      },
      '/deezer-api': {
        target: 'https://api.deezer.com',
        changeOrigin: true,
        rewrite: path => path.replace(/^\/deezer-api/, '')
      }
    }
  },
})