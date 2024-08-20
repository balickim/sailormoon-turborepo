import { resolve } from 'path'
import { defineConfig, externalizeDepsPlugin } from 'electron-vite'
import viteReact from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  main: {
    plugins: [externalizeDepsPlugin()]
  },
  preload: {
    plugins: [
      externalizeDepsPlugin(),
      TanStackRouterVite({ routesDirectory: path.resolve(__dirname, 'src/renderer/src/routes') }),
      viteReact()
    ]
  },
  renderer: {
    resolve: {
      alias: {
        '@renderer': resolve('src/renderer/src')
      }
    },
    plugins: [react()]
  }
})
