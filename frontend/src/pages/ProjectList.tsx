import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getProjects, deleteProject } from '../services/projectService';
import ProjectCard from '../components/molecules/ProjectCard';
import Button from '../components/atoms/Button';
import type { Project } from '../types/project';

/**
 * 案件一覧ページ
 */
const ProjectList: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 案件一覧の取得
  const fetchProjects = async () => {
    setLoading(true);
    setError(null);

    try {
      const data = await getProjects();
      setProjects(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : '案件の取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  // 案件削除
  const handleDelete = async (id: number) => {
    if (!window.confirm('この案件を削除してもよろしいですか？')) {
      return;
    }

    try {
      await deleteProject(id);
      // 削除後、一覧を再取得
      await fetchProjects();
    } catch (err) {
      setError(err instanceof Error ? err.message : '案件の削除に失敗しました');
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto px-4 py-8">
        {/* ヘッダー */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-4xl font-bold text-gray-900 mb-2">案件管理</h1>
            <p className="text-gray-600">
              セキュリティチェックシート抽出アプリ - 案件一覧
            </p>
          </div>
          <div className="flex gap-2">
            <Link to="/">
              <Button variant="secondary">ホームへ戻る</Button>
            </Link>
            <Link to="/projects/new">
              <Button variant="success">+ 新規案件作成</Button>
            </Link>
          </div>
        </div>

        {/* エラー表示 */}
        {error && (
          <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {/* ローディング表示 */}
        {loading && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">読み込み中...</p>
          </div>
        )}

        {/* 案件一覧 */}
        {!loading && projects.length === 0 && (
          <div className="text-center py-12 bg-white rounded-lg shadow">
            <p className="text-gray-500 text-lg mb-4">
              案件がまだありません
            </p>
            <Link to="/projects/new">
              <Button variant="primary">最初の案件を作成</Button>
            </Link>
          </div>
        )}

        {!loading && projects.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {projects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}

        {/* 統計情報 */}
        {!loading && projects.length > 0 && (
          <div className="mt-8 bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold mb-4">統計情報</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="text-center">
                <p className="text-3xl font-bold text-blue-600">{projects.length}</p>
                <p className="text-gray-600">総案件数</p>
              </div>
              <div className="text-center">
                <p className="text-3xl font-bold text-green-600">
                  {projects.filter((p) => p.status === 'active').length}
                </p>
                <p className="text-gray-600">活動中</p>
              </div>
              <div className="text-center">
                <p className="text-3xl font-bold text-yellow-600">
                  {projects.filter((p) => p.status === 'in_progress').length}
                </p>
                <p className="text-gray-600">進行中</p>
              </div>
              <div className="text-center">
                <p className="text-3xl font-bold text-blue-600">
                  {projects.filter((p) => p.status === 'completed').length}
                </p>
                <p className="text-gray-600">完了</p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ProjectList;
