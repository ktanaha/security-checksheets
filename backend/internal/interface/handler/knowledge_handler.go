package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/security-checksheets/backend/internal/domain"
	"github.com/security-checksheets/backend/internal/usecase"
)

// KnowledgeHandler はナレッジに関するHTTPハンドラー
type KnowledgeHandler struct {
	useCase usecase.KnowledgeUseCase
}

// NewKnowledgeHandler は新しいKnowledgeHandlerを生成する
func NewKnowledgeHandler(useCase usecase.KnowledgeUseCase) *KnowledgeHandler {
	return &KnowledgeHandler{useCase: useCase}
}

// CreateKnowledgeRequest はナレッジ作成リクエスト
type CreateKnowledgeRequest struct {
	ProjectID     int    `json:"project_id" binding:"required"`
	FileID        *int   `json:"file_id"`
	SheetName     string `json:"sheet_name"`
	SourceRange   string `json:"source_range"`
	Question      string `json:"question" binding:"required"`
	Answer        string `json:"answer"`
	DepartmentID  *int   `json:"department_id"`
	QuestionGroup string `json:"question_group"`
	CreatedBy     string `json:"created_by"`
}

// BulkCreateKnowledgeRequest は一括作成リクエスト
type BulkCreateKnowledgeRequest struct {
	Items []CreateKnowledgeRequest `json:"items" binding:"required"`
}

// CreateKnowledge はナレッジアイテムを作成する
// @Summary ナレッジ作成
// @Description ナレッジアイテムを作成する
// @Tags knowledge
// @Accept json
// @Produce json
// @Param body body CreateKnowledgeRequest true "ナレッジ作成リクエスト"
// @Success 201 {object} domain.KnowledgeItem
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/knowledge [post]
func (h *KnowledgeHandler) CreateKnowledge(c *gin.Context) {
	var req CreateKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := domain.NewKnowledgeItem(
		req.ProjectID,
		req.FileID,
		req.SheetName,
		req.SourceRange,
		req.Question,
		req.Answer,
		req.DepartmentID,
		req.CreatedBy,
	)

	if req.QuestionGroup != "" {
		item.QuestionGroup = req.QuestionGroup
	}

	if err := h.useCase.CreateKnowledge(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// BulkCreateKnowledge は複数のナレッジアイテムを一括作成する
// @Summary ナレッジ一括作成
// @Description 複数のナレッジアイテムを一括作成する
// @Tags knowledge
// @Accept json
// @Produce json
// @Param body body BulkCreateKnowledgeRequest true "一括作成リクエスト"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/knowledge/bulk [post]
func (h *KnowledgeHandler) BulkCreateKnowledge(c *gin.Context) {
	var req BulkCreateKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]*domain.KnowledgeItem, len(req.Items))
	for i, itemReq := range req.Items {
		items[i] = domain.NewKnowledgeItem(
			itemReq.ProjectID,
			itemReq.FileID,
			itemReq.SheetName,
			itemReq.SourceRange,
			itemReq.Question,
			itemReq.Answer,
			itemReq.DepartmentID,
			itemReq.CreatedBy,
		)
		if itemReq.QuestionGroup != "" {
			items[i].QuestionGroup = itemReq.QuestionGroup
		}
	}

	if err := h.useCase.BulkCreateKnowledge(items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "一括作成が完了しました", "count": len(items)})
}

// GetKnowledge はナレッジアイテムを取得する
// @Summary ナレッジ取得
// @Description 指定されたIDのナレッジアイテムを取得する
// @Tags knowledge
// @Produce json
// @Param id path int true "ナレッジID"
// @Success 200 {object} domain.KnowledgeItem
// @Failure 404 {object} gin.H
// @Router /api/knowledge/{id} [get]
func (h *KnowledgeHandler) GetKnowledge(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	item, err := h.useCase.GetKnowledge(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ナレッジアイテムが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// ListKnowledgeByProject は案件に紐づくナレッジアイテムを取得する
// @Summary 案件のナレッジ一覧取得
// @Description 指定された案件に紐づくナレッジアイテムを取得する
// @Tags knowledge
// @Produce json
// @Param id path int true "案件ID"
// @Success 200 {array} domain.KnowledgeItem
// @Failure 500 {object} gin.H
// @Router /api/projects/{id}/knowledge [get]
func (h *KnowledgeHandler) ListKnowledgeByProject(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な案件IDです"})
		return
	}

	items, err := h.useCase.GetKnowledgeByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// UpdateKnowledge はナレッジアイテムを更新する
// @Summary ナレッジ更新
// @Description ナレッジアイテムを更新する
// @Tags knowledge
// @Accept json
// @Produce json
// @Param id path int true "ナレッジID"
// @Param body body domain.KnowledgeItem true "ナレッジアイテム"
// @Success 200 {object} domain.KnowledgeItem
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/knowledge/{id} [put]
func (h *KnowledgeHandler) UpdateKnowledge(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	var item domain.KnowledgeItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.ID = id
	if err := h.useCase.UpdateKnowledge(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteKnowledge はナレッジアイテムを削除する
// @Summary ナレッジ削除
// @Description ナレッジアイテムを削除する
// @Tags knowledge
// @Param id path int true "ナレッジID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/knowledge/{id} [delete]
func (h *KnowledgeHandler) DeleteKnowledge(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	if err := h.useCase.DeleteKnowledge(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchKnowledge はナレッジアイテムを検索する
// @Summary ナレッジ検索
// @Description ナレッジアイテムを検索する
// @Tags knowledge
// @Produce json
// @Param q query string false "検索クエリ"
// @Param project_id query int false "案件ID"
// @Param department_id query int false "部門ID"
// @Param status query string false "ステータス"
// @Success 200 {array} domain.KnowledgeItem
// @Failure 500 {object} gin.H
// @Router /api/knowledge/search [get]
func (h *KnowledgeHandler) SearchKnowledge(c *gin.Context) {
	query := c.Query("q")
	filters := make(map[string]interface{})

	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		if projectID, err := strconv.Atoi(projectIDStr); err == nil {
			filters["project_id"] = projectID
		}
	}

	if departmentIDStr := c.Query("department_id"); departmentIDStr != "" {
		if departmentID, err := strconv.Atoi(departmentIDStr); err == nil {
			filters["department_id"] = departmentID
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	items, err := h.useCase.SearchKnowledge(query, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}
