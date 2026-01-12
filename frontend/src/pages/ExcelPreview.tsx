import React, { useState, useEffect } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { parseExcelFile, getSheetPreview } from '../services/excelService';
import ExcelTable from '../components/organisms/ExcelTable';
import Button from '../components/atoms/Button';
import type { SheetInfo, SheetPreviewResponse, CellRange } from '../types/excel';

/**
 * Excelプレビューページ
 */
const ExcelPreview: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [filePath, setFilePath] = useState('');
  const [sheets, setSheets] = useState<SheetInfo[]>([]);
  const [selectedSheet, setSelectedSheet] = useState<string | null>(null);
  const [previewData, setPreviewData] = useState<SheetPreviewResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedRange, setSelectedRange] = useState<CellRange | null>(null);

  // クエリパラメータからファイルパスを取得し、自動で解析
  useEffect(() => {
    const filePathParam = searchParams.get('filePath');
    if (filePathParam) {
      setFilePath(filePathParam);
      // 自動で解析を実行
      parseFile(filePathParam);
    }
  }, [searchParams]);

  // Excelファイル解析（共通関数）
  const parseFile = async (path: string) => {
    if (!path) {
      setError('ファイルパスを入力してください');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await parseExcelFile(path);
      setSheets(result.sheets);
      if (result.sheets.length > 0) {
        setSelectedSheet(result.sheets[0].name);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ファイルの解析に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // Excelファイル解析（ボタンクリック用）
  const handleParseFile = async () => {
    await parseFile(filePath);
  };

  // シートプレビュー取得
  useEffect(() => {
    if (!filePath || !selectedSheet) return;

    const fetchPreview = async () => {
      setLoading(true);
      setError(null);

      try {
        const result = await getSheetPreview({
          filePath,
          sheetName: selectedSheet,
        });
        setPreviewData(result);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'プレビューの取得に失敗しました');
      } finally {
        setLoading(false);
      }
    };

    fetchPreview();
  }, [filePath, selectedSheet]);

  // 範囲選択ハンドラー
  const handleRangeSelect = (range: CellRange) => {
    setSelectedRange(range);
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto p-6">
        {/* ヘッダー */}
        <div className="mb-6">
          <Button variant="secondary" onClick={() => navigate(-1)}>
            ← 戻る
          </Button>
        </div>

        <h1 className="text-3xl font-bold mb-6">Excelプレビュー</h1>

        {/* ファイルパス入力 */}
      <div className="mb-6">
        <div className="flex gap-2">
          <input
            type="text"
            value={filePath}
            onChange={(e) => setFilePath(e.target.value)}
            placeholder="/app/uploads/test.xlsx"
            className="flex-1 px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            onClick={handleParseFile}
            disabled={loading}
            className="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"
          >
            {loading ? '解析中...' : '解析'}
          </button>
        </div>
      </div>

      {/* エラー表示 */}
      {error && (
        <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
          {error}
        </div>
      )}

      {/* シート選択 */}
      {sheets.length > 0 && (
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            シート選択
          </label>
          <select
            value={selectedSheet || ''}
            onChange={(e) => setSelectedSheet(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {sheets.map((sheet) => (
              <option key={sheet.index} value={sheet.name}>
                {sheet.name} ({sheet.row_count}行 × {sheet.column_count}列)
              </option>
            ))}
          </select>
        </div>
      )}

      {/* 選択範囲表示 */}
      {selectedRange && (
        <div className="mb-4 p-4 bg-blue-100 border border-blue-400 text-blue-700 rounded">
          選択範囲: 行 {selectedRange.startRow}~{selectedRange.endRow}, 列{' '}
          {selectedRange.startColumn}~{selectedRange.endColumn}
        </div>
      )}

      {/* プレビュー表示 */}
      {loading && <div className="text-center py-8">読み込み中...</div>}

      {previewData && !loading && (
        <div className="bg-white rounded shadow">
          <div className="p-4 border-b border-gray-200">
            <h2 className="text-xl font-semibold">
              {previewData.sheet_name} - プレビュー
            </h2>
            <p className="text-sm text-gray-600">
              {previewData.total_rows}行 × {previewData.total_columns}列
            </p>
          </div>
          <div className="p-4">
            <ExcelTable
              cells={previewData.cells}
              totalRows={previewData.total_rows}
              totalColumns={previewData.total_columns}
              onRangeSelect={handleRangeSelect}
            />
          </div>
        </div>
      )}
      </div>
    </div>
  );
};

export default ExcelPreview;
