package repository

import (
	"database/sql"
	"testing"

	"github.com/security-checksheets/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq"
)

// setupTestDB はテスト用のDBコネクションをセットアップする
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=postgres port=5432 user=admin password=password dbname=security_checksheets sslmode=disable")
	require.NoError(t, err)

	// テスト前にデータをクリア
	_, err = db.Exec("DELETE FROM knowledge_items")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM uploaded_files")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM projects")
	require.NoError(t, err)

	return db
}

func TestProjectRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)
	project := domain.NewProject("テスト株式会社", "テスト案件", "山田太郎")

	err := repo.Create(project)
	assert.NoError(t, err)
	assert.NotZero(t, project.ID, "IDが設定されるべき")
}

func TestProjectRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)
	project := domain.NewProject("テスト株式会社", "テスト案件", "山田太郎")
	err := repo.Create(project)
	require.NoError(t, err)

	// 取得テスト
	fetched, err := repo.GetByID(project.ID)
	assert.NoError(t, err)
	assert.Equal(t, project.CustomerName, fetched.CustomerName)
	assert.Equal(t, project.Description, fetched.Description)
	assert.Equal(t, project.Owner, fetched.Owner)
	assert.Equal(t, "active", fetched.Status)
}

func TestProjectRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)

	// 存在しないIDで取得
	_, err := repo.GetByID(99999)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestProjectRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)

	// 複数の案件を作成
	project1 := domain.NewProject("テスト株式会社1", "案件1", "山田太郎")
	project2 := domain.NewProject("テスト株式会社2", "案件2", "佐藤花子")

	require.NoError(t, repo.Create(project1))
	require.NoError(t, repo.Create(project2))

	// 全件取得
	projects, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, projects, 2)
}

func TestProjectRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)
	project := domain.NewProject("テスト株式会社", "テスト案件", "山田太郎")
	require.NoError(t, repo.Create(project))

	// 更新
	project.CustomerName = "更新株式会社"
	project.Description = "更新された案件"
	project.Status = "completed"

	err := repo.Update(project)
	assert.NoError(t, err)

	// 更新されたことを確認
	fetched, err := repo.GetByID(project.ID)
	require.NoError(t, err)
	assert.Equal(t, "更新株式会社", fetched.CustomerName)
	assert.Equal(t, "更新された案件", fetched.Description)
	assert.Equal(t, "completed", fetched.Status)
}

func TestProjectRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)
	project := domain.NewProject("テスト株式会社", "テスト案件", "山田太郎")
	require.NoError(t, repo.Create(project))

	// 削除
	err := repo.Delete(project.ID)
	assert.NoError(t, err)

	// 削除されたことを確認
	_, err = repo.GetByID(project.ID)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}
