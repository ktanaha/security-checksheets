"""Excel処理APIのPydanticモデル定義"""

from typing import List, Optional, Any
from pydantic import BaseModel, Field


class SheetInfo(BaseModel):
    """シート情報"""
    name: str = Field(..., description="シート名")
    index: int = Field(..., description="シートのインデックス（0始まり）")
    row_count: int = Field(..., description="行数")
    column_count: int = Field(..., description="列数")


class CellData(BaseModel):
    """セルデータ"""
    row: int = Field(..., description="行番号（1始まり）")
    column: int = Field(..., description="列番号（1始まり）")
    value: Optional[Any] = Field(None, description="セルの値")
    formatted_value: Optional[str] = Field(None, description="フォーマット済みの値")
    is_merged: bool = Field(False, description="結合セルかどうか")
    merge_range: Optional[str] = Field(None, description="結合セル範囲（例: A1:B2）")


class SheetPreviewRequest(BaseModel):
    """シートプレビューリクエスト"""
    file_path: str = Field(..., description="Excelファイルのパス")
    sheet_name: Optional[str] = Field(None, description="シート名（未指定の場合は最初のシート）")
    start_row: int = Field(1, ge=1, description="開始行（1始まり）")
    end_row: Optional[int] = Field(None, ge=1, description="終了行（未指定の場合はすべて）")
    start_column: int = Field(1, ge=1, description="開始列（1始まり）")
    end_column: Optional[int] = Field(None, ge=1, description="終了列（未指定の場合はすべて）")


class SheetPreviewResponse(BaseModel):
    """シートプレビューレスポンス"""
    sheet_name: str = Field(..., description="シート名")
    cells: List[CellData] = Field(..., description="セルデータのリスト")
    row_count: int = Field(..., description="取得した行数")
    column_count: int = Field(..., description="取得した列数")


class ParseExcelRequest(BaseModel):
    """Excel解析リクエスト"""
    file_path: str = Field(..., description="Excelファイルのパス")


class ParseExcelResponse(BaseModel):
    """Excel解析レスポンス"""
    file_name: str = Field(..., description="ファイル名")
    file_path: str = Field(..., description="ファイルパス")
    sheets: List[SheetInfo] = Field(..., description="シート情報のリスト")
    total_sheets: int = Field(..., description="総シート数")


class ErrorResponse(BaseModel):
    """エラーレスポンス"""
    error: str = Field(..., description="エラーメッセージ")
    detail: Optional[str] = Field(None, description="エラー詳細")


class ExtractQARequest(BaseModel):
    """Q/A抽出リクエスト"""
    file_path: str = Field(..., description="Excelファイルのパス")
    sheet_name: str = Field(..., description="シート名")
    start_row: int = Field(..., ge=1, description="開始行（1始まり）")
    end_row: int = Field(..., ge=1, description="終了行")
    question_column: int = Field(..., ge=1, description="質問列番号（1始まり）")
    answer_column: int = Field(..., ge=1, description="回答列番号（1始まり）")
    department_column: Optional[int] = Field(None, ge=1, description="担当部門列番号（1始まり、任意）")
    skip_header_rows: int = Field(1, ge=0, description="スキップするヘッダー行数")


class QAItem(BaseModel):
    """抽出されたQ/Aアイテム"""
    row_number: int = Field(..., description="元の行番号")
    question: str = Field(..., description="質問")
    answer: Optional[str] = Field(None, description="回答")
    department: Optional[str] = Field(None, description="担当部門")


class ExtractQAResponse(BaseModel):
    """Q/A抽出レスポンス"""
    file_path: str = Field(..., description="ファイルパス")
    sheet_name: str = Field(..., description="シート名")
    source_range: str = Field(..., description="抽出元の範囲")
    items: List[QAItem] = Field(..., description="抽出されたQ/Aアイテムのリスト")
    total_items: int = Field(..., description="抽出されたアイテム数")
