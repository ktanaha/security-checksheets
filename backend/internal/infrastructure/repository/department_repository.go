package repository

import (
	"database/sql"

	"github.com/security-checksheets/backend/internal/domain"
)

// DepartmentRepositoryImpl はDepartmentRepositoryの実装
type DepartmentRepositoryImpl struct {
	db *sql.DB
}

// NewDepartmentRepository は新しいDepartmentRepositoryを生成する
func NewDepartmentRepository(db *sql.DB) domain.DepartmentRepository {
	return &DepartmentRepositoryImpl{db: db}
}

// GetAll はすべての部門を取得する
func (r *DepartmentRepositoryImpl) GetAll() ([]*domain.Department, error) {
	query := `
		SELECT id, name, display_order, is_active, created_at
		FROM departments
		WHERE is_active = true
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	departments := []*domain.Department{}
	for rows.Next() {
		dept := &domain.Department{}
		err := rows.Scan(
			&dept.ID,
			&dept.Name,
			&dept.DisplayOrder,
			&dept.IsActive,
			&dept.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		departments = append(departments, dept)
	}

	return departments, nil
}

// GetByID は指定されたIDの部門を取得する
func (r *DepartmentRepositoryImpl) GetByID(id int) (*domain.Department, error) {
	query := `
		SELECT id, name, display_order, is_active, created_at
		FROM departments
		WHERE id = $1
	`

	dept := &domain.Department{}
	err := r.db.QueryRow(query, id).Scan(
		&dept.ID,
		&dept.Name,
		&dept.DisplayOrder,
		&dept.IsActive,
		&dept.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return dept, nil
}
