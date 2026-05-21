import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    allowedHosts: [
      'thegioibatdong.site',
      'api.thegioibatdong.site',
      'assets.thegioibatdong.site'
    ]
  }
})
