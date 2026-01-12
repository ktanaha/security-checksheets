package domain

import "time"

// UploadedFile はアップロードされたファイルを表すドメインエンティティ
type UploadedFile struct {
	ID         int       `json:"id"`
	ProjectID  int       `json:"project_id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	FileSize   int64     `json:"file_size"`
	UploadedBy string    `json:"uploaded_by"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// FileRepository はファイルデータアクセスのインターフェース
type FileRepository interface {
	Create(file *UploadedFile) error
	GetByID(id int) (*UploadedFile, error)
	GetByProjectID(projectID int) ([]*UploadedFile, error)
	Delete(id int) error
}

// NewUploadedFile は新規ファイルエンティティを生成する
func NewUploadedFile(projectID int, fileName, filePath string, fileSize int64, uploadedBy string) *UploadedFile {
	return &UploadedFile{
		ProjectID:  projectID,
		FileName:   fileName,
		FilePath:   filePath,
		FileSize:   fileSize,
		UploadedBy: uploadedBy,
		UploadedAt: time.Now(),
	}
}

// Validate はファイルデータの妥当性を検証する
func (f *UploadedFile) Validate() error {
	if f.ProjectID == 0 {
		return &ValidationError{Field: "project_id", Message: "案件IDは必須です"}
	}
	if f.FileName == "" {
		return &ValidationError{Field: "file_name", Message: "ファイル名は必須です"}
	}
	if f.FilePath == "" {
		return &ValidationError{Field: "file_path", Message: "ファイルパスは必須です"}
	}
	return nil
}
