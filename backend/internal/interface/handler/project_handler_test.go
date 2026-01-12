package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/security-checksheets/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectUseCase はProjectUseCaseのモック
type MockProjectUseCase struct {
	mock.Mock
}

func (m *MockProjectUseCase) CreateProject(project *domain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectUseCase) GetProject(id int) (*domain.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectUseCase) ListProjects() ([]*domain.Project, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectUseCase) UpdateProject(project *domain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectUseCase) DeleteProject(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestProjectHandler_CreateProject(t *testing.T) {
	mockUseCase := new(MockProjectUseCase)
	handler := NewProjectHandler(mockUseCase)

	router := setupRouter()
	router.POST("/api/projects", handler.CreateProject)

	reqBody := map[string]interface{}{
		"customer_name": "テスト株式会社",
		"description":   "テスト案件",
		"owner":         "山田太郎",
	}
	body, _ := json.Marshal(reqBody)

	mockUseCase.On("CreateProject", mock.AnythingOfType("*domain.Project")).Return(nil)

	req, _ := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestProjectHandler_CreateProject_InvalidJSON(t *testing.T) {
	mockUseCase := new(MockProjectUseCase)
	handler := NewProjectHandler(mockUseCase)

	router := setupRouter()
	router.POST("/api/projects", handler.CreateProject)

	req, _ := http.NewRequest("POST", "/api/projects", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectHandler_GetProject(t *testing.T) {
	mockUseCase := new(MockProjectUseCase)
	handler := NewProjectHandler(mockUseCase)

	router := setupRouter()
	router.GET("/api/projects/:id", handler.GetProject)

	expectedProject := &domain.Project{
		ID:           1,
		CustomerName: "テスト株式会社",
		Description:  "テスト案件",
		Owner:        "山田太郎",
		Status:       "active",
	}

	mockUseCase.On("GetProject", 1).Return(expectedProject, nil)

	req, _ := http.NewRequest("GET", "/api/projects/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Project
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedProject.CustomerName, response.CustomerName)

	mockUseCase.AssertExpectations(t)
}

func TestProjectHandler_GetProject_NotFound(t *testing.T) {
	mockUseCase := new(MockProjectUseCase)
	handler := NewProjectHandler(mockUseCase)

	router := setupRouter()
	router.GET("/api/projects/:id", handler.GetProject)

	mockUseCase.On("GetProject", 999).Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/api/projects/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestProjectHandler_ListProjects(t *testing.T) {
	mockUseCase := new(MockProjectUseCase)
	handler := NewProjectHandler(mockUseCase)

	router := setupRouter()
	router.GET("/api/projects", handler.ListProjects)

	expectedProjects := []*domain.Project{
		{ID: 1, CustomerName: "テスト株式会社1", Status: "active"},
		{ID: 2, CustomerName: "テスト株式会社2", Status: "active"},
	}

	mockUseCase.On("ListProjects").Return(expectedProjects, nil)

	req, _ := http.NewRequest("GET", "/api/projects", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*domain.Project
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 2)

	mockUseCase.AssertExpectations(t)
}

func TestProjectHandler_UpdateProject(t *testing.T) {
	mockUseCase := new(MockProjectUseCase)
	handler := NewProjectHandler(mockUseCase)

	router := setupRouter()
	router.PUT("/api/projects/:id", handler.UpdateProject)

	reqBody := map[string]interface{}{
		"customer_name": "更新株式会社",
		"description":   "更新された案件",
		"owner":         "山田太郎",
		"status":        "completed",
	}
	body, _ := json.Marshal(reqBody)

	mockUseCase.On("UpdateProject", mock.AnythingOfType("*domain.Project")).Return(nil)

	req, _ := http.NewRequest("PUT", "/api/projects/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestProjectHandler_DeleteProject(t *testing.T) {
	mockUseCase := new(MockProjectUseCase)
	handler := NewProjectHandler(mockUseCase)

	router := setupRouter()
	router.DELETE("/api/projects/:id", handler.DeleteProject)

	mockUseCase.On("DeleteProject", 1).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/projects/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockUseCase.AssertExpectations(t)
}
