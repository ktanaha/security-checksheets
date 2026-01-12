import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ExcelTable from '../ExcelTable';
import type { CellData } from '../../../types/excel';

describe('ExcelTable', () => {
  const mockCells: CellData[] = [
    {
      row: 1,
      column: 1,
      value: 'A1',
      formatted_value: 'A1',
      is_merged: false,
    },
    {
      row: 1,
      column: 2,
      value: 'B1',
      formatted_value: 'B1',
      is_merged: false,
    },
    {
      row: 2,
      column: 1,
      value: 'A2',
      formatted_value: 'A2',
      is_merged: false,
    },
    {
      row: 2,
      column: 2,
      value: 'B2',
      formatted_value: 'B2',
      is_merged: false,
    },
  ];

  it('セルデータを表示する', () => {
    render(<ExcelTable cells={mockCells} totalRows={2} totalColumns={2} />);

    expect(screen.getByText('A1')).toBeInTheDocument();
    expect(screen.getByText('B1')).toBeInTheDocument();
    expect(screen.getByText('A2')).toBeInTheDocument();
    expect(screen.getByText('B2')).toBeInTheDocument();
  });

  it('行ヘッダーと列ヘッダーを表示する', () => {
    render(<ExcelTable cells={mockCells} totalRows={2} totalColumns={2} />);

    // 列ヘッダー（A, B）
    expect(screen.getByText('A')).toBeInTheDocument();
    expect(screen.getByText('B')).toBeInTheDocument();

    // 行ヘッダー（1, 2）
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
  });

  it('セルをクリックすると選択される', async () => {
    const user = userEvent.setup();
    const { container } = render(
      <ExcelTable cells={mockCells} totalRows={2} totalColumns={2} />
    );

    const cellA1 = screen.getByText('A1');
    await user.click(cellA1);

    // 選択されたセルはbg-blue-100クラスを持つ
    const selectedCell = container.querySelector('.bg-blue-100');
    expect(selectedCell).toBeInTheDocument();
  });

  it('範囲選択コールバックが呼ばれる', async () => {
    const handleRangeSelect = jest.fn();
    const user = userEvent.setup();

    render(
      <ExcelTable
        cells={mockCells}
        totalRows={2}
        totalColumns={2}
        onRangeSelect={handleRangeSelect}
      />
    );

    // A1をクリック
    await user.click(screen.getByText('A1'));
    // Shiftキーを押しながらB2をクリック（範囲選択）
    await user.keyboard('{Shift>}');
    await user.click(screen.getByText('B2'));
    await user.keyboard('{/Shift}');

    expect(handleRangeSelect).toHaveBeenCalledWith({
      startRow: 1,
      endRow: 2,
      startColumn: 1,
      endColumn: 2,
    });
  });

  it('結合セルを正しく表示する', () => {
    const mergedCells: CellData[] = [
      {
        row: 1,
        column: 1,
        value: 'Merged',
        formatted_value: 'Merged',
        is_merged: true,
        merge_range: 'A1:B2',
      },
      {
        row: 1,
        column: 3,
        value: 'C1',
        formatted_value: 'C1',
        is_merged: false,
      },
    ];

    const { container } = render(
      <ExcelTable cells={mergedCells} totalRows={2} totalColumns={3} />
    );

    expect(screen.getByText('Merged')).toBeInTheDocument();

    // 結合セルはcolspan/rowspanを持つ
    const mergedCell = container.querySelector('td[colspan="2"]');
    expect(mergedCell).toBeInTheDocument();
  });

  it('空のセルを正しく表示する', () => {
    const sparseCells: CellData[] = [
      {
        row: 1,
        column: 1,
        value: 'A1',
        formatted_value: 'A1',
        is_merged: false,
      },
      // (1,2)は空
      {
        row: 2,
        column: 1,
        value: null,
        formatted_value: '',
        is_merged: false,
      },
    ];

    render(<ExcelTable cells={sparseCells} totalRows={2} totalColumns={2} />);

    expect(screen.getByText('A1')).toBeInTheDocument();
    // 空セルも表示される（テーブル構造維持のため）
  });
});
