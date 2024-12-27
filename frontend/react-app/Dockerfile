# ベースイメージを指定
FROM node:20

# 作業ディレクトリを設定
WORKDIR /app

# 必要なファイルをコピー
COPY package.json package-lock.json ./
RUN npm install

# アプリケーションコードをコピー
COPY . ./

EXPOSE 3000


# ビルド
CMD ["npm", "run", "dev", "--", "--host"]

