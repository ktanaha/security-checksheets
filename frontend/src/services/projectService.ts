import axios from 'axios';
import type {
  Project,
  CreateProjectRequest,
  UpdateProjectRequest,
  UploadedFile,
  UploadFileRequest,
} from '../types/project';

const API_BASE_URL = 'http://localhost:8080';

/**
 * 案件一覧を取得
 */
export const getProjects = async (): Promise<Project[]> => {
  const response = await axios.get<Project[]>(`${API_BASE_URL}/api/projects`);
  return response.data;
};

/**
 * 案件詳細を取得
 */
export const getProject = async (id: number): Promise<Project> => {
  const response = await axios.get<Project>(`${API_BASE_URL}/api/projects/${id}`);
  return response.data;
};

/**
 * 案件を作成
 */
export const createProject = async (request: CreateProjectRequest): Promise<Project> => {
  const response = await axios.post<Project>(`${API_BASE_URL}/api/projects`, request);
  return response.data;
};

/**
 * 案件を更新
 */
export const updateProject = async (
  id: number,
  request: UpdateProjectRequest
): Promise<Project> => {
  const response = await axios.put<Project>(`${API_BASE_URL}/api/projects/${id}`, request);
  return response.data;
};

/**
 * 案件を削除
 */
export const deleteProject = async (id: number): Promise<void> => {
  await axios.delete(`${API_BASE_URL}/api/projects/${id}`);
};

/**
 * 案件のファイル一覧を取得
 */
export const getProjectFiles = async (projectId: number): Promise<UploadedFile[]> => {
  const response = await axios.get<UploadedFile[]>(
    `${API_BASE_URL}/api/projects/${projectId}/files`
  );
  return response.data;
};

/**
 * ファイルをアップロード
 */
export const uploadFile = async (
  projectId: number,
  request: UploadFileRequest
): Promise<UploadedFile> => {
  const formData = new FormData();
  formData.append('file', request.file);
  formData.append('uploaded_by', request.uploaded_by);

  const response = await axios.post<UploadedFile>(
    `${API_BASE_URL}/api/projects/${projectId}/files`,
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }
  );
  return response.data;
};

/**
 * ファイルを削除
 */
export const deleteFile = async (fileId: number): Promise<void> => {
  await axios.delete(`${API_BASE_URL}/api/files/${fileId}`);
};
