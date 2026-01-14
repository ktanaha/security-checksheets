package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/security-checksheets/backend/internal/domain"
)

// DepartmentHandler は部門に関するHTTPハンドラー
type DepartmentHandler struct {
	repo domain.DepartmentRepository
}

// NewDepartmentHandler は新しいDepartmentHandlerを生成する
func NewDepartmentHandler(repo domain.DepartmentRepository) *DepartmentHandler {
	return &DepartmentHandler{repo: repo}
}

// ListDepartments はすべての部門を取得する
// @Summary 部門一覧取得
// @Description すべての有効な部門を取得する
// @Tags departments
// @Produce json
// @Success 200 {array} domain.Department
// @Failure 500 {object} gin.H
// @Router /api/departments [get]
func (h *DepartmentHandler) ListDepartments(c *gin.Context) {
	departments, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, departments)
}
