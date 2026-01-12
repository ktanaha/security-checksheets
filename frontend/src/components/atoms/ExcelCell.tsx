import React from 'react';

interface ExcelCellProps {
  row: number;
  column: number;
  value: string | number | null;
  formattedValue: string;
  isMerged: boolean;
  mergeRange?: string;
  colSpan?: number;
  rowSpan?: number;
  isSelected: boolean;
  onCellClick: (row: number, column: number) => void;
}

/**
 * Excel風のセルコンポーネント
 */
const ExcelCell: React.FC<ExcelCellProps> = ({
  row,
  column,
  formattedValue,
  isMerged,
  colSpan,
  rowSpan,
  isSelected,
  onCellClick,
}) => {
  const handleClick = () => {
    onCellClick(row, column);
  };

  const baseClasses = 'border border-gray-300 px-2 py-1 text-sm cursor-pointer hover:bg-gray-50';
  const selectedClasses = isSelected ? 'bg-blue-100 border-blue-500' : 'bg-white';

  return (
    <td
      className={`${baseClasses} ${selectedClasses}`}
      colSpan={isMerged && colSpan ? colSpan : undefined}
      rowSpan={isMerged && rowSpan ? rowSpan : undefined}
      onClick={handleClick}
    >
      {formattedValue}
    </td>
  );
};

export default ExcelCell;
