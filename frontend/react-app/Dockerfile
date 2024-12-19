# ベースイメージを指定
FROM node:20

# 作業ディレクトリを設定
WORKDIR /react-app

# 必要なファイルをコピー
COPY package*.json ./
RUN npm install

# アプリケーションコードをコピー
COPY . .

# ビルド
RUN npm run build

# 開放するポート
EXPOSE 3000

# サーバーを起動
CMD ["npx", "serve", "-s", "build"]
