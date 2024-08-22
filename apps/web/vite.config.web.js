import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'
import path from 'path'

export default defineConfig({
  root: path.resolve(__dirname, 'src/renderer'),
  build: {
    outDir: path.resolve(__dirname, 'dist-web'),
    emptyOutDir: true
  },
  plugins: [
    react(),
    TanStackRouterVite({
      routesDirectory: path.resolve(__dirname, 'src/renderer/src/routes'),
      generatedRouteTree: './src/renderer/src/routeTree.gen.ts'
    })
  ],
  resolve: {
    alias: {
      '@renderer': path.resolve(__dirname, 'src/renderer/src')
    }
  }
})
