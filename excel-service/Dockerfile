FROM python:3.11-slim

# 作業ディレクトリ設定
WORKDIR /app

# システム依存パッケージのインストール
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# requirements.txtをコピー
COPY requirements.txt .

# Python依存関係のインストール
RUN pip install --no-cache-dir -r requirements.txt

# アプリケーションコードをコピー
COPY . .

# ポート公開
EXPOSE 8000

# デフォルトコマンド
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]
