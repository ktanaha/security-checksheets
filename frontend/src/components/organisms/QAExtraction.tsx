import React, { useState, useEffect } from 'react';
import { extractQA } from '../../services/excelService';
import { getDepartments, bulkCreateKnowledge } from '../../services/knowledgeService';
import Button from '../atoms/Button';
import type { ExtractQARequest, QAItem, Department } from '../../types/knowledge';

interface QAExtractionProps {
  filePath: string;
  sheetName: string;
  projectId: number;
  onSuccess?: () => void;
}

/**
 * Q/A抽出コンポーネント
 */
const QAExtraction: React.FC<QAExtractionProps> = ({
  filePath,
  sheetName,
  projectId,
  onSuccess,
}) => {
  // 抽出設定
  const [startRow, setStartRow] = useState(1);
  const [endRow, setEndRow] = useState(10);
  const [questionColumn, setQuestionColumn] = useState(2);
  const [answerColumn, setAnswerColumn] = useState(3);
  const [departmentColumn, setDepartmentColumn] = useState<number | undefined>(4);
  const [skipHeaderRows, setSkipHeaderRows] = useState(1);

  // 抽出結果
  const [extractedItems, setExtractedItems] = useState<QAItem[]>([]);
  const [sourceRange, setSourceRange] = useState('');

  // 部門一覧
  const [departments, setDepartments] = useState<Department[]>([]);
  const [departmentMap, setDepartmentMap] = useState<Map<string, number>>(new Map());

  // 状態管理
  const [extracting, setExtracting] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [createdBy, setCreatedBy] = useState('');

  // 部門一覧を取得
  useEffect(() => {
    const fetchDepartments = async () => {
      try {
        const data = await getDepartments();
        setDepartments(data);

        // 部門名 -> IDのマップを作成
        const map = new Map<string, number>();
        data.forEach((dept) => {
          map.set(dept.name, dept.id);
        });
        setDepartmentMap(map);
      } catch (err) {
        console.error('部門一覧の取得に失敗しました:', err);
      }
    };
    fetchDepartments();
  }, []);

  // Q/A抽出を実行
  const handleExtract = async () => {
    setExtracting(true);
    setError(null);
    setSuccessMessage(null);

    try {
      const request: ExtractQARequest = {
        file_path: filePath,
        sheet_name: sheetName,
        start_row: startRow,
        end_row: endRow,
        question_column: questionColumn,
        answer_column: answerColumn,
        department_column: departmentColumn,
        skip_header_rows: skipHeaderRows,
      };

      const result = await extractQA(request);
      setExtractedItems(result.items);
      setSourceRange(result.source_range);
      setSuccessMessage(`${result.total_items}件のQ/Aを抽出しました`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Q/Aの抽出に失敗しました');
    } finally {
      setExtracting(false);
    }
  };

  // ナレッジとして保存
  const handleSave = async () => {
    if (extractedItems.length === 0) {
      setError('保存するQ/Aがありません');
      return;
    }

    if (!createdBy.trim()) {
      setError('作成者名を入力してください');
      return;
    }

    setSaving(true);
    setError(null);
    setSuccessMessage(null);

    try {
      const items = extractedItems.map((item) => {
        // 部門名からIDを取得
        let departmentId: number | undefined = undefined;
        if (item.department) {
          departmentId = departmentMap.get(item.department);
        }

        return {
          project_id: projectId,
          sheet_name: sheetName,
          source_range: `${sourceRange}:Row${item.row_number}`,
          question: item.question,
          answer: item.answer || '',
          department_id: departmentId,
          created_by: createdBy.trim(),
        };
      });

      await bulkCreateKnowledge({ items });
      setSuccessMessage(`${items.length}件のナレッジを保存しました`);
      setExtractedItems([]);

      // 成功コールバックを呼び出す
      if (onSuccess) {
        setTimeout(() => onSuccess(), 1500);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ナレッジの保存に失敗しました');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Q/A抽出・ナレッジ化</h2>

      {/* エラー表示 */}
      {error && (
        <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
          {error}
        </div>
      )}

      {/* 成功メッセージ */}
      {successMessage && (
        <div className="mb-4 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
          {successMessage}
        </div>
      )}

      {/* 抽出設定 */}
      <div className="mb-6 p-4 bg-gray-50 rounded-lg">
        <h3 className="text-lg font-semibold mb-4">抽出設定</h3>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              開始行
            </label>
            <input
              type="number"
              value={startRow}
              onChange={(e) => setStartRow(parseInt(e.target.value) || 1)}
              min="1"
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              終了行
            </label>
            <input
              type="number"
              value={endRow}
              onChange={(e) => setEndRow(parseInt(e.target.value) || 1)}
              min="1"
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              質問列
            </label>
            <input
              type="number"
              value={questionColumn}
              onChange={(e) => setQuestionColumn(parseInt(e.target.value) || 1)}
              min="1"
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              回答列
            </label>
            <input
              type="number"
              value={answerColumn}
              onChange={(e) => setAnswerColumn(parseInt(e.target.value) || 1)}
              min="1"
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              部門列（任意）
            </label>
            <input
              type="number"
              value={departmentColumn || ''}
              onChange={(e) =>
                setDepartmentColumn(e.target.value ? parseInt(e.target.value) : undefined)
              }
              min="1"
              placeholder="なし"
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              ヘッダー行数
            </label>
            <input
              type="number"
              value={skipHeaderRows}
              onChange={(e) => setSkipHeaderRows(parseInt(e.target.value) || 0)}
              min="0"
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>
        </div>

        <div className="mt-4">
          <Button
            variant="primary"
            onClick={handleExtract}
            disabled={extracting || !filePath || !sheetName}
          >
            {extracting ? '抽出中...' : 'Q/Aを抽出'}
          </Button>
        </div>
      </div>

      {/* 抽出結果 */}
      {extractedItems.length > 0 && (
        <div className="mb-6">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-semibold">
              抽出結果（{extractedItems.length}件）
            </h3>
            <div className="flex items-center gap-2">
              <input
                type="text"
                value={createdBy}
                onChange={(e) => setCreatedBy(e.target.value)}
                placeholder="作成者名"
                className="px-3 py-2 border border-gray-300 rounded-md"
              />
              <Button
                variant="success"
                onClick={handleSave}
                disabled={saving || !createdBy.trim()}
              >
                {saving ? '保存中...' : 'ナレッジとして保存'}
              </Button>
            </div>
          </div>

          <div className="space-y-3 max-h-96 overflow-y-auto">
            {extractedItems.map((item, index) => (
              <div
                key={index}
                className="p-4 bg-gray-50 rounded-lg border border-gray-200"
              >
                <div className="flex items-start justify-between mb-2">
                  <span className="text-xs text-gray-500">行 {item.row_number}</span>
                  {item.department && (
                    <span className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
                      {item.department}
                    </span>
                  )}
                </div>
                <div className="mb-2">
                  <span className="text-sm font-semibold text-gray-700">Q: </span>
                  <span className="text-sm text-gray-900">{item.question}</span>
                </div>
                {item.answer && (
                  <div>
                    <span className="text-sm font-semibold text-gray-700">A: </span>
                    <span className="text-sm text-gray-600">{item.answer}</span>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* 使い方の説明 */}
      {extractedItems.length === 0 && (
        <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
          <h4 className="font-semibold text-blue-900 mb-2">💡 使い方</h4>
          <ul className="list-disc list-inside text-blue-800 space-y-1 text-sm">
            <li>開始行・終了行: データが含まれる行の範囲を指定</li>
            <li>質問列・回答列: 質問と回答が入力されている列番号を指定</li>
            <li>部門列: 担当部門が記載されている列番号（任意）</li>
            <li>ヘッダー行数: スキップする見出し行の数</li>
          </ul>
        </div>
      )}
    </div>
  );
};

export default QAExtraction;
