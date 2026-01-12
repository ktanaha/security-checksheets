import axios from 'axios';
import { parseExcelFile, getSheetPreview } from '../excelService';
import type { ParseExcelResponse, SheetPreviewResponse } from '../../types/excel';

// axiosのモック
jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('excelService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('parseExcelFile', () => {
    it('正常にExcelファイルを解析できる', async () => {
      // Arrange
      const filePath = '/app/uploads/test.xlsx';
      const mockResponse: ParseExcelResponse = {
        file_name: 'test.xlsx',
        sheets: [
          {
            name: 'Sheet1',
            index: 0,
            row_count: 10,
            column_count: 5,
          },
        ],
        total_sheets: 1,
      };

      mockedAxios.post.mockResolvedValueOnce({ data: mockResponse });

      // Act
      const result = await parseExcelFile(filePath);

      // Assert
      expect(mockedAxios.post).toHaveBeenCalledWith(
        'http://localhost:8000/excel/parse',
        { file_path: filePath }
      );
      expect(result).toEqual(mockResponse);
    });

    it('エラー時に例外をスローする', async () => {
      // Arrange
      const filePath = '/app/uploads/nonexistent.xlsx';
      mockedAxios.post.mockRejectedValueOnce(new Error('File not found'));

      // Act & Assert
      await expect(parseExcelFile(filePath)).rejects.toThrow('File not found');
    });
  });

  describe('getSheetPreview', () => {
    it('デフォルトパラメータでプレビューを取得できる', async () => {
      // Arrange
      const filePath = '/app/uploads/test.xlsx';
      const mockResponse: SheetPreviewResponse = {
        sheet_name: 'Sheet1',
        cells: [
          {
            row: 1,
            column: 1,
            value: 'Test',
            formatted_value: 'Test',
            is_merged: false,
          },
        ],
        total_rows: 10,
        total_columns: 5,
      };

      mockedAxios.post.mockResolvedValueOnce({ data: mockResponse });

      // Act
      const result = await getSheetPreview({ filePath });

      // Assert
      expect(mockedAxios.post).toHaveBeenCalledWith(
        'http://localhost:8000/excel/preview',
        { file_path: filePath }
      );
      expect(result).toEqual(mockResponse);
    });

    it('すべてのパラメータを指定してプレビューを取得できる', async () => {
      // Arrange
      const params = {
        filePath: '/app/uploads/test.xlsx',
        sheetName: 'Sheet2',
        startRow: 1,
        endRow: 10,
        startColumn: 1,
        endColumn: 5,
      };

      const mockResponse: SheetPreviewResponse = {
        sheet_name: 'Sheet2',
        cells: [],
        total_rows: 10,
        total_columns: 5,
      };

      mockedAxios.post.mockResolvedValueOnce({ data: mockResponse });

      // Act
      const result = await getSheetPreview(params);

      // Assert
      expect(mockedAxios.post).toHaveBeenCalledWith(
        'http://localhost:8000/excel/preview',
        {
          file_path: params.filePath,
          sheet_name: params.sheetName,
          start_row: params.startRow,
          end_row: params.endRow,
          start_column: params.startColumn,
          end_column: params.endColumn,
        }
      );
      expect(result).toEqual(mockResponse);
    });

    it('エラー時に例外をスローする', async () => {
      // Arrange
      const params = { filePath: '/app/uploads/test.xlsx', sheetName: 'NonExistent' };
      mockedAxios.post.mockRejectedValueOnce(new Error('Sheet not found'));

      // Act & Assert
      await expect(getSheetPreview(params)).rejects.toThrow('Sheet not found');
    });
  });
});
