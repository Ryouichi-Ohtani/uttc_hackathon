import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/v1': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          // Core React libraries
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],

          // State management and data fetching
          'data-vendor': ['@tanstack/react-query', 'zustand', 'axios'],

          // UI and animation libraries
          'ui-vendor': ['framer-motion', 'react-hot-toast'],

          // 3D and graphics libraries
          '3d-vendor': ['three', '@react-three/fiber', '@react-three/drei'],

          // Charts and visualization
          'chart-vendor': ['chart.js', 'react-chartjs-2', 'recharts'],

          // Utilities
          'util-vendor': ['date-fns', 'clsx'],
        },
      },
    },
    // Increase chunk size warning limit
    chunkSizeWarningLimit: 1000,
    // Enable source maps for production debugging
    sourcemap: false,
  },
})
