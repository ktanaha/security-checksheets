// ナレッジ管理関連の型定義

/**
 * ナレッジアイテム
 */
export interface KnowledgeItem {
  id: number;
  project_id: number;
  file_id?: number;
  sheet_name: string;
  source_range: string;
  question: string;
  answer: string;
  department_id?: number;
  question_group: string;
  status: KnowledgeStatus;
  version: number;
  created_by: string;
  created_at: string;
  updated_at: string;
}

/**
 * ナレッジステータス
 */
export type KnowledgeStatus = 'draft' | 'published' | 'archived';

/**
 * ナレッジ作成リクエスト
 */
export interface CreateKnowledgeRequest {
  project_id: number;
  file_id?: number;
  sheet_name: string;
  source_range: string;
  question: string;
  answer: string;
  department_id?: number;
  question_group?: string;
  created_by: string;
}

/**
 * ナレッジ一括作成リクエスト
 */
export interface BulkCreateKnowledgeRequest {
  items: CreateKnowledgeRequest[];
}

/**
 * 部門
 */
export interface Department {
  id: number;
  name: string;
  display_order: number;
  is_active: boolean;
  created_at: string;
}

/**
 * Q/A抽出リクエスト
 */
export interface ExtractQARequest {
  file_path: string;
  sheet_name: string;
  start_row: number;
  end_row: number;
  question_column: number;
  answer_column: number;
  department_column?: number;
  skip_header_rows: number;
}

/**
 * 抽出されたQ/Aアイテム
 */
export interface QAItem {
  row_number: number;
  question: string;
  answer?: string;
  department?: string;
}

/**
 * Q/A抽出レスポンス
 */
export interface ExtractQAResponse {
  file_path: string;
  sheet_name: string;
  source_range: string;
  items: QAItem[];
  total_items: number;
}

/**
 * ナレッジ検索フィルタ
 */
export interface KnowledgeSearchFilters {
  q?: string;
  project_id?: number;
  department_id?: number;
  status?: KnowledgeStatus;
}
