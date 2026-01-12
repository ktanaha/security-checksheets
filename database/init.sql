-- セキュリティチェックシート抽出アプリ データベース初期化スクリプト

-- 日本語全文検索用の拡張機能を有効化
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- タイムゾーン設定
SET timezone = 'Asia/Tokyo';

-- projects（案件）テーブル
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    description TEXT,
    owner VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- プロジェクトテーブルにインデックス
CREATE INDEX idx_projects_customer ON projects(customer_name);
CREATE INDEX idx_projects_status ON projects(status);

-- uploaded_files（アップロードファイル）テーブル
CREATE TABLE uploaded_files (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT,
    uploaded_by VARCHAR(255),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, file_name)
);

-- ファイルテーブルにインデックス
CREATE INDEX idx_files_project ON uploaded_files(project_id);

-- departments（部門マスタ）テーブル
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    display_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 部門マスタ初期データ
INSERT INTO departments (name, display_order) VALUES
    ('情報システム', 1),
    ('セキュリティ', 2),
    ('法務', 3),
    ('開発', 4),
    ('CS', 5),
    ('営業', 6);

-- knowledge_items（ナレッジQ/A）テーブル
CREATE TABLE knowledge_items (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    file_id INTEGER REFERENCES uploaded_files(id) ON DELETE SET NULL,
    sheet_name VARCHAR(255),
    source_range VARCHAR(100),
    question TEXT NOT NULL,
    answer TEXT,
    department_id INTEGER REFERENCES departments(id),
    question_group VARCHAR(100),
    status VARCHAR(50) DEFAULT 'draft',
    version INTEGER DEFAULT 1,
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ナレッジテーブルにインデックス
CREATE INDEX idx_knowledge_project ON knowledge_items(project_id);
CREATE INDEX idx_knowledge_department ON knowledge_items(department_id);
CREATE INDEX idx_knowledge_status ON knowledge_items(status);

-- 日本語全文検索用インデックス（pg_trgmを使用）
CREATE INDEX idx_knowledge_question_trgm ON knowledge_items USING gin (question gin_trgm_ops);
CREATE INDEX idx_knowledge_answer_trgm ON knowledge_items USING gin (answer gin_trgm_ops);

-- extraction_sessions（抽出セッション・将来用）テーブル
CREATE TABLE extraction_sessions (
    id SERIAL PRIMARY KEY,
    file_id INTEGER NOT NULL REFERENCES uploaded_files(id) ON DELETE CASCADE,
    sheet_name VARCHAR(255),
    selected_range VARCHAR(100),
    excluded_ranges TEXT[],
    settings JSONB,
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 抽出セッションテーブルにインデックス
CREATE INDEX idx_extraction_file ON extraction_sessions(file_id);

-- updated_atを自動更新するトリガー関数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- projectsテーブルのupdated_atトリガー
CREATE TRIGGER update_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- knowledge_itemsテーブルのupdated_atトリガー
CREATE TRIGGER update_knowledge_items_updated_at
    BEFORE UPDATE ON knowledge_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 初期化完了ログ
DO $$
BEGIN
    RAISE NOTICE 'データベース初期化が完了しました';
    RAISE NOTICE 'テーブル数: %', (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public');
END $$;
