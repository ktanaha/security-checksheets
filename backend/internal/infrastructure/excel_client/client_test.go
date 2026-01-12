package excel_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExcelClient_ParseExcel(t *testing.T) {
	// Excel Serviceが起動していることを前提とする統合テスト
	client := NewExcelClient("http://excel-service:8000")

	// テスト用ファイルパス（事前にアップロードされたファイル）
	filePath := "/app/uploads/test_security_check.xlsx"

	result, err := client.ParseExcel(filePath)
	require.NoError(t, err, "Excel解析は成功するべき")
	assert.NotNil(t, result)

	// ファイル情報の確認
	assert.Equal(t, "test_security_check.xlsx", result.FileName)
	assert.Equal(t, filePath, result.FilePath)

	// シート情報の確認
	assert.Equal(t, 1, result.TotalSheets)
	assert.Len(t, result.Sheets, 1)

	sheet := result.Sheets[0]
	assert.Equal(t, "セキュリティチェック", sheet.Name)
	assert.Equal(t, 0, sheet.Index)
	assert.Greater(t, sheet.RowCount, 0)
	assert.Greater(t, sheet.ColumnCount, 0)
}

func TestExcelClient_GetSheetPreview(t *testing.T) {
	// Excel Serviceが起動していることを前提とする統合テスト
	client := NewExcelClient("http://excel-service:8000")

	// テスト用ファイルパス
	filePath := "/app/uploads/test_security_check.xlsx"
	sheetName := "セキュリティチェック"
	startRow := 1
	endRow := 2

	result, err := client.GetSheetPreview(filePath, &sheetName, &startRow, &endRow, nil, nil)
	require.NoError(t, err, "シートプレビュー取得は成功するべき")
	assert.NotNil(t, result)

	// プレビュー情報の確認
	assert.Equal(t, "セキュリティチェック", result.SheetName)
	assert.Greater(t, len(result.Cells), 0, "セルデータが取得されるべき")
	assert.Equal(t, 2, result.RowCount)
	assert.Greater(t, result.ColumnCount, 0)

	// セルデータの確認
	cells := result.Cells
	assert.Greater(t, len(cells), 0)

	// A1セル（ヘッダー）を確認
	found := false
	for _, cell := range cells {
		if cell.Row == 1 && cell.Column == 1 {
			found = true
			assert.Equal(t, "部門", cell.Value)
			break
		}
	}
	assert.True(t, found, "A1セルが見つかるべき")
}

func TestExcelClient_ParseExcel_NonExistentFile(t *testing.T) {
	client := NewExcelClient("http://excel-service:8000")

	// 存在しないファイルパス
	filePath := "/nonexistent/file.xlsx"

	_, err := client.ParseExcel(filePath)
	assert.Error(t, err, "存在しないファイルの場合はエラーになるべき")
}

func TestExcelClient_GetSheetPreview_NonExistentSheet(t *testing.T) {
	client := NewExcelClient("http://excel-service:8000")

	filePath := "/app/uploads/test_security_check.xlsx"
	sheetName := "NonExistentSheet"

	_, err := client.GetSheetPreview(filePath, &sheetName, nil, nil, nil, nil)
	assert.Error(t, err, "存在しないシート名の場合はエラーになるべき")
}
