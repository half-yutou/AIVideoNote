import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, resolve(__dirname, '..'), '')
  const port = parseInt(env.FRONTEND_PORT || '3015', 10)
  const backendPort = env.BACKEND_PORT || '8080'
  const backend = `http://127.0.0.1:${backendPort}`

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src'),
      },
    },
    server: {
      port,
      proxy: {
        '/api': backend,
        '/uploads': backend,
      },
    },
  }
})
