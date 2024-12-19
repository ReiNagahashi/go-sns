# ベースイメージを指定
FROM golang:1.23.4

# 作業ディレクトリを設定
WORKDIR /backend

# Goモジュールとアプリケーションコードをコピー
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# アプリケーションをビルド
RUN go build -o main .

# 開放するポート
EXPOSE 8080

# バックエンドを起動
CMD ["./main"]
