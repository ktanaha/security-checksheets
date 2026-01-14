package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/security-checksheets/backend/internal/domain"
)

// KnowledgeRepositoryImpl はKnowledgeRepositoryの実装
type KnowledgeRepositoryImpl struct {
	db *sql.DB
}

// NewKnowledgeRepository は新しいKnowledgeRepositoryを生成する
func NewKnowledgeRepository(db *sql.DB) domain.KnowledgeRepository {
	return &KnowledgeRepositoryImpl{db: db}
}

// Create は新規ナレッジアイテムを作成する
func (r *KnowledgeRepositoryImpl) Create(item *domain.KnowledgeItem) error {
	query := `
		INSERT INTO knowledge_items (
			project_id, file_id, sheet_name, source_range, question, answer,
			department_id, question_group, status, version, created_by, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		item.ProjectID,
		item.FileID,
		item.SheetName,
		item.SourceRange,
		item.Question,
		item.Answer,
		item.DepartmentID,
		item.QuestionGroup,
		item.Status,
		item.Version,
		item.CreatedBy,
		time.Now(),
		time.Now(),
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	return err
}

// GetByID は指定されたIDのナレッジアイテムを取得する
func (r *KnowledgeRepositoryImpl) GetByID(id int) (*domain.KnowledgeItem, error) {
	query := `
		SELECT id, project_id, file_id, sheet_name, source_range, question, answer,
		       department_id, question_group, status, version, created_by, created_at, updated_at
		FROM knowledge_items
		WHERE id = $1
	`

	item := &domain.KnowledgeItem{}
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.ProjectID,
		&item.FileID,
		&item.SheetName,
		&item.SourceRange,
		&item.Question,
		&item.Answer,
		&item.DepartmentID,
		&item.QuestionGroup,
		&item.Status,
		&item.Version,
		&item.CreatedBy,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

// GetByProjectID は指定された案件のすべてのナレッジアイテムを取得する
func (r *KnowledgeRepositoryImpl) GetByProjectID(projectID int) ([]*domain.KnowledgeItem, error) {
	query := `
		SELECT id, project_id, file_id, sheet_name, source_range, question, answer,
		       department_id, question_group, status, version, created_by, created_at, updated_at
		FROM knowledge_items
		WHERE project_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*domain.KnowledgeItem{}
	for rows.Next() {
		item := &domain.KnowledgeItem{}
		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.FileID,
			&item.SheetName,
			&item.SourceRange,
			&item.Question,
			&item.Answer,
			&item.DepartmentID,
			&item.QuestionGroup,
			&item.Status,
			&item.Version,
			&item.CreatedBy,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Update はナレッジアイテムを更新する
func (r *KnowledgeRepositoryImpl) Update(item *domain.KnowledgeItem) error {
	query := `
		UPDATE knowledge_items
		SET project_id = $1, file_id = $2, sheet_name = $3, source_range = $4,
		    question = $5, answer = $6, department_id = $7, question_group = $8,
		    status = $9, version = $10, updated_at = $11
		WHERE id = $12
	`

	_, err := r.db.Exec(
		query,
		item.ProjectID,
		item.FileID,
		item.SheetName,
		item.SourceRange,
		item.Question,
		item.Answer,
		item.DepartmentID,
		item.QuestionGroup,
		item.Status,
		item.Version,
		time.Now(),
		item.ID,
	)

	return err
}

// Delete はナレッジアイテムを削除する
func (r *KnowledgeRepositoryImpl) Delete(id int) error {
	query := `DELETE FROM knowledge_items WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Search はナレッジアイテムを検索する
func (r *KnowledgeRepositoryImpl) Search(query string, filters map[string]interface{}) ([]*domain.KnowledgeItem, error) {
	// 基本クエリ
	sql := `
		SELECT id, project_id, file_id, sheet_name, source_range, question, answer,
		       department_id, question_group, status, version, created_by, created_at, updated_at
		FROM knowledge_items
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	// クエリ文字列による検索（質問または回答に含まれる）
	if query != "" {
		sql += fmt.Sprintf(" AND (question ILIKE $%d OR answer ILIKE $%d)", argIndex, argIndex+1)
		searchTerm := "%" + query + "%"
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	// フィルタ条件の追加
	if projectID, ok := filters["project_id"].(int); ok && projectID > 0 {
		sql += fmt.Sprintf(" AND project_id = $%d", argIndex)
		args = append(args, projectID)
		argIndex++
	}

	if departmentID, ok := filters["department_id"].(int); ok && departmentID > 0 {
		sql += fmt.Sprintf(" AND department_id = $%d", argIndex)
		args = append(args, departmentID)
		argIndex++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		sql += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	sql += " ORDER BY created_at DESC"

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*domain.KnowledgeItem{}
	for rows.Next() {
		item := &domain.KnowledgeItem{}
		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.FileID,
			&item.SheetName,
			&item.SourceRange,
			&item.Question,
			&item.Answer,
			&item.DepartmentID,
			&item.QuestionGroup,
			&item.Status,
			&item.Version,
			&item.CreatedBy,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
