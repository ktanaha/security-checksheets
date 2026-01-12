"""
Excel処理専用FastAPIサービス
セキュリティチェックシート抽出アプリのExcel解析を担当
"""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from datetime import datetime
import os

# ルートのインポート
from app.api.routes import router as excel_router

# FastAPIアプリケーションの初期化
app = FastAPI(
    title="Security Checksheets Excel Service",
    description="Excel処理専用APIサービス",
    version="1.0.0"
)

# CORS設定
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "http://localhost:8080"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# ルートの登録
app.include_router(excel_router)


@app.get("/")
async def root():
    """ルートエンドポイント"""
    return {
        "service": "security-checksheets-excel-service",
        "version": "1.0.0",
        "status": "running"
    }


@app.get("/health")
async def health_check():
    """ヘルスチェックエンドポイント"""
    return {
        "status": "healthy",
        "service": "security-checksheets-excel-service",
        "time": datetime.now().isoformat()
    }


@app.get("/api/ping")
async def ping():
    """疎通確認エンドポイント"""
    return {"message": "pong"}


# アプリケーション起動時の処理
@app.on_event("startup")
async def startup_event():
    """アプリケーション起動時の初期化処理"""
    print("=" * 50)
    print("Excel Service starting...")
    print(f"Environment: {os.getenv('ENV', 'development')}")
    print("=" * 50)


# アプリケーション終了時の処理
@app.on_event("shutdown")
async def shutdown_event():
    """アプリケーション終了時のクリーンアップ処理"""
    print("=" * 50)
    print("Excel Service shutting down...")
    print("=" * 50)


if __name__ == "__main__":
    import uvicorn

    port = int(os.getenv("PORT", "8000"))
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=port,
        reload=True
    )
