FROM node:18-alpine

# 作業ディレクトリ設定
WORKDIR /app

# package.jsonとpackage-lock.jsonをコピー
COPY package*.json ./

# 依存関係のインストール
RUN npm install

# アプリケーションコードをコピー
COPY . .

# ポート公開
EXPOSE 3000

# デフォルトコマンド
CMD ["npm", "run", "dev", "--", "--host"]
