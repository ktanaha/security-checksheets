import React, { useState } from 'react';
import ExcelCell from '../atoms/ExcelCell';
import type { CellData, CellRange } from '../../types/excel';

interface ExcelTableProps {
  cells: CellData[];
  totalRows: number;
  totalColumns: number;
  onRangeSelect?: (range: CellRange) => void;
}

/**
 * Excel風のテーブルコンポーネント
 */
const ExcelTable: React.FC<ExcelTableProps> = ({
  cells,
  totalRows,
  totalColumns,
  onRangeSelect,
}) => {
  const [selectedCell, setSelectedCell] = useState<{ row: number; column: number } | null>(null);
  const [isShiftPressed, setIsShiftPressed] = useState(false);
  const [rangeStart, setRangeStart] = useState<{ row: number; column: number } | null>(null);

  // セルのマップを作成（高速検索用）
  const cellMap = React.useMemo(() => {
    const map = new Map<string, CellData>();
    cells.forEach((cell) => {
      const key = `${cell.row}-${cell.column}`;
      map.set(key, cell);
    });
    return map;
  }, [cells]);

  // 列名を取得（A, B, C, ...）
  const getColumnName = (index: number): string => {
    let name = '';
    let num = index;
    while (num > 0) {
      const remainder = (num - 1) % 26;
      name = String.fromCharCode(65 + remainder) + name;
      num = Math.floor((num - 1) / 26);
    }
    return name;
  };

  // セル結合情報を解析
  const parseMergeRange = (
    mergeRange: string
  ): { colSpan: number; rowSpan: number } | null => {
    // 例: "A1:B2" -> colSpan=2, rowSpan=2
    const match = mergeRange.match(/([A-Z]+)(\d+):([A-Z]+)(\d+)/);
    if (!match) return null;

    const startCol = match[1];
    const startRow = parseInt(match[2]);
    const endCol = match[3];
    const endRow = parseInt(match[4]);

    // 列名をインデックスに変換
    const colToIndex = (col: string): number => {
      let index = 0;
      for (let i = 0; i < col.length; i++) {
        index = index * 26 + (col.charCodeAt(i) - 64);
      }
      return index;
    };

    const colSpan = colToIndex(endCol) - colToIndex(startCol) + 1;
    const rowSpan = endRow - startRow + 1;

    return { colSpan, rowSpan };
  };

  // セルクリックハンドラー
  const handleCellClick = (row: number, column: number) => {
    if (isShiftPressed && rangeStart) {
      // 範囲選択
      const range: CellRange = {
        startRow: Math.min(rangeStart.row, row),
        endRow: Math.max(rangeStart.row, row),
        startColumn: Math.min(rangeStart.column, column),
        endColumn: Math.max(rangeStart.column, column),
      };
      onRangeSelect?.(range);
    } else {
      // 単一セル選択
      setSelectedCell({ row, column });
      setRangeStart({ row, column });
    }
  };

  // キーボードイベントリスナー
  React.useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Shift') {
        setIsShiftPressed(true);
      }
    };

    const handleKeyUp = (e: KeyboardEvent) => {
      if (e.key === 'Shift') {
        setIsShiftPressed(false);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    window.addEventListener('keyup', handleKeyUp);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
      window.removeEventListener('keyup', handleKeyUp);
    };
  }, []);

  // セルが結合セルの一部かどうかをチェック
  const isCellPartOfMerge = (row: number, column: number): boolean => {
    for (const cell of cells) {
      if (cell.is_merged && cell.merge_range) {
        const match = cell.merge_range.match(/([A-Z]+)(\d+):([A-Z]+)(\d+)/);
        if (match) {
          const startRow = parseInt(match[2]);
          const endRow = parseInt(match[4]);

          const colToIndex = (col: string): number => {
            let index = 0;
            for (let i = 0; i < col.length; i++) {
              index = index * 26 + (col.charCodeAt(i) - 64);
            }
            return index;
          };

          const startCol = colToIndex(match[1]);
          const endCol = colToIndex(match[3]);

          if (
            row >= startRow &&
            row <= endRow &&
            column >= startCol &&
            column <= endCol &&
            !(row === cell.row && column === cell.column)
          ) {
            return true; // このセルは結合セルの一部なのでスキップ
          }
        }
      }
    }
    return false;
  };

  return (
    <div className="overflow-auto border border-gray-300 rounded">
      <table className="border-collapse">
        <thead>
          <tr>
            <th className="border border-gray-300 bg-gray-100 px-2 py-1 text-xs font-semibold sticky left-0 z-10">
              {/* 左上の空セル */}
            </th>
            {Array.from({ length: totalColumns }, (_, i) => (
              <th
                key={i}
                className="border border-gray-300 bg-gray-100 px-2 py-1 text-xs font-semibold"
              >
                {getColumnName(i + 1)}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: totalRows }, (_, rowIndex) => (
            <tr key={rowIndex}>
              <th className="border border-gray-300 bg-gray-100 px-2 py-1 text-xs font-semibold sticky left-0 z-10">
                {rowIndex + 1}
              </th>
              {Array.from({ length: totalColumns }, (_, colIndex) => {
                const row = rowIndex + 1;
                const column = colIndex + 1;

                // 結合セルの一部の場合はスキップ
                if (isCellPartOfMerge(row, column)) {
                  return null;
                }

                const key = `${row}-${column}`;
                const cellData = cellMap.get(key);

                // セルデータが存在しない場合は空セル
                if (!cellData) {
                  return (
                    <ExcelCell
                      key={key}
                      row={row}
                      column={column}
                      value={null}
                      formattedValue=""
                      isMerged={false}
                      isSelected={selectedCell?.row === row && selectedCell?.column === column}
                      onCellClick={handleCellClick}
                    />
                  );
                }

                // 結合セルの処理
                let colSpan: number | undefined;
                let rowSpan: number | undefined;

                if (cellData.is_merged && cellData.merge_range) {
                  const mergeInfo = parseMergeRange(cellData.merge_range);
                  if (mergeInfo) {
                    colSpan = mergeInfo.colSpan;
                    rowSpan = mergeInfo.rowSpan;
                  }
                }

                return (
                  <ExcelCell
                    key={key}
                    row={row}
                    column={column}
                    value={cellData.value}
                    formattedValue={cellData.formatted_value}
                    isMerged={cellData.is_merged}
                    mergeRange={cellData.merge_range}
                    colSpan={colSpan}
                    rowSpan={rowSpan}
                    isSelected={selectedCell?.row === row && selectedCell?.column === column}
                    onCellClick={handleCellClick}
                  />
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ExcelTable;
