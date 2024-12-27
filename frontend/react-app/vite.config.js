import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000, // フロントエンド開発サーバーをポート3000で起動
    host: true, // 必須: Dockerで動作する場合は必要 (0.0.0.0をバインド)
  },
})
