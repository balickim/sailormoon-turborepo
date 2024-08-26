import { CSSProperties, useEffect, useState, useMemo, useRef } from 'react'
import {
  Column,
  ColumnDef,
  ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  SortingState,
  useReactTable
} from '@tanstack/react-table'
import { useQuery } from '@tanstack/react-query'
import SlipsApi from '@renderer/api/slips/routes'
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
import AddSlipDialog from './AddSlipDialog'

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

export function Slips() {
  const slipsApi = new SlipsApi()

  const [sorting, setSorting] = useState<SortingState>([])
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10
  })
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [globalFilter, setGlobalFilter] = useState('')

  const STORAGE_VERSION = '1.0'
  const dialogRef = useRef()

  const { data, refetch } = useQuery({
    queryKey: ['getSlips', sorting, pagination, columnFilters, globalFilter],
    queryFn: () => {
      const params: Record<string, unknown> = {
        page: pagination.pageIndex + 1,
        page_size: pagination.pageSize,
        global_filter: globalFilter
      }

      if (sorting.length > 0) {
        params.sort_by = sorting[0].id
        params.sort_order = sorting[0].desc ? 'desc' : 'asc'
      }

      if (columnFilters.length > 0) {
        params.filters = columnFilters
      }

      return slipsApi.getSlips(params)
    },
    refetchOnWindowFocus: true
  })

  const [tableData, setTableData] = useState<unknown[]>([])

  const defaultColumns = useMemo<ColumnDef<unknown>[]>(
    () => [
      {
        id: 'number',
        accessorKey: 'number',
        header: 'Nr',
        size: 50,
        cell: (info) => info.getValue(),
        filterFn: 'equals'
      },
      {
        id: 'boats',
        accessorKey: 'boats',
        header: 'Łodzie',
        size: 300,
        cell: ({ getValue }) => {
          const boats = getValue() as Array<{ name: string; type: string }> | undefined
          return (
            <ul>
              {boats?.map((boat, index) => (
                <li key={index}>
                  <a href={boat.name}>
                    {boat.name} ({boat.type})
                  </a>
                </li>
              )) || 'No boats'}
            </ul>
          )
        },
        filterFn: 'includesString'
      },
      {
        id: 'notes',
        accessorKey: 'notes',
        header: 'Notatki',
        filterFn: 'includesString'
      },
      {
        id: 'CreatedAt',
        accessorKey: 'CreatedAt',
        header: 'Utworzono'
      }
    ],
    []
  )

  const loadColumnOrder = () => {
    const stored = JSON.parse(localStorage.getItem('columnOrder') || 'null')
    if (stored && stored.version === STORAGE_VERSION) {
      return stored.columnOrder
    }
    return defaultColumns.map((col) => col.id as string)
  }

  const [columnOrder, setColumnOrder] = useState<string[]>(loadColumnOrder())

  useEffect(() => {
    if (data) {
      setTableData(data.data)
    }
  }, [data])

  useEffect(() => {
    localStorage.setItem(
      'columnOrder',
      JSON.stringify({
        version: STORAGE_VERSION,
        columnOrder: columnOrder
      })
    )
  }, [columnOrder])

  const table = useReactTable({
    data: tableData,
    columns: defaultColumns,
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
    pageCount: data?.meta.total ?? -1,
    debugTable: true
  })

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event
    console.log('Drag end:', active?.id, 'over:', over?.id)
    if (active && over && active.id !== over.id) {
      setColumnOrder((columnOrder) => {
        const oldIndex = columnOrder.indexOf(active.id as string)
        const newIndex = columnOrder.indexOf(over.id as string)
        console.log('Reordering:', oldIndex, 'to', newIndex)
        return arrayMove(columnOrder, oldIndex, newIndex)
      })
    }
  }

  const sensors = useSensors(
    useSensor(MouseSensor, {}),
    useSensor(TouchSensor, {}),
    useSensor(KeyboardSensor, {})
  )

  const resetToDefaultOrder = () => {
    const defaultOrder = defaultColumns.map((col) => col.id as string)
    setColumnOrder(defaultOrder)
    localStorage.setItem(
      'columnOrder',
      JSON.stringify({
        version: STORAGE_VERSION,
        columnOrder: defaultOrder
      })
    )
  }

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

      <button onClick={() => resetToDefaultOrder()}> Reset to default </button>
      <div>
        <input
          value={globalFilter ?? ''}
          onChange={(e) => setGlobalFilter(e.target.value)}
          placeholder="Search all columns..."
          className="p-2 font-lg shadow border border-block"
        />
      </div>

      <button className="btn btn-primary" onClick={() => dialogRef.current?.open()}>
        Add
      </button>
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

      <AddSlipDialog ref={dialogRef} />
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
