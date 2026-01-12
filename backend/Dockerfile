FROM golang:1.21-alpine AS builder

# 必要なパッケージのインストール
RUN apk add --no-cache git

# ワーキングディレクトリ設定
WORKDIR /app

# go.modとgo.sumをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# Airをインストール（ホットリロード用）
RUN go install github.com/cosmtrek/air@v1.49.0

# 開発環境用
FROM golang:1.21-alpine

RUN apk add --no-cache git

WORKDIR /app

# Airをコピー
COPY --from=builder /go/bin/air /usr/local/bin/air

# ポート公開
EXPOSE 8080

# デフォルトコマンド
CMD ["air", "-c", ".air.toml"]
