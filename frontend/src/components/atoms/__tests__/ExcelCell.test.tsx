import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ExcelCell from '../ExcelCell';

// ヘルパー関数: テーブル構造でセルをレンダリング
const renderCell = (cell: React.ReactElement) => {
  return render(
    <table>
      <tbody>
        <tr>{cell}</tr>
      </tbody>
    </table>
  );
};

describe('ExcelCell', () => {
  it('通常のセル値を表示する', () => {
    renderCell(
      <ExcelCell
        row={1}
        column={1}
        value="テスト"
        formattedValue="テスト"
        isMerged={false}
        isSelected={false}
        onCellClick={() => {}}
      />
    );

    expect(screen.getByText('テスト')).toBeInTheDocument();
  });

  it('数値を表示する', () => {
    renderCell(
      <ExcelCell
        row={1}
        column={1}
        value={123}
        formattedValue="123"
        isMerged={false}
        isSelected={false}
        onCellClick={() => {}}
      />
    );

    expect(screen.getByText('123')).toBeInTheDocument();
  });

  it('nullの場合は空文字を表示する', () => {
    const { container } = renderCell(
      <ExcelCell
        row={1}
        column={1}
        value={null}
        formattedValue=""
        isMerged={false}
        isSelected={false}
        onCellClick={() => {}}
      />
    );

    const cell = container.querySelector('td');
    expect(cell?.textContent).toBe('');
  });

  it('選択時にスタイルが変わる', () => {
    const { container } = renderCell(
      <ExcelCell
        row={1}
        column={1}
        value="テスト"
        formattedValue="テスト"
        isMerged={false}
        isSelected={true}
        onCellClick={() => {}}
      />
    );

    const cell = container.querySelector('td');
    expect(cell).toHaveClass('bg-blue-100');
  });

  it('クリック時にコールバックが呼ばれる', async () => {
    const handleClick = jest.fn();
    const user = userEvent.setup();

    renderCell(
      <ExcelCell
        row={1}
        column={1}
        value="テスト"
        formattedValue="テスト"
        isMerged={false}
        isSelected={false}
        onCellClick={handleClick}
      />
    );

    const cell = screen.getByText('テスト');
    await user.click(cell);

    expect(handleClick).toHaveBeenCalledWith(1, 1);
  });

  it('結合セルの場合はcolSpan/rowSpanが設定される', () => {
    const { container } = renderCell(
      <ExcelCell
        row={1}
        column={1}
        value="結合セル"
        formattedValue="結合セル"
        isMerged={true}
        mergeRange="A1:B2"
        colSpan={2}
        rowSpan={2}
        isSelected={false}
        onCellClick={() => {}}
      />
    );

    const cell = container.querySelector('td');
    expect(cell).toHaveAttribute('colspan', '2');
    expect(cell).toHaveAttribute('rowspan', '2');
  });
});
