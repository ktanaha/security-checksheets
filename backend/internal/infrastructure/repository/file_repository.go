package repository

import (
	"database/sql"
	"time"

	"github.com/security-checksheets/backend/internal/domain"
)

// FileRepositoryImpl はFileRepositoryの実装
type FileRepositoryImpl struct {
	db *sql.DB
}

// NewFileRepository は新しいFileRepositoryを生成する
func NewFileRepository(db *sql.DB) domain.FileRepository {
	return &FileRepositoryImpl{db: db}
}

// Create は新規ファイルを登録する
func (r *FileRepositoryImpl) Create(file *domain.UploadedFile) error {
	query := `
		INSERT INTO uploaded_files (project_id, file_name, file_path, file_size, uploaded_by, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, uploaded_at
	`

	err := r.db.QueryRow(
		query,
		file.ProjectID,
		file.FileName,
		file.FilePath,
		file.FileSize,
		file.UploadedBy,
		time.Now(),
	).Scan(&file.ID, &file.UploadedAt)

	return err
}

// GetByID は指定されたIDのファイルを取得する
func (r *FileRepositoryImpl) GetByID(id int) (*domain.UploadedFile, error) {
	query := `
		SELECT id, project_id, file_name, file_path, file_size, uploaded_by, uploaded_at
		FROM uploaded_files
		WHERE id = $1
	`

	file := &domain.UploadedFile{}
	err := r.db.QueryRow(query, id).Scan(
		&file.ID,
		&file.ProjectID,
		&file.FileName,
		&file.FilePath,
		&file.FileSize,
		&file.UploadedBy,
		&file.UploadedAt,
	)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetByProjectID は指定された案件のすべてのファイルを取得する
func (r *FileRepositoryImpl) GetByProjectID(projectID int) ([]*domain.UploadedFile, error) {
	query := `
		SELECT id, project_id, file_name, file_path, file_size, uploaded_by, uploaded_at
		FROM uploaded_files
		WHERE project_id = $1
		ORDER BY uploaded_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []*domain.UploadedFile{}
	for rows.Next() {
		file := &domain.UploadedFile{}
		err := rows.Scan(
			&file.ID,
			&file.ProjectID,
			&file.FileName,
			&file.FilePath,
			&file.FileSize,
			&file.UploadedBy,
			&file.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// Delete はファイルを削除する
func (r *FileRepositoryImpl) Delete(id int) error {
	query := `DELETE FROM uploaded_files WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
