"""Excel処理APIのルート定義"""

from fastapi import APIRouter, HTTPException, status
from app.api.models import (
    ParseExcelRequest,
    ParseExcelResponse,
    SheetPreviewRequest,
    SheetPreviewResponse,
    SheetInfo,
    CellData,
    ErrorResponse
)
from app.services.excel_parser import ExcelParser

router = APIRouter(prefix="/excel", tags=["excel"])
excel_parser = ExcelParser()


@router.post("/parse", response_model=ParseExcelResponse)
async def parse_excel(request: ParseExcelRequest):
    """
    Excelファイルを解析してシート情報を取得する

    Args:
        request: Excel解析リクエスト

    Returns:
        Excel解析結果

    Raises:
        HTTPException: ファイルが見つからない場合やエラーが発生した場合
    """
    try:
        result = excel_parser.parse_excel(request.file_path)

        return ParseExcelResponse(
            file_name=result['file_name'],
            file_path=result['file_path'],
            sheets=[
                SheetInfo(
                    name=sheet['name'],
                    index=sheet['index'],
                    row_count=sheet['row_count'],
                    column_count=sheet['column_count']
                )
                for sheet in result['sheets']
            ],
            total_sheets=result['total_sheets']
        )

    except FileNotFoundError as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(e)
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Excel解析中にエラーが発生しました: {str(e)}"
        )


@router.post("/preview", response_model=SheetPreviewResponse)
async def get_sheet_preview(request: SheetPreviewRequest):
    """
    シートのプレビューデータを取得する

    Args:
        request: シートプレビューリクエスト

    Returns:
        シートプレビューデータ

    Raises:
        HTTPException: ファイルが見つからない場合やエラーが発生した場合
    """
    try:
        result = excel_parser.get_sheet_preview(
            file_path=request.file_path,
            sheet_name=request.sheet_name,
            start_row=request.start_row,
            end_row=request.end_row,
            start_column=request.start_column,
            end_column=request.end_column
        )

        return SheetPreviewResponse(
            sheet_name=result['sheet_name'],
            cells=[
                CellData(
                    row=cell['row'],
                    column=cell['column'],
                    value=cell['value'],
                    formatted_value=cell['formatted_value'],
                    is_merged=cell['is_merged'],
                    merge_range=cell['merge_range']
                )
                for cell in result['cells']
            ],
            row_count=result['row_count'],
            column_count=result['column_count']
        )

    except FileNotFoundError as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(e)
        )
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"シートプレビュー取得中にエラーが発生しました: {str(e)}"
        )
