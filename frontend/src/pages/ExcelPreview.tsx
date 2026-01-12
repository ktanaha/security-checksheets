import React, { useState, useEffect } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { parseExcelFile, getSheetPreview } from '../services/excelService';
import { getProjects, uploadFile } from '../services/projectService';
import ExcelTable from '../components/organisms/ExcelTable';
import Button from '../components/atoms/Button';
import type { SheetInfo, SheetPreviewResponse, CellRange } from '../types/excel';
import type { Project } from '../types/project';

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

  // ファイルアップロード関連の状態
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadedBy, setUploadedBy] = useState('');

  // 案件一覧を取得
  useEffect(() => {
    const fetchProjects = async () => {
      try {
        const data = await getProjects();
        setProjects(data);
        if (data.length > 0) {
          setSelectedProjectId(data[0].id);
          setUploadedBy(data[0].owner);
        }
      } catch (err) {
        console.error('案件一覧の取得に失敗しました:', err);
      }
    };
    fetchProjects();
  }, []);

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

  // ファイル選択ハンドラー
  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Excelファイルのみ許可
    const allowedTypes = [
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', // .xlsx
      'application/vnd.ms-excel', // .xls
    ];

    if (!allowedTypes.includes(file.type)) {
      setError('Excelファイル（.xlsx, .xls）のみアップロード可能です');
      e.target.value = '';
      return;
    }

    setSelectedFile(file);
    setError(null);
  };

  // ファイルアップロードハンドラー
  const handleFileUpload = async () => {
    if (!selectedFile || !selectedProjectId) {
      setError('ファイルと案件を選択してください');
      return;
    }

    setUploading(true);
    setError(null);

    try {
      const uploadedFile = await uploadFile(selectedProjectId, {
        file: selectedFile,
        uploaded_by: uploadedBy || '匿名',
      });

      // アップロード成功後、自動的にプレビューを表示
      setFilePath(uploadedFile.file_path);
      await parseFile(uploadedFile.file_path);

      // 成功メッセージ
      alert(`ファイル「${uploadedFile.file_name}」をアップロードしました`);
      setSelectedFile(null);

      // ファイル入力をリセット
      const fileInput = document.getElementById('file-upload') as HTMLInputElement;
      if (fileInput) fileInput.value = '';
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ファイルのアップロードに失敗しました');
    } finally {
      setUploading(false);
    }
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

        {/* ファイルアップロードセクション */}
        <div className="mb-8 bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-semibold mb-4">ファイルアップロード</h2>

          {/* 案件選択 */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              案件を選択 <span className="text-red-500">*</span>
            </label>
            <select
              value={selectedProjectId || ''}
              onChange={(e) => {
                const projectId = parseInt(e.target.value);
                setSelectedProjectId(projectId);
                const project = projects.find((p) => p.id === projectId);
                if (project) setUploadedBy(project.owner);
              }}
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={uploading}
            >
              <option value="">案件を選択してください</option>
              {projects.map((project) => (
                <option key={project.id} value={project.id}>
                  {project.customer_name} - {project.owner}
                </option>
              ))}
            </select>
          </div>

          {/* ファイル選択とアップロード */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Excelファイル <span className="text-red-500">*</span>
              </label>
              <input
                id="file-upload"
                type="file"
                accept=".xlsx,.xls"
                onChange={handleFileSelect}
                disabled={uploading}
                className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 disabled:opacity-50"
              />
              {selectedFile && (
                <p className="mt-2 text-sm text-gray-600">
                  選択中: {selectedFile.name} ({(selectedFile.size / 1024).toFixed(1)} KB)
                </p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                アップロード者
              </label>
              <input
                type="text"
                value={uploadedBy}
                onChange={(e) => setUploadedBy(e.target.value)}
                disabled={uploading}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
                placeholder="名前を入力"
              />
            </div>
          </div>

          <div className="mt-4">
            <Button
              variant="primary"
              onClick={handleFileUpload}
              disabled={!selectedFile || !selectedProjectId || uploading}
            >
              {uploading ? 'アップロード中...' : 'アップロードしてプレビュー'}
            </Button>
          </div>
        </div>

        {/* ファイルパス直接入力（既存ファイル用） */}
        <div className="mb-6 bg-gray-50 rounded-lg p-4">
          <h3 className="text-sm font-semibold text-gray-700 mb-2">
            または、既存ファイルのパスを入力
          </h3>
          <div className="flex gap-2">
            <input
              type="text"
              value={filePath}
              onChange={(e) => setFilePath(e.target.value)}
              placeholder="/app/uploads/project_15/test.xlsx"
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
