import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import {
  getProject,
  getProjectFiles,
  uploadFile,
  deleteFile,
} from '../services/projectService';
import { parseExcelFile } from '../services/excelService';
import Button from '../components/atoms/Button';
import type { Project, UploadedFile } from '../types/project';

/**
 * 案件詳細ページ
 */
const ProjectDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const projectId = parseInt(id || '0');

  const [project, setProject] = useState<Project | null>(null);
  const [files, setFiles] = useState<UploadedFile[]>([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [uploadedBy, setUploadedBy] = useState('');

  // 案件詳細とファイル一覧の取得
  const fetchData = async () => {
    setLoading(true);
    setError(null);

    try {
      const [projectData, filesData] = await Promise.all([
        getProject(projectId),
        getProjectFiles(projectId),
      ]);
      setProject(projectData);
      setFiles(filesData);
      setUploadedBy(projectData.owner); // デフォルトで担当者名をセット
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [projectId]);

  // ファイルアップロード
  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
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

    setUploading(true);
    setError(null);

    try {
      await uploadFile(projectId, {
        file,
        uploaded_by: uploadedBy || project?.owner || '不明',
      });
      // アップロード後、ファイル一覧を再取得
      const filesData = await getProjectFiles(projectId);
      setFiles(filesData);
      e.target.value = ''; // ファイル選択をリセット
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ファイルのアップロードに失敗しました');
    } finally {
      setUploading(false);
    }
  };

  // ファイル削除
  const handleDeleteFile = async (fileId: number) => {
    if (!window.confirm('このファイルを削除してもよろしいですか？')) {
      return;
    }

    try {
      await deleteFile(fileId);
      // 削除後、ファイル一覧を再取得
      const filesData = await getProjectFiles(projectId);
      setFiles(filesData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ファイルの削除に失敗しました');
    }
  };

  // ファイルサイズのフォーマット
  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  // 日時のフォーマット
  const formatDateTime = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleString('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const statusLabels = {
    active: '活動中',
    in_progress: '進行中',
    completed: '完了',
    archived: 'アーカイブ',
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">読み込み中...</p>
        </div>
      </div>
    );
  }

  if (!project) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600 mb-4">案件が見つかりません</p>
          <Button onClick={() => navigate('/projects')}>一覧へ戻る</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto px-4 py-8">
        {/* ヘッダー */}
        <div className="mb-8">
          <Button variant="secondary" onClick={() => navigate('/projects')}>
            ← 一覧へ戻る
          </Button>
        </div>

        {/* エラー表示 */}
        {error && (
          <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {/* 案件情報 */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex justify-between items-start mb-4">
            <h1 className="text-3xl font-bold text-gray-900">{project.customer_name}</h1>
            <span className="px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
              {statusLabels[project.status]}
            </span>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <p className="text-sm text-gray-500">担当者</p>
              <p className="text-lg text-gray-900">{project.owner}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">作成日</p>
              <p className="text-lg text-gray-900">{formatDateTime(project.created_at)}</p>
            </div>
          </div>

          <div className="mb-4">
            <p className="text-sm text-gray-500 mb-2">案件説明</p>
            <p className="text-gray-900">{project.description}</p>
          </div>
        </div>

        {/* ファイルアップロードセクション */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">ファイル管理</h2>

          {/* アップロードフォーム */}
          <div className="mb-6 p-4 bg-gray-50 rounded-lg">
            <div className="flex items-center gap-4">
              <div className="flex-1">
                <label
                  htmlFor="file-upload"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Excelファイルをアップロード
                </label>
                <input
                  id="file-upload"
                  type="file"
                  accept=".xlsx,.xls"
                  onChange={handleFileUpload}
                  disabled={uploading}
                  className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 disabled:opacity-50"
                />
              </div>
              <div className="w-48">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  アップロード者
                </label>
                <input
                  type="text"
                  value={uploadedBy}
                  onChange={(e) => setUploadedBy(e.target.value)}
                  disabled={uploading}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm disabled:opacity-50"
                />
              </div>
            </div>
            {uploading && (
              <p className="mt-2 text-sm text-blue-600">アップロード中...</p>
            )}
          </div>

          {/* ファイル一覧 */}
          {files.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <p>まだファイルがアップロードされていません</p>
              <p className="text-sm mt-2">上記からExcelファイルをアップロードしてください</p>
            </div>
          ) : (
            <div className="space-y-3">
              {files.map((file) => (
                <div
                  key={file.id}
                  className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <div className="flex-1">
                    <p className="font-medium text-gray-900">{file.file_name}</p>
                    <p className="text-sm text-gray-500">
                      {formatFileSize(file.file_size)} • {file.uploaded_by} •{' '}
                      {formatDateTime(file.uploaded_at)}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    <Link
                      to={`/excel-preview?filePath=${encodeURIComponent(file.file_path)}`}
                    >
                      <Button variant="primary" size="sm">
                        プレビュー
                      </Button>
                    </Link>
                    <Button
                      variant="danger"
                      size="sm"
                      onClick={() => handleDeleteFile(file.id)}
                    >
                      削除
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProjectDetail;
