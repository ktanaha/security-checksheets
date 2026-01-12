import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createProject } from '../services/projectService';
import Button from '../components/atoms/Button';
import type { CreateProjectRequest } from '../types/project';

/**
 * 案件作成フォームページ
 */
const ProjectForm: React.FC = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState<CreateProjectRequest>({
    customer_name: '',
    description: '',
    owner: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // フォーム入力の変更ハンドラー
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  // フォーム送信ハンドラー
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const project = await createProject(formData);
      // 作成成功後、案件詳細ページへ遷移
      navigate(`/projects/${project.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : '案件の作成に失敗しました');
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          {/* ヘッダー */}
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-2">
              新規案件作成
            </h1>
            <p className="text-gray-600">
              セキュリティチェックシート抽出の新しい案件を作成します
            </p>
          </div>

          {/* エラー表示 */}
          {error && (
            <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
              {error}
            </div>
          )}

          {/* フォーム */}
          <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow-md p-6">
            {/* 顧客名 */}
            <div className="mb-6">
              <label
                htmlFor="customer_name"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                顧客名 <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="customer_name"
                name="customer_name"
                value={formData.customer_name}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="株式会社〇〇"
              />
            </div>

            {/* 説明 */}
            <div className="mb-6">
              <label
                htmlFor="description"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                案件説明 <span className="text-red-500">*</span>
              </label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleChange}
                required
                rows={4}
                className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="この案件の詳細説明を入力してください"
              />
            </div>

            {/* 担当者 */}
            <div className="mb-6">
              <label
                htmlFor="owner"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                担当者 <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="owner"
                name="owner"
                value={formData.owner}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="山田太郎"
              />
            </div>

            {/* ボタン */}
            <div className="flex gap-4">
              <Button
                type="submit"
                variant="primary"
                disabled={loading}
                className="flex-1"
              >
                {loading ? '作成中...' : '案件を作成'}
              </Button>
              <Button
                type="button"
                variant="secondary"
                onClick={() => navigate('/projects')}
                disabled={loading}
              >
                キャンセル
              </Button>
            </div>
          </form>

          {/* 注意事項 */}
          <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h3 className="font-semibold text-blue-900 mb-2">📝 案件作成後の流れ</h3>
            <ul className="list-disc list-inside text-blue-800 space-y-1 text-sm">
              <li>案件作成後、Excelファイルをアップロードできます</li>
              <li>アップロードしたファイルはプレビュー表示可能です</li>
              <li>プレビュー画面から範囲選択してナレッジを抽出できます</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProjectForm;
