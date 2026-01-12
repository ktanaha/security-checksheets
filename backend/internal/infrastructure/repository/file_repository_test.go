package repository

import (
	"database/sql"
	"testing"

	"github.com/security-checksheets/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// テスト用の案件を作成
	projectRepo := NewProjectRepository(db)
	project := domain.NewProject("テスト株式会社", "テスト案件", "山田太郎")
	require.NoError(t, projectRepo.Create(project))

	// ファイルリポジトリのテスト
	fileRepo := NewFileRepository(db)
	file := domain.NewUploadedFile(project.ID, "test.xlsx", "/uploads/test.xlsx", 12345, "山田太郎")

	err := fileRepo.Create(file)
	assert.NoError(t, err)
	assert.NotZero(t, file.ID, "IDが設定されるべき")
}

func TestFileRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// テスト用の案件を作成
	projectRepo := NewProjectRepository(db)
	project := domain.NewProject("テスト株式会社", "テスト案件", "山田太郎")
	require.NoError(t, projectRepo.Create(project))

	// ファイルを作成
	fileRepo := NewFileRepository(db)
	file := domain.NewUploadedFile(project.ID, "test.xlsx", "/uploads/test.xlsx", 12345, "山田太郎")
	require.NoError(t, fileRepo.Create(file))

	// 取得テスト
	fetched, err := fileRepo.GetByID(file.ID)
	assert.NoError(t, err)
	assert.Equal(t, file.FileName, fetched.FileName)
	assert.Equal(t, file.FilePath, fetched.FilePath)
	assert.Equal(t, file.FileSize, fetched.FileSize)
	assert.Equal(t, file.ProjectID, fetched.ProjectID)
}

func TestFileRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	fileRepo := NewFileRepository(db)

	// 存在しないIDで取得
	_, err := fileRepo.GetByID(99999)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestFileRepository_GetByProjectID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// テスト用の案件を作成
	projectRepo := NewProjectRepository(db)
	project := domain.NewProject("テスト株式会社", "テスト案件", "山田太郎")
	require.NoError(t, projectRepo.Create(project))

	// 複数のファイルを作成
	fileRepo := NewFileRepository(db)
	file1 := domain.NewUploadedFile(project.ID, "test1.xlsx", "/uploads/test1.xlsx", 12345, "山田太郎")
	file2 := domain.NewUploadedFile(project.ID, "test2.xlsx", "/uploads/test2.xlsx", 67890, "山田太郎")

	require.NoError(t, fileRepo.Create(file1))
	require.NoError(t, fileRepo.Create(file2))

	// プロジェクトIDでファイルを取得
	files, err := fileRepo.GetByProjectID(project.ID)
	assert.NoError(t, err)
	assert.Len(t, files, 2)
}

func TestFileRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// テスト用の案件を作成
	projectRepo := NewProjectRepository(db)
	project := domain.NewProject("テスト株式会社", "テスト案件", "山田太郎")
	require.NoError(t, projectRepo.Create(project))

	// ファイルを作成
	fileRepo := NewFileRepository(db)
	file := domain.NewUploadedFile(project.ID, "test.xlsx", "/uploads/test.xlsx", 12345, "山田太郎")
	require.NoError(t, fileRepo.Create(file))

	// 削除
	err := fileRepo.Delete(file.ID)
	assert.NoError(t, err)

	// 削除されたことを確認
	_, err = fileRepo.GetByID(file.ID)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}
