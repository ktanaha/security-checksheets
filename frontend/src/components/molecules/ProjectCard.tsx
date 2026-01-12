import React from 'react';
import { Link } from 'react-router-dom';
import type { Project } from '../../types/project';
import Button from '../atoms/Button';

interface ProjectCardProps {
  project: Project;
  onDelete?: (id: number) => void;
}

/**
 * 案件カードコンポーネント
 */
const ProjectCard: React.FC<ProjectCardProps> = ({ project, onDelete }) => {
  const statusColors = {
    active: 'bg-green-100 text-green-800',
    in_progress: 'bg-yellow-100 text-yellow-800',
    completed: 'bg-blue-100 text-blue-800',
    archived: 'bg-gray-100 text-gray-800',
  };

  const statusLabels = {
    active: '活動中',
    in_progress: '進行中',
    completed: '完了',
    archived: 'アーカイブ',
  };

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    });
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-xl font-semibold text-gray-900">{project.customer_name}</h3>
        <span
          className={`px-3 py-1 rounded-full text-sm font-medium ${
            statusColors[project.status]
          }`}
        >
          {statusLabels[project.status]}
        </span>
      </div>

      <p className="text-gray-600 mb-4 line-clamp-2">{project.description}</p>

      <div className="flex items-center text-sm text-gray-500 mb-4">
        <span className="mr-4">
          <strong>担当:</strong> {project.owner}
        </span>
        <span>
          <strong>作成日:</strong> {formatDate(project.created_at)}
        </span>
      </div>

      <div className="flex gap-2">
        <Link to={`/projects/${project.id}`} className="flex-1">
          <Button variant="primary" size="sm" className="w-full">
            詳細を見る
          </Button>
        </Link>
        {onDelete && (
          <Button
            variant="danger"
            size="sm"
            onClick={() => onDelete(project.id)}
          >
            削除
          </Button>
        )}
      </div>
    </div>
  );
};

export default ProjectCard;
