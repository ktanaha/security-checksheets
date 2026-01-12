// Excel関連の型定義

/**
 * セルデータ
 */
export interface CellData {
  row: number;
  column: number;
  value: string | number | null;
  formatted_value: string;
  is_merged: boolean;
  merge_range?: string;
}

/**
 * シート情報
 */
export interface SheetInfo {
  name: string;
  index: number;
  row_count: number;
  column_count: number;
}

/**
 * Excel解析レスポンス
 */
export interface ParseExcelResponse {
  file_name: string;
  sheets: SheetInfo[];
  total_sheets: number;
}

/**
 * シートプレビューレスポンス
 */
export interface SheetPreviewResponse {
  sheet_name: string;
  cells: CellData[];
  total_rows: number;
  total_columns: number;
}

/**
 * プレビューリクエストパラメータ
 */
export interface PreviewRequestParams {
  filePath: string;
  sheetName?: string;
  startRow?: number;
  endRow?: number;
  startColumn?: number;
  endColumn?: number;
}

/**
 * セル範囲選択
 */
export interface CellRange {
  startRow: number;
  endRow: number;
  startColumn: number;
  endColumn: number;
}

/**
 * エラーレスポンス
 */
export interface ErrorResponse {
  detail: string;
}
