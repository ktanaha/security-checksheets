package repository

import (
	"database/sql"
	"time"

	"github.com/security-checksheets/backend/internal/domain"
)

// ProjectRepositoryImpl はProjectRepositoryの実装
type ProjectRepositoryImpl struct {
	db *sql.DB
}

// NewProjectRepository は新しいProjectRepositoryを生成する
func NewProjectRepository(db *sql.DB) domain.ProjectRepository {
	return &ProjectRepositoryImpl{db: db}
}

// Create は新規案件を作成する
func (r *ProjectRepositoryImpl) Create(project *domain.Project) error {
	query := `
		INSERT INTO projects (customer_name, description, owner, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		project.CustomerName,
		project.Description,
		project.Owner,
		project.Status,
		time.Now(),
		time.Now(),
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)

	return err
}

// GetByID は指定されたIDの案件を取得する
func (r *ProjectRepositoryImpl) GetByID(id int) (*domain.Project, error) {
	query := `
		SELECT id, customer_name, description, owner, status, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	project := &domain.Project{}
	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.CustomerName,
		&project.Description,
		&project.Owner,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return project, nil
}

// GetAll はすべての案件を取得する
func (r *ProjectRepositoryImpl) GetAll() ([]*domain.Project, error) {
	query := `
		SELECT id, customer_name, description, owner, status, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []*domain.Project{}
	for rows.Next() {
		project := &domain.Project{}
		err := rows.Scan(
			&project.ID,
			&project.CustomerName,
			&project.Description,
			&project.Owner,
			&project.Status,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// Update は案件情報を更新する
func (r *ProjectRepositoryImpl) Update(project *domain.Project) error {
	query := `
		UPDATE projects
		SET customer_name = $1, description = $2, owner = $3, status = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(
		query,
		project.CustomerName,
		project.Description,
		project.Owner,
		project.Status,
		time.Now(),
		project.ID,
	)

	return err
}

// Delete は案件を削除する
func (r *ProjectRepositoryImpl) Delete(id int) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
