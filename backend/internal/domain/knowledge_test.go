package domain

import (
	"testing"
)

func TestNewKnowledgeItem(t *testing.T) {
	// 正常ケース
	projectID := 1
	fileID := 2
	sheetName := "セキュリティチェック"
	sourceRange := "A1:D10"
	question := "アクセス制御は適切に設定されていますか？"
	answer := "はい、全ユーザーに対して多要素認証を実施しています"
	departmentID := 1
	createdBy := "テストユーザー"

	item := NewKnowledgeItem(
		projectID,
		&fileID,
		sheetName,
		sourceRange,
		question,
		answer,
		&departmentID,
		createdBy,
	)

	if item.ProjectID != projectID {
		t.Errorf("ProjectID = %d, want %d", item.ProjectID, projectID)
	}
	if *item.FileID != fileID {
		t.Errorf("FileID = %d, want %d", *item.FileID, fileID)
	}
	if item.SheetName != sheetName {
		t.Errorf("SheetName = %s, want %s", item.SheetName, sheetName)
	}
	if item.SourceRange != sourceRange {
		t.Errorf("SourceRange = %s, want %s", item.SourceRange, sourceRange)
	}
	if item.Question != question {
		t.Errorf("Question = %s, want %s", item.Question, question)
	}
	if item.Answer != answer {
		t.Errorf("Answer = %s, want %s", item.Answer, answer)
	}
	if *item.DepartmentID != departmentID {
		t.Errorf("DepartmentID = %d, want %d", *item.DepartmentID, departmentID)
	}
	if item.Status != "draft" {
		t.Errorf("Status = %s, want draft", item.Status)
	}
	if item.Version != 1 {
		t.Errorf("Version = %d, want 1", item.Version)
	}
	if item.CreatedBy != createdBy {
		t.Errorf("CreatedBy = %s, want %s", item.CreatedBy, createdBy)
	}
}

func TestKnowledgeItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    *KnowledgeItem
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常ケース",
			item: &KnowledgeItem{
				ProjectID: 1,
				Question:  "テスト質問",
				Status:    "draft",
			},
			wantErr: false,
		},
		{
			name: "ProjectIDが0",
			item: &KnowledgeItem{
				ProjectID: 0,
				Question:  "テスト質問",
			},
			wantErr: true,
			errMsg:  "project_idは必須です",
		},
		{
			name: "質問が空",
			item: &KnowledgeItem{
				ProjectID: 1,
				Question:  "",
			},
			wantErr: true,
			errMsg:  "質問は必須です",
		},
		{
			name: "質問が長すぎる",
			item: &KnowledgeItem{
				ProjectID: 1,
				Question:  string(make([]byte, 10001)),
			},
			wantErr: true,
			errMsg:  "質問は10000文字以内で入力してください",
		},
		{
			name: "回答が長すぎる",
			item: &KnowledgeItem{
				ProjectID: 1,
				Question:  "テスト質問",
				Answer:    string(make([]byte, 50001)),
			},
			wantErr: true,
			errMsg:  "回答は50000文字以内で入力してください",
		},
		{
			name: "無効なステータス",
			item: &KnowledgeItem{
				ProjectID: 1,
				Question:  "テスト質問",
				Status:    "invalid",
			},
			wantErr: true,
			errMsg:  "ステータスはdraft, published, archivedのいずれかである必要があります",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestKnowledgeItem_UpdateQuestion(t *testing.T) {
	item := &KnowledgeItem{
		ProjectID: 1,
		Question:  "古い質問",
		Version:   1,
	}

	newQuestion := "新しい質問"
	item.UpdateQuestion(newQuestion)

	if item.Question != newQuestion {
		t.Errorf("Question = %s, want %s", item.Question, newQuestion)
	}
	if item.Version != 2 {
		t.Errorf("Version = %d, want 2", item.Version)
	}
}

func TestKnowledgeItem_UpdateAnswer(t *testing.T) {
	item := &KnowledgeItem{
		ProjectID: 1,
		Question:  "質問",
		Answer:    "古い回答",
		Version:   1,
	}

	newAnswer := "新しい回答"
	item.UpdateAnswer(newAnswer)

	if item.Answer != newAnswer {
		t.Errorf("Answer = %s, want %s", item.Answer, newAnswer)
	}
	if item.Version != 2 {
		t.Errorf("Version = %d, want 2", item.Version)
	}
}

func TestKnowledgeItem_Publish(t *testing.T) {
	item := &KnowledgeItem{
		ProjectID: 1,
		Question:  "質問",
		Status:    "draft",
	}

	item.Publish()

	if item.Status != "published" {
		t.Errorf("Status = %s, want published", item.Status)
	}
}

func TestKnowledgeItem_Archive(t *testing.T) {
	item := &KnowledgeItem{
		ProjectID: 1,
		Question:  "質問",
		Status:    "published",
	}

	item.Archive()

	if item.Status != "archived" {
		t.Errorf("Status = %s, want archived", item.Status)
	}
}
