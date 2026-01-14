package domain

import (
	"errors"
	"time"
)

// KnowledgeItem はナレッジQ/Aのドメインモデル
type KnowledgeItem struct {
	ID            int       `json:"id"`
	ProjectID     int       `json:"project_id"`
	FileID        *int      `json:"file_id,omitempty"`
	SheetName     string    `json:"sheet_name"`
	SourceRange   string    `json:"source_range"`
	Question      string    `json:"question"`
	Answer        string    `json:"answer"`
	DepartmentID  *int      `json:"department_id,omitempty"`
	QuestionGroup string    `json:"question_group"`
	Status        string    `json:"status"`
	Version       int       `json:"version"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// KnowledgeRepository はナレッジリポジトリのインターフェース
type KnowledgeRepository interface {
	Create(item *KnowledgeItem) error
	GetByID(id int) (*KnowledgeItem, error)
	GetByProjectID(projectID int) ([]*KnowledgeItem, error)
	Update(item *KnowledgeItem) error
	Delete(id int) error
	Search(query string, filters map[string]interface{}) ([]*KnowledgeItem, error)
}

// NewKnowledgeItem は新しいナレッジアイテムを生成する
func NewKnowledgeItem(
	projectID int,
	fileID *int,
	sheetName string,
	sourceRange string,
	question string,
	answer string,
	departmentID *int,
	createdBy string,
) *KnowledgeItem {
	now := time.Now()
	return &KnowledgeItem{
		ProjectID:    projectID,
		FileID:       fileID,
		SheetName:    sheetName,
		SourceRange:  sourceRange,
		Question:     question,
		Answer:       answer,
		DepartmentID: departmentID,
		Status:       "draft",
		Version:      1,
		CreatedBy:    createdBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Validate はナレッジアイテムのバリデーションを行う
func (k *KnowledgeItem) Validate() error {
	if k.ProjectID == 0 {
		return errors.New("project_idは必須です")
	}

	if k.Question == "" {
		return errors.New("質問は必須です")
	}

	if len(k.Question) > 10000 {
		return errors.New("質問は10000文字以内で入力してください")
	}

	if len(k.Answer) > 50000 {
		return errors.New("回答は50000文字以内で入力してください")
	}

	// ステータスの検証
	validStatuses := map[string]bool{
		"draft":     true,
		"published": true,
		"archived":  true,
	}
	if k.Status != "" && !validStatuses[k.Status] {
		return errors.New("ステータスはdraft, published, archivedのいずれかである必要があります")
	}

	return nil
}

// UpdateQuestion は質問を更新する
func (k *KnowledgeItem) UpdateQuestion(question string) {
	k.Question = question
	k.Version++
	k.UpdatedAt = time.Now()
}

// UpdateAnswer は回答を更新する
func (k *KnowledgeItem) UpdateAnswer(answer string) {
	k.Answer = answer
	k.Version++
	k.UpdatedAt = time.Now()
}

// Publish はナレッジアイテムを公開する
func (k *KnowledgeItem) Publish() {
	k.Status = "published"
	k.UpdatedAt = time.Now()
}

// Archive はナレッジアイテムをアーカイブする
func (k *KnowledgeItem) Archive() {
	k.Status = "archived"
	k.UpdatedAt = time.Now()
}
