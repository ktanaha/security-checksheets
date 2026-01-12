// 案件管理関連の型定義

/**
 * 案件
 */
export interface Project {
  id: number;
  customer_name: string;
  description: string;
  owner: string;
  status: ProjectStatus;
  created_at: string;
  updated_at: string;
}

/**
 * 案件ステータス
 */
export type ProjectStatus = 'active' | 'in_progress' | 'completed' | 'archived';

/**
 * 案件作成リクエスト
 */
export interface CreateProjectRequest {
  customer_name: string;
  description: string;
  owner: string;
}

/**
 * 案件更新リクエスト
 */
export interface UpdateProjectRequest {
  customer_name: string;
  description: string;
  owner: string;
  status: ProjectStatus;
}

/**
 * アップロード済みファイル
 */
export interface UploadedFile {
  id: number;
  project_id: number;
  file_name: string;
  file_path: string;
  file_size: number;
  uploaded_by: string;
  uploaded_at: string;
}

/**
 * ファイルアップロードリクエスト
 */
export interface UploadFileRequest {
  file: File;
  uploaded_by: string;
}
