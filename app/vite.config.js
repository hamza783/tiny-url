import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': {
        // When running inside Docker, use the service name, not localhost
        target: 'http://api-gateway:8080',
        changeOrigin: true
      },
    },
  },
})
