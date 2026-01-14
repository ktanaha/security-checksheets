package domain

import "time"

// Department は部門のドメインモデル
type Department struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// DepartmentRepository は部門リポジトリのインターフェース
type DepartmentRepository interface {
	GetAll() ([]*Department, error)
	GetByID(id int) (*Department, error)
}
