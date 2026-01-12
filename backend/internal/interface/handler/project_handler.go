package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/security-checksheets/backend/internal/domain"
	"github.com/security-checksheets/backend/internal/usecase"
)

// ProjectHandler は案件に関するHTTPハンドラー
type ProjectHandler struct {
	useCase usecase.ProjectUseCase
}

// NewProjectHandler は新しいProjectHandlerを生成する
func NewProjectHandler(useCase usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{useCase: useCase}
}

// CreateProjectRequest は案件作成リクエスト
type CreateProjectRequest struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Description  string `json:"description"`
	Owner        string `json:"owner"`
}

// UpdateProjectRequest は案件更新リクエスト
type UpdateProjectRequest struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Description  string `json:"description"`
	Owner        string `json:"owner"`
	Status       string `json:"status"`
}

// CreateProject は新規案件を作成する
// @Summary 案件作成
// @Description 新規案件を作成する
// @Tags projects
// @Accept json
// @Produce json
// @Param project body CreateProjectRequest true "案件情報"
// @Success 201 {object} domain.Project
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	project := domain.NewProject(req.CustomerName, req.Description, req.Owner)

	if err := h.useCase.CreateProject(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetProject は指定されたIDの案件を取得する
// @Summary 案件取得
// @Description 指定されたIDの案件を取得する
// @Tags projects
// @Produce json
// @Param id path int true "案件ID"
// @Success 200 {object} domain.Project
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	project, err := h.useCase.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "案件が見つかりません"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// ListProjects はすべての案件を取得する
// @Summary 案件一覧取得
// @Description すべての案件を取得する
// @Tags projects
// @Produce json
// @Success 200 {array} domain.Project
// @Failure 500 {object} gin.H
// @Router /api/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	projects, err := h.useCase.ListProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// UpdateProject は案件情報を更新する
// @Summary 案件更新
// @Description 案件情報を更新する
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "案件ID"
// @Param project body UpdateProjectRequest true "案件情報"
// @Success 200 {object} domain.Project
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	project := &domain.Project{
		ID:           id,
		CustomerName: req.CustomerName,
		Description:  req.Description,
		Owner:        req.Owner,
		Status:       req.Status,
	}

	if err := h.useCase.UpdateProject(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject は案件を削除する
// @Summary 案件削除
// @Description 案件を削除する
// @Tags projects
// @Param id path int true "案件ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	if err := h.useCase.DeleteProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
