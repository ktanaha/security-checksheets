"""Excel処理APIのルート定義"""

from fastapi import APIRouter, HTTPException, status
from app.api.models import (
    ParseExcelRequest,
    ParseExcelResponse,
    SheetPreviewRequest,
    SheetPreviewResponse,
    SheetInfo,
    CellData,
    ErrorResponse,
    ExtractQARequest,
    ExtractQAResponse,
    QAItem
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


@router.post("/extract-qa", response_model=ExtractQAResponse)
async def extract_qa(request: ExtractQARequest):
    """
    ExcelシートからQ/Aを抽出する

    Args:
        request: Q/A抽出リクエスト

    Returns:
        抽出されたQ/Aデータ

    Raises:
        HTTPException: ファイルが見つからない場合やエラーが発生した場合
    """
    try:
        result = excel_parser.extract_qa(
            file_path=request.file_path,
            sheet_name=request.sheet_name,
            start_row=request.start_row,
            end_row=request.end_row,
            question_column=request.question_column,
            answer_column=request.answer_column,
            department_column=request.department_column,
            skip_header_rows=request.skip_header_rows
        )

        return ExtractQAResponse(
            file_path=result['file_path'],
            sheet_name=result['sheet_name'],
            source_range=result['source_range'],
            items=[
                QAItem(
                    row_number=item['row_number'],
                    question=item['question'],
                    answer=item['answer'],
                    department=item['department']
                )
                for item in result['items']
            ],
            total_items=result['total_items']
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
            detail=f"Q/A抽出中にエラーが発生しました: {str(e)}"
        )
