package excel_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ExcelClient はExcel処理API（Python）のクライアント
type ExcelClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewExcelClient は新しいExcelClientを生成する
func NewExcelClient(baseURL string) *ExcelClient {
	return &ExcelClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SheetInfo はシート情報
type SheetInfo struct {
	Name        string `json:"name"`
	Index       int    `json:"index"`
	RowCount    int    `json:"row_count"`
	ColumnCount int    `json:"column_count"`
}

// ParseExcelResponse はExcel解析APIのレスポンス
type ParseExcelResponse struct {
	FileName    string      `json:"file_name"`
	FilePath    string      `json:"file_path"`
	Sheets      []SheetInfo `json:"sheets"`
	TotalSheets int         `json:"total_sheets"`
}

// CellData はセルデータ
type CellData struct {
	Row            int         `json:"row"`
	Column         int         `json:"column"`
	Value          interface{} `json:"value"`
	FormattedValue *string     `json:"formatted_value"`
	IsMerged       bool        `json:"is_merged"`
	MergeRange     *string     `json:"merge_range"`
}

// SheetPreviewResponse はシートプレビューAPIのレスポンス
type SheetPreviewResponse struct {
	SheetName   string     `json:"sheet_name"`
	Cells       []CellData `json:"cells"`
	RowCount    int        `json:"row_count"`
	ColumnCount int        `json:"column_count"`
}

// ParseExcel はExcelファイルを解析する
func (c *ExcelClient) ParseExcel(filePath string) (*ParseExcelResponse, error) {
	requestBody := map[string]string{
		"file_path": filePath,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストJSONのマーシャルに失敗しました: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/excel/parse", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("Excel解析APIリクエストに失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Excel解析APIがエラーを返しました (status: %d): %s", resp.StatusCode, string(body))
	}

	var result ParseExcelResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}

	return &result, nil
}

// GetSheetPreview はシートのプレビューを取得する
func (c *ExcelClient) GetSheetPreview(
	filePath string,
	sheetName *string,
	startRow, endRow, startColumn, endColumn *int,
) (*SheetPreviewResponse, error) {
	requestBody := map[string]interface{}{
		"file_path": filePath,
	}

	if sheetName != nil {
		requestBody["sheet_name"] = *sheetName
	}
	if startRow != nil {
		requestBody["start_row"] = *startRow
	}
	if endRow != nil {
		requestBody["end_row"] = *endRow
	}
	if startColumn != nil {
		requestBody["start_column"] = *startColumn
	}
	if endColumn != nil {
		requestBody["end_column"] = *endColumn
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストJSONのマーシャルに失敗しました: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/excel/preview", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("シートプレビューAPIリクエストに失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("シートプレビューAPIがエラーを返しました (status: %d): %s", resp.StatusCode, string(body))
	}

	var result SheetPreviewResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}

	return &result, nil
}
