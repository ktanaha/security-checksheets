package usecase

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/security-checksheets/backend/internal/domain"
)

// FileUseCase はファイルに関するビジネスロジックを提供する
type FileUseCase interface {
	UploadFile(projectID int, fileHeader *multipart.FileHeader, uploadedBy string) (*domain.UploadedFile, error)
	GetFile(id int) (*domain.UploadedFile, error)
	GetFilesByProject(projectID int) ([]*domain.UploadedFile, error)
	DeleteFile(id int) error
}

// FileUseCaseImpl はFileUseCaseの実装
type FileUseCaseImpl struct {
	fileRepo       domain.FileRepository
	projectRepo    domain.ProjectRepository
	uploadBasePath string
}

// NewFileUseCase は新しいFileUseCaseを生成する
func NewFileUseCase(fileRepo domain.FileRepository, projectRepo domain.ProjectRepository, uploadBasePath string) FileUseCase {
	return &FileUseCaseImpl{
		fileRepo:       fileRepo,
		projectRepo:    projectRepo,
		uploadBasePath: uploadBasePath,
	}
}

// UploadFile はファイルをアップロードする
func (u *FileUseCaseImpl) UploadFile(projectID int, fileHeader *multipart.FileHeader, uploadedBy string) (*domain.UploadedFile, error) {
	// 案件の存在確認
	_, err := u.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("案件が存在しません: %w", err)
	}

	// アップロードディレクトリの作成
	uploadDir := filepath.Join(u.uploadBasePath, fmt.Sprintf("project_%d", projectID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("アップロードディレクトリの作成に失敗しました: %w", err)
	}

	// ファイル名の重複を避けるためにタイムスタンプを付与
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s", timestamp, fileHeader.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	// ファイルの保存
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("ファイルのオープンに失敗しました: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの作成に失敗しました: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("ファイルのコピーに失敗しました: %w", err)
	}

	// ファイル情報をDBに保存
	file := domain.NewUploadedFile(
		projectID,
		fileHeader.Filename,
		filePath,
		fileHeader.Size,
		uploadedBy,
	)

	if err := file.Validate(); err != nil {
		// バリデーションエラーの場合はファイルを削除
		os.Remove(filePath)
		return nil, err
	}

	if err := u.fileRepo.Create(file); err != nil {
		// DB保存エラーの場合はファイルを削除
		os.Remove(filePath)
		return nil, fmt.Errorf("ファイル情報の保存に失敗しました: %w", err)
	}

	return file, nil
}

// GetFile は指定されたIDのファイルを取得する
func (u *FileUseCaseImpl) GetFile(id int) (*domain.UploadedFile, error) {
	return u.fileRepo.GetByID(id)
}

// GetFilesByProject は指定された案件のすべてのファイルを取得する
func (u *FileUseCaseImpl) GetFilesByProject(projectID int) ([]*domain.UploadedFile, error) {
	return u.fileRepo.GetByProjectID(projectID)
}

// DeleteFile はファイルを削除する
func (u *FileUseCaseImpl) DeleteFile(id int) error {
	// ファイル情報を取得
	file, err := u.fileRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("ファイルが見つかりません: %w", err)
	}

	// DBから削除
	if err := u.fileRepo.Delete(id); err != nil {
		return fmt.Errorf("ファイル情報の削除に失敗しました: %w", err)
	}

	// 物理ファイルを削除
	if err := os.Remove(file.FilePath); err != nil {
		// ファイルが既に存在しない場合はエラーを無視
		if !os.IsNotExist(err) {
			return fmt.Errorf("物理ファイルの削除に失敗しました: %w", err)
		}
	}

	return nil
}
