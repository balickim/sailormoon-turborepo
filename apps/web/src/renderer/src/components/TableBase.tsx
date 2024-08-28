import { CSSProperties } from 'react'
import {
  Column,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable
} from '@tanstack/react-table'
import {
  DndContext,
  KeyboardSensor,
  MouseSensor,
  TouchSensor,
  closestCenter,
  type DragEndEvent,
  useSensor,
  useSensors
} from '@dnd-kit/core'
import { restrictToHorizontalAxis } from '@dnd-kit/modifiers'
import { arrayMove, SortableContext, horizontalListSortingStrategy } from '@dnd-kit/sortable'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'

const DraggableTableHeader = ({ header }) => {
  const { attributes, isDragging, listeners, setNodeRef, transform } = useSortable({
    id: header.column.id
  })

  const style: CSSProperties = {
    opacity: isDragging ? 0.8 : 1,
    position: 'relative',
    transform: CSS.Translate.toString(transform),
    transition: 'width transform 0.2s ease-in-out',
    whiteSpace: 'nowrap',
    width: header.column.getSize(),
    zIndex: isDragging ? 1 : 0,
    borderWidth: '2px'
  }

  return (
    <th colSpan={header.colSpan} ref={setNodeRef} style={style}>
      <div
        onClick={header.column.getToggleSortingHandler()}
        style={{ cursor: 'pointer' }}
        className="flex justify-center gap-1"
      >
        {header.isPlaceholder
          ? null
          : flexRender(header.column.columnDef.header, header.getContext())}
        <span>
          {{
            asc: ' 🔼',
            desc: ' 🔽'
          }[header.column.getIsSorted() as string] ?? null}
        </span>
        <button {...attributes} {...listeners}>
          🟰
        </button>
      </div>
      {header.column.getCanFilter() ? (
        <div>
          <Filter column={header.column} />
        </div>
      ) : null}
    </th>
  )
}

const DragAlongCell = ({ cell }) => {
  const { isDragging, setNodeRef, transform } = useSortable({
    id: cell.column.id
  })

  const style: CSSProperties = {
    opacity: isDragging ? 0.8 : 1,
    position: 'relative',
    transform: CSS.Translate.toString(transform),
    transition: 'width transform 0.2s ease-in-out',
    width: cell.column.getSize(),
    zIndex: isDragging ? 1 : 0
  }

  return (
    <td style={style} ref={setNodeRef} className="border-2 pl-3">
      {flexRender(cell.column.columnDef.cell, cell.getContext())}
    </td>
  )
}

export const TableBase = ({
  columns,
  data,
  columnOrder,
  setColumnOrder,
  sorting,
  setSorting,
  pagination,
  setPagination,
  columnFilters,
  setColumnFilters,
  globalFilter,
  setGlobalFilter,
  pageCount,
  refetch
}) => {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    state: {
      columnOrder,
      sorting,
      pagination,
      columnFilters,
      globalFilter
    },
    onColumnOrderChange: setColumnOrder,
    onSortingChange: (updatedSorting) => {
      setSorting(updatedSorting)
      setPagination({ pageIndex: 0, pageSize: pagination.pageSize })
      refetch()
    },
    onPaginationChange: setPagination,
    onColumnFiltersChange: setColumnFilters,
    onGlobalFilterChange: setGlobalFilter,
    manualPagination: true,
    pageCount: pageCount,
    debugTable: true
  })

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event
    if (active && over && active.id !== over.id) {
      setColumnOrder((columnOrder) => {
        const oldIndex = columnOrder.indexOf(active.id as string)
        const newIndex = columnOrder.indexOf(over.id as string)
        return arrayMove(columnOrder, oldIndex, newIndex)
      })
    }
  }

  const sensors = useSensors(
    useSensor(MouseSensor, {}),
    useSensor(TouchSensor, {}),
    useSensor(KeyboardSensor, {})
  )

  return (
    <>
      <div className="pagination">
        <button onClick={() => table.setPageIndex(0)} disabled={!table.getCanPreviousPage()}>
          {'<<'}
        </button>
        <button onClick={() => table.previousPage()} disabled={!table.getCanPreviousPage()}>
          {'<'}
        </button>
        <span>
          Page{' '}
          <strong>
            {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}
          </strong>
        </span>
        <button onClick={() => table.nextPage()} disabled={!table.getCanNextPage()}>
          {'>'}
        </button>
        <button
          onClick={() => table.setPageIndex(table.getPageCount() - 1)}
          disabled={!table.getCanNextPage()}
        >
          {'>>'}
        </button>
        <select
          value={table.getState().pagination.pageSize}
          onChange={(e) => {
            table.setPageSize(Number(e.target.value))
          }}
        >
          {[10, 20, 30, 40, 50].map((pageSize) => (
            <option key={pageSize} value={pageSize}>
              Show {pageSize}
            </option>
          ))}
        </select>
      </div>

      <DndContext
        collisionDetection={closestCenter}
        modifiers={[restrictToHorizontalAxis]}
        onDragEnd={handleDragEnd}
        sensors={sensors}
      >
        <div className="w-full overflow-x-auto p-2">
          <table className="min-w-full border-2">
            <thead className="bg-gray-100 h-16">
              {table.getHeaderGroups().map((headerGroup) => (
                <tr key={headerGroup.id}>
                  <SortableContext items={columnOrder} strategy={horizontalListSortingStrategy}>
                    {headerGroup.headers.map((header) => (
                      <DraggableTableHeader key={header.id} header={header} />
                    ))}
                  </SortableContext>
                </tr>
              ))}
            </thead>
            <tbody>
              {table.getRowModel().rows.map((row) => (
                <tr key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <SortableContext
                      key={cell.id}
                      items={columnOrder}
                      strategy={horizontalListSortingStrategy}
                    >
                      <DragAlongCell key={cell.id} cell={cell} />
                    </SortableContext>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </DndContext>
    </>
  )
}

function Filter({ column }: { column: Column<unknown, unknown> }) {
  const columnFilterValue = column.getFilterValue()

  return (
    <input
      type="text"
      value={(columnFilterValue ?? '') as string}
      onChange={(e) => column.setFilterValue(e.target.value)}
      placeholder={`Search...`}
      className="w-36 border shadow rounded"
    />
  )
}
