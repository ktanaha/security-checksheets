package usecase

import (
	"errors"
	"testing"

	"github.com/security-checksheets/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectRepository はProjectRepositoryのモック
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(project *domain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(id int) (*domain.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) GetAll() ([]*domain.Project, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(project *domain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestProjectUseCase_CreateProject(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	usecase := NewProjectUseCase(mockRepo)

	project := domain.NewProject("テスト株式会社", "テスト案件", "山田太郎")

	mockRepo.On("Create", project).Return(nil)

	err := usecase.CreateProject(project)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_CreateProject_ValidationError(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	usecase := NewProjectUseCase(mockRepo)

	// 顧客名が空の案件（バリデーションエラー）
	project := &domain.Project{
		CustomerName: "",
		Description:  "テスト案件",
		Owner:        "山田太郎",
	}

	err := usecase.CreateProject(project)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "顧客名は必須です")

	// Createが呼ばれないことを確認
	mockRepo.AssertNotCalled(t, "Create")
}

func TestProjectUseCase_GetProject(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	usecase := NewProjectUseCase(mockRepo)

	expectedProject := &domain.Project{
		ID:           1,
		CustomerName: "テスト株式会社",
		Description:  "テスト案件",
		Owner:        "山田太郎",
		Status:       "active",
	}

	mockRepo.On("GetByID", 1).Return(expectedProject, nil)

	project, err := usecase.GetProject(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedProject, project)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_GetProject_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	usecase := NewProjectUseCase(mockRepo)

	mockRepo.On("GetByID", 999).Return(nil, errors.New("not found"))

	project, err := usecase.GetProject(999)
	assert.Error(t, err)
	assert.Nil(t, project)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_ListProjects(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	usecase := NewProjectUseCase(mockRepo)

	expectedProjects := []*domain.Project{
		{ID: 1, CustomerName: "テスト株式会社1", Status: "active"},
		{ID: 2, CustomerName: "テスト株式会社2", Status: "active"},
	}

	mockRepo.On("GetAll").Return(expectedProjects, nil)

	projects, err := usecase.ListProjects()
	assert.NoError(t, err)
	assert.Len(t, projects, 2)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_UpdateProject(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	usecase := NewProjectUseCase(mockRepo)

	project := &domain.Project{
		ID:           1,
		CustomerName: "更新株式会社",
		Description:  "更新された案件",
		Owner:        "山田太郎",
		Status:       "completed",
	}

	mockRepo.On("Update", project).Return(nil)

	err := usecase.UpdateProject(project)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_UpdateProject_ValidationError(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	usecase := NewProjectUseCase(mockRepo)

	// 顧客名が空の案件（バリデーションエラー）
	project := &domain.Project{
		ID:           1,
		CustomerName: "",
		Description:  "テスト案件",
		Owner:        "山田太郎",
	}

	err := usecase.UpdateProject(project)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "顧客名は必須です")

	// Updateが呼ばれないことを確認
	mockRepo.AssertNotCalled(t, "Update")
}

func TestProjectUseCase_DeleteProject(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	usecase := NewProjectUseCase(mockRepo)

	mockRepo.On("Delete", 1).Return(nil)

	err := usecase.DeleteProject(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
