package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/security-checksheets/backend/internal/usecase"
)

// FileHandler はファイルに関するHTTPハンドラー
type FileHandler struct {
	useCase usecase.FileUseCase
}

// NewFileHandler は新しいFileHandlerを生成する
func NewFileHandler(useCase usecase.FileUseCase) *FileHandler {
	return &FileHandler{useCase: useCase}
}

// UploadFile はファイルをアップロードする
// @Summary ファイルアップロード
// @Description 案件にファイルをアップロードする
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "案件ID"
// @Param file formData file true "アップロードファイル"
// @Param uploaded_by formData string false "アップロード者"
// @Success 201 {object} domain.UploadedFile
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/projects/{id}/files [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な案件IDです"})
		return
	}

	// マルチパートフォームからファイルを取得
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ファイルが指定されていません"})
		return
	}

	// アップロード者（オプション）
	uploadedBy := c.PostForm("uploaded_by")
	if uploadedBy == "" {
		uploadedBy = "anonymous"
	}

	// ファイルアップロード処理
	file, err := h.useCase.UploadFile(projectID, fileHeader, uploadedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, file)
}

// GetFile は指定されたIDのファイルを取得する
// @Summary ファイル詳細取得
// @Description 指定されたIDのファイル情報を取得する
// @Tags files
// @Produce json
// @Param id path int true "ファイルID"
// @Success 200 {object} domain.UploadedFile
// @Failure 404 {object} gin.H
// @Router /api/files/{id} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なファイルIDです"})
		return
	}

	file, err := h.useCase.GetFile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ファイルが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, file)
}

// ListFilesByProject は指定された案件のすべてのファイルを取得する
// @Summary 案件のファイル一覧取得
// @Description 指定された案件に紐づくすべてのファイルを取得する
// @Tags files
// @Produce json
// @Param id path int true "案件ID"
// @Success 200 {array} domain.UploadedFile
// @Failure 500 {object} gin.H
// @Router /api/projects/{id}/files [get]
func (h *FileHandler) ListFilesByProject(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な案件IDです"})
		return
	}

	files, err := h.useCase.GetFilesByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

// DeleteFile はファイルを削除する
// @Summary ファイル削除
// @Description ファイルを削除する
// @Tags files
// @Param id path int true "ファイルID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なファイルIDです"})
		return
	}

	if err := h.useCase.DeleteFile(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
