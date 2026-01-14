import axios from 'axios';
import type {
  ParseExcelResponse,
  SheetPreviewResponse,
  PreviewRequestParams,
} from '../types/excel';
import type { ExtractQARequest, ExtractQAResponse } from '../types/knowledge';

const EXCEL_SERVICE_URL = 'http://localhost:8000';

/**
 * Excelファイルを解析してシート情報を取得
 * @param filePath ファイルパス
 * @returns Excel解析結果
 */
export const parseExcelFile = async (filePath: string): Promise<ParseExcelResponse> => {
  const response = await axios.post<ParseExcelResponse>(`${EXCEL_SERVICE_URL}/excel/parse`, {
    file_path: filePath,
  });
  return response.data;
};

/**
 * シートのプレビューデータを取得
 * @param params プレビューリクエストパラメータ
 * @returns シートプレビューデータ
 */
export const getSheetPreview = async (
  params: PreviewRequestParams
): Promise<SheetPreviewResponse> => {
  const requestBody: Record<string, unknown> = {
    file_path: params.filePath,
  };

  // オプショナルパラメータの追加
  if (params.sheetName !== undefined) {
    requestBody.sheet_name = params.sheetName;
  }
  if (params.startRow !== undefined) {
    requestBody.start_row = params.startRow;
  }
  if (params.endRow !== undefined) {
    requestBody.end_row = params.endRow;
  }
  if (params.startColumn !== undefined) {
    requestBody.start_column = params.startColumn;
  }
  if (params.endColumn !== undefined) {
    requestBody.end_column = params.endColumn;
  }

  const response = await axios.post<SheetPreviewResponse>(
    `${EXCEL_SERVICE_URL}/excel/preview`,
    requestBody
  );
  return response.data;
};

/**
 * ExcelシートからQ/Aを抽出
 * @param request Q/A抽出リクエスト
 * @returns Q/A抽出結果
 */
export const extractQA = async (request: ExtractQARequest): Promise<ExtractQAResponse> => {
  const response = await axios.post<ExtractQAResponse>(
    `${EXCEL_SERVICE_URL}/excel/extract-qa`,
    request
  );
  return response.data;
};
