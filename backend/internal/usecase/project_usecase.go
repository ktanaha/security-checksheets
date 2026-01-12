package usecase

import "github.com/security-checksheets/backend/internal/domain"

// ProjectUseCase は案件に関するビジネスロジックを提供する
type ProjectUseCase interface {
	CreateProject(project *domain.Project) error
	GetProject(id int) (*domain.Project, error)
	ListProjects() ([]*domain.Project, error)
	UpdateProject(project *domain.Project) error
	DeleteProject(id int) error
}

// ProjectUseCaseImpl はProjectUseCaseの実装
type ProjectUseCaseImpl struct {
	repo domain.ProjectRepository
}

// NewProjectUseCase は新しいProjectUseCaseを生成する
func NewProjectUseCase(repo domain.ProjectRepository) ProjectUseCase {
	return &ProjectUseCaseImpl{repo: repo}
}

// CreateProject は新規案件を作成する
func (u *ProjectUseCaseImpl) CreateProject(project *domain.Project) error {
	// バリデーション
	if err := project.Validate(); err != nil {
		return err
	}

	// リポジトリに保存
	return u.repo.Create(project)
}

// GetProject は指定されたIDの案件を取得する
func (u *ProjectUseCaseImpl) GetProject(id int) (*domain.Project, error) {
	return u.repo.GetByID(id)
}

// ListProjects はすべての案件を取得する
func (u *ProjectUseCaseImpl) ListProjects() ([]*domain.Project, error) {
	return u.repo.GetAll()
}

// UpdateProject は案件情報を更新する
func (u *ProjectUseCaseImpl) UpdateProject(project *domain.Project) error {
	// バリデーション
	if err := project.Validate(); err != nil {
		return err
	}

	// リポジトリで更新
	return u.repo.Update(project)
}

// DeleteProject は案件を削除する
func (u *ProjectUseCaseImpl) DeleteProject(id int) error {
	return u.repo.Delete(id)
}
