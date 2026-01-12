package domain

import "time"

// Project は案件を表すドメインエンティティ
type Project struct {
	ID          int       `json:"id"`
	CustomerName string    `json:"customer_name"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectRepository は案件データアクセスのインターフェース
type ProjectRepository interface {
	Create(project *Project) error
	GetByID(id int) (*Project, error)
	GetAll() ([]*Project, error)
	Update(project *Project) error
	Delete(id int) error
}

// NewProject は新規案件を生成する
func NewProject(customerName, description, owner string) *Project {
	return &Project{
		CustomerName: customerName,
		Description:  description,
		Owner:        owner,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Validate は案件データの妥当性を検証する
func (p *Project) Validate() error {
	if p.CustomerName == "" {
		return &ValidationError{Field: "customer_name", Message: "顧客名は必須です"}
	}
	return nil
}

// ValidationError はバリデーションエラーを表す
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
