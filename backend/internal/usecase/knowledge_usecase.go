package usecase

import (
	"fmt"

	"github.com/security-checksheets/backend/internal/domain"
)

// KnowledgeUseCase はナレッジに関するビジネスロジックを提供する
type KnowledgeUseCase interface {
	CreateKnowledge(item *domain.KnowledgeItem) error
	GetKnowledge(id int) (*domain.KnowledgeItem, error)
	GetKnowledgeByProject(projectID int) ([]*domain.KnowledgeItem, error)
	UpdateKnowledge(item *domain.KnowledgeItem) error
	DeleteKnowledge(id int) error
	SearchKnowledge(query string, filters map[string]interface{}) ([]*domain.KnowledgeItem, error)
	BulkCreateKnowledge(items []*domain.KnowledgeItem) error
}

// KnowledgeUseCaseImpl はKnowledgeUseCaseの実装
type KnowledgeUseCaseImpl struct {
	knowledgeRepo domain.KnowledgeRepository
	projectRepo   domain.ProjectRepository
}

// NewKnowledgeUseCase は新しいKnowledgeUseCaseを生成する
func NewKnowledgeUseCase(
	knowledgeRepo domain.KnowledgeRepository,
	projectRepo domain.ProjectRepository,
) KnowledgeUseCase {
	return &KnowledgeUseCaseImpl{
		knowledgeRepo: knowledgeRepo,
		projectRepo:   projectRepo,
	}
}

// CreateKnowledge はナレッジアイテムを作成する
func (u *KnowledgeUseCaseImpl) CreateKnowledge(item *domain.KnowledgeItem) error {
	// 案件の存在確認
	_, err := u.projectRepo.GetByID(item.ProjectID)
	if err != nil {
		return fmt.Errorf("案件が存在しません: %w", err)
	}

	// バリデーション
	if err := item.Validate(); err != nil {
		return err
	}

	// 作成
	return u.knowledgeRepo.Create(item)
}

// GetKnowledge はナレッジアイテムを取得する
func (u *KnowledgeUseCaseImpl) GetKnowledge(id int) (*domain.KnowledgeItem, error) {
	return u.knowledgeRepo.GetByID(id)
}

// GetKnowledgeByProject は案件に紐づくナレッジアイテムを取得する
func (u *KnowledgeUseCaseImpl) GetKnowledgeByProject(projectID int) ([]*domain.KnowledgeItem, error) {
	// 案件の存在確認
	_, err := u.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("案件が存在しません: %w", err)
	}

	return u.knowledgeRepo.GetByProjectID(projectID)
}

// UpdateKnowledge はナレッジアイテムを更新する
func (u *KnowledgeUseCaseImpl) UpdateKnowledge(item *domain.KnowledgeItem) error {
	// 存在確認
	_, err := u.knowledgeRepo.GetByID(item.ID)
	if err != nil {
		return fmt.Errorf("ナレッジアイテムが存在しません: %w", err)
	}

	// バリデーション
	if err := item.Validate(); err != nil {
		return err
	}

	// 更新
	return u.knowledgeRepo.Update(item)
}

// DeleteKnowledge はナレッジアイテムを削除する
func (u *KnowledgeUseCaseImpl) DeleteKnowledge(id int) error {
	// 存在確認
	_, err := u.knowledgeRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("ナレッジアイテムが存在しません: %w", err)
	}

	return u.knowledgeRepo.Delete(id)
}

// SearchKnowledge はナレッジアイテムを検索する
func (u *KnowledgeUseCaseImpl) SearchKnowledge(query string, filters map[string]interface{}) ([]*domain.KnowledgeItem, error) {
	return u.knowledgeRepo.Search(query, filters)
}

// BulkCreateKnowledge は複数のナレッジアイテムを一括作成する
func (u *KnowledgeUseCaseImpl) BulkCreateKnowledge(items []*domain.KnowledgeItem) error {
	for _, item := range items {
		// バリデーション
		if err := item.Validate(); err != nil {
			return fmt.Errorf("バリデーションエラー (項目: %s): %w", item.Question, err)
		}

		// 作成
		if err := u.knowledgeRepo.Create(item); err != nil {
			return fmt.Errorf("作成エラー (項目: %s): %w", item.Question, err)
		}
	}

	return nil
}
