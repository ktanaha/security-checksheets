"""Excel解析サービスのテスト"""

import pytest
from app.services.excel_parser import ExcelParser


class TestExcelParser:
    """ExcelParserのテストクラス"""

    def test_parse_excel_basic(self, temp_excel_file):
        """基本的なExcel解析のテスト"""
        parser = ExcelParser()
        result = parser.parse_excel(temp_excel_file)

        assert result is not None
        assert result['file_name'] == 'test_sample.xlsx'
        assert result['total_sheets'] == 3
        assert len(result['sheets']) == 3

        # Sheet1の確認
        sheet1 = result['sheets'][0]
        assert sheet1['name'] == 'Sheet1'
        assert sheet1['index'] == 0
        assert sheet1['row_count'] > 0
        assert sheet1['column_count'] > 0

    def test_parse_excel_sheet_info(self, temp_excel_file):
        """シート情報の正確性をテスト"""
        parser = ExcelParser()
        result = parser.parse_excel(temp_excel_file)

        sheets = result['sheets']

        # Sheet1: 3行2列のデータ
        assert sheets[0]['name'] == 'Sheet1'
        assert sheets[0]['row_count'] == 3
        assert sheets[0]['column_count'] == 2

        # Sheet2: 3行3列のデータ（セル結合あり）
        assert sheets[1]['name'] == 'Sheet2'
        assert sheets[1]['row_count'] == 3
        assert sheets[1]['column_count'] == 3

        # Sheet3: 空のシート
        assert sheets[2]['name'] == 'Sheet3'

    def test_get_sheet_preview_default(self, temp_excel_file):
        """デフォルトパラメータでのシートプレビュー取得テスト"""
        parser = ExcelParser()
        result = parser.get_sheet_preview(temp_excel_file)

        assert result is not None
        assert result['sheet_name'] == 'Sheet1'
        assert len(result['cells']) > 0
        assert result['row_count'] > 0
        assert result['column_count'] > 0

    def test_get_sheet_preview_specific_sheet(self, temp_excel_file):
        """特定シートのプレビュー取得テスト"""
        parser = ExcelParser()
        result = parser.get_sheet_preview(temp_excel_file, sheet_name='Sheet2')

        assert result['sheet_name'] == 'Sheet2'
        assert len(result['cells']) > 0

    def test_get_sheet_preview_range(self, temp_excel_file):
        """範囲指定でのプレビュー取得テスト"""
        parser = ExcelParser()
        result = parser.get_sheet_preview(
            temp_excel_file,
            start_row=1,
            end_row=2,
            start_column=1,
            end_column=2
        )

        assert result is not None
        # 範囲内のセルのみ取得されている
        cells = result['cells']
        for cell in cells:
            assert 1 <= cell['row'] <= 2
            assert 1 <= cell['column'] <= 2

    def test_get_sheet_preview_merged_cells(self, temp_excel_file):
        """セル結合を含むプレビュー取得テスト"""
        parser = ExcelParser()
        result = parser.get_sheet_preview(temp_excel_file, sheet_name='Sheet2')

        cells = result['cells']

        # セル結合情報が含まれているか確認
        merged_cells = [c for c in cells if c['is_merged']]
        assert len(merged_cells) > 0

        # A2セル（結合セルの開始）を確認
        a2_cell = next((c for c in cells if c['row'] == 2 and c['column'] == 1), None)
        assert a2_cell is not None
        assert a2_cell['is_merged'] is True
        assert a2_cell['merge_range'] is not None

    def test_parse_large_file(self, large_excel_file):
        """大量データを含むファイルの解析テスト"""
        parser = ExcelParser()
        result = parser.parse_excel(large_excel_file)

        assert result is not None
        assert result['total_sheets'] == 1

        sheet = result['sheets'][0]
        assert sheet['name'] == 'LargeSheet'
        assert sheet['row_count'] == 101  # ヘッダー + 100行

    def test_parse_empty_file(self, empty_excel_file):
        """空のファイルの解析テスト"""
        parser = ExcelParser()
        result = parser.parse_excel(empty_excel_file)

        assert result is not None
        assert result['total_sheets'] == 1

    def test_parse_nonexistent_file(self):
        """存在しないファイルの解析テスト"""
        parser = ExcelParser()

        with pytest.raises(FileNotFoundError):
            parser.parse_excel('/nonexistent/path/file.xlsx')

    def test_get_sheet_preview_nonexistent_sheet(self, temp_excel_file):
        """存在しないシート名の指定テスト"""
        parser = ExcelParser()

        with pytest.raises(ValueError):
            parser.get_sheet_preview(temp_excel_file, sheet_name='NonExistentSheet')

    def test_cell_value_formatting(self, temp_excel_file):
        """セル値のフォーマット確認テスト"""
        parser = ExcelParser()
        result = parser.get_sheet_preview(temp_excel_file)

        cells = result['cells']

        # A1セル（"質問"）を確認
        a1_cell = next((c for c in cells if c['row'] == 1 and c['column'] == 1), None)
        assert a1_cell is not None
        assert a1_cell['value'] == '質問'
        assert a1_cell['formatted_value'] == '質問'
