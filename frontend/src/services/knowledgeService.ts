import axios from 'axios';
import type {
  KnowledgeItem,
  CreateKnowledgeRequest,
  BulkCreateKnowledgeRequest,
  Department,
  KnowledgeSearchFilters,
} from '../types/knowledge';

const API_BASE_URL = 'http://localhost:8080';

/**
 * ナレッジアイテムを作成
 */
export const createKnowledge = async (
  request: CreateKnowledgeRequest
): Promise<KnowledgeItem> => {
  const response = await axios.post<KnowledgeItem>(
    `${API_BASE_URL}/api/knowledge`,
    request
  );
  return response.data;
};

/**
 * ナレッジアイテムを一括作成
 */
export const bulkCreateKnowledge = async (
  request: BulkCreateKnowledgeRequest
): Promise<{ message: string; count: number }> => {
  const response = await axios.post<{ message: string; count: number }>(
    `${API_BASE_URL}/api/knowledge/bulk`,
    request
  );
  return response.data;
};

/**
 * ナレッジアイテムを取得
 */
export const getKnowledge = async (id: number): Promise<KnowledgeItem> => {
  const response = await axios.get<KnowledgeItem>(
    `${API_BASE_URL}/api/knowledge/${id}`
  );
  return response.data;
};

/**
 * 案件のナレッジ一覧を取得
 */
export const getProjectKnowledge = async (projectId: number): Promise<KnowledgeItem[]> => {
  const response = await axios.get<KnowledgeItem[]>(
    `${API_BASE_URL}/api/projects/${projectId}/knowledge`
  );
  return response.data;
};

/**
 * ナレッジアイテムを更新
 */
export const updateKnowledge = async (
  id: number,
  item: KnowledgeItem
): Promise<KnowledgeItem> => {
  const response = await axios.put<KnowledgeItem>(
    `${API_BASE_URL}/api/knowledge/${id}`,
    item
  );
  return response.data;
};

/**
 * ナレッジアイテムを削除
 */
export const deleteKnowledge = async (id: number): Promise<void> => {
  await axios.delete(`${API_BASE_URL}/api/knowledge/${id}`);
};

/**
 * ナレッジを検索
 */
export const searchKnowledge = async (
  filters: KnowledgeSearchFilters
): Promise<KnowledgeItem[]> => {
  const params = new URLSearchParams();
  if (filters.q) params.append('q', filters.q);
  if (filters.project_id) params.append('project_id', filters.project_id.toString());
  if (filters.department_id)
    params.append('department_id', filters.department_id.toString());
  if (filters.status) params.append('status', filters.status);

  const response = await axios.get<KnowledgeItem[]>(
    `${API_BASE_URL}/api/knowledge/search?${params.toString()}`
  );
  return response.data;
};

/**
 * 部門一覧を取得
 */
export const getDepartments = async (): Promise<Department[]> => {
  const response = await axios.get<Department[]>(`${API_BASE_URL}/api/departments`);
  return response.data;
};
