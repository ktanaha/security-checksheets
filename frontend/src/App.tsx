import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import ExcelPreview from './pages/ExcelPreview';

function Home() {
  const [backendHealth, setBackendHealth] = useState<string>('checking...');
  const [excelHealth, setExcelHealth] = useState<string>('checking...');

  useEffect(() => {
    // Goバックエンドのヘルスチェック
    fetch('http://localhost:8080/health')
      .then((res) => res.json())
      .then((data) => setBackendHealth(data.status))
      .catch(() => setBackendHealth('error'));

    // Python Excel Serviceのヘルスチェック
    fetch('http://localhost:8000/health')
      .then((res) => res.json())
      .then((data) => setExcelHealth(data.status))
      .catch(() => setExcelHealth('error'));
  }, []);

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-4xl font-bold text-center mb-8 text-blue-600">
          セキュリティチェックシート抽出アプリ
        </h1>

        <div className="bg-white rounded-lg shadow-md p-6 max-w-2xl mx-auto">
          <h2 className="text-2xl font-semibold mb-4">システムステータス</h2>

          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 bg-gray-50 rounded">
              <span className="font-medium">フロントエンド (React + Vite)</span>
              <span className="px-3 py-1 bg-green-500 text-white rounded-full text-sm">
                Running
              </span>
            </div>

            <div className="flex items-center justify-between p-4 bg-gray-50 rounded">
              <span className="font-medium">バックエンド API (Go)</span>
              <span
                className={`px-3 py-1 rounded-full text-sm text-white ${
                  backendHealth === 'healthy'
                    ? 'bg-green-500'
                    : backendHealth === 'checking...'
                    ? 'bg-yellow-500'
                    : 'bg-red-500'
                }`}
              >
                {backendHealth}
              </span>
            </div>

            <div className="flex items-center justify-between p-4 bg-gray-50 rounded">
              <span className="font-medium">Excel Service (Python)</span>
              <span
                className={`px-3 py-1 rounded-full text-sm text-white ${
                  excelHealth === 'healthy'
                    ? 'bg-green-500'
                    : excelHealth === 'checking...'
                    ? 'bg-yellow-500'
                    : 'bg-red-500'
                }`}
              >
                {excelHealth}
              </span>
            </div>
          </div>

          <div className="mt-6 p-4 bg-blue-50 rounded">
            <h3 className="font-semibold text-blue-900 mb-2">機能</h3>
            <ul className="list-disc list-inside text-blue-800 space-y-1">
              <li>
                <Link to="/excel-preview" className="underline hover:text-blue-600">
                  Excelプレビュー
                </Link>
              </li>
              <li>案件管理機能（実装済み）</li>
              <li>ファイルアップロード機能（実装済み）</li>
              <li>ナレッジ検索機能（予定）</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/excel-preview" element={<ExcelPreview />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
