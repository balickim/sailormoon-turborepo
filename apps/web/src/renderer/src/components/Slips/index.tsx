import { CSSProperties, useEffect, useState, useMemo } from 'react'
import {
  ColumnDef,
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

import './index.css'

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
    zIndex: isDragging ? 1 : 0
  }

  return (
    <th colSpan={header.colSpan} ref={setNodeRef} style={style}>
      <div onClick={header.column.getToggleSortingHandler()} style={{ cursor: 'pointer' }}>
        {header.isPlaceholder
          ? null
          : flexRender(header.column.columnDef.header, header.getContext())}
        <span>
          {{
            asc: ' 🔼',
            desc: ' 🔽'
          }[header.column.getIsSorted() as string] ?? null}
        </span>
      </div>
      <button {...attributes} {...listeners}>
        🟰
      </button>
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
    <td style={style} ref={setNodeRef}>
      {flexRender(cell.column.columnDef.cell, cell.getContext())}
    </td>
  )
}

export function Slips() {
  const slipsApi = new SlipsApi()

  const [sorting, setSorting] = useState<SortingState>([])

  const { data, refetch } = useQuery({
    queryKey: ['getSlips', sorting],
    queryFn: () => {
      const sort_by = sorting[0]?.id
      const sort_order = sorting[0]?.desc ? 'desc' : 'asc'

      const params: {
        page: number
        page_size: number
        sort_by?: string
        sort_order?: 'asc' | 'desc'
      } = {
        page: 1,
        page_size: 100
      }

      if (sort_by) {
        params.sort_by = sort_by
        params.sort_order = sort_order
      }

      return slipsApi.getSlips(params)
    },
    refetchOnWindowFocus: true
  })

  const [tableData, setTableData] = useState<unknown[]>([])

  const defaultColumns = useMemo<ColumnDef<unknown>[]>(
    () => [
      {
        accessorKey: 'number',
        cell: (info) => info.getValue(),
        id: 'number',
        size: 150
      },
      {
        accessorKey: 'boats',
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
        header: 'Boats',
        id: 'boats',
        size: 300
      },
      {
        accessorKey: 'notes',
        header: 'Notatki',
        id: 'notes',
        size: 500
      }
    ],
    []
  )

  const storedColumnOrder = JSON.parse(localStorage.getItem('columnOrder') || 'null')
  const [columnOrder, setColumnOrder] = useState<string[]>(
    storedColumnOrder || defaultColumns.map((col) => col.id)
  )

  useEffect(() => {
    if (data) {
      setTableData(data.data)
    }
  }, [data])

  useEffect(() => {
    localStorage.setItem('columnOrder', JSON.stringify(columnOrder))
  }, [columnOrder])

  const table = useReactTable({
    data: tableData,
    columns: defaultColumns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    state: {
      columnOrder,
      sorting
    },
    onColumnOrderChange: setColumnOrder,
    onSortingChange: (updatedSorting) => {
      setSorting(updatedSorting)
      refetch()
    },
    debugTable: true
  })

  // Handle the drag & drop reordering of columns
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
    <DndContext
      collisionDetection={closestCenter}
      modifiers={[restrictToHorizontalAxis]}
      onDragEnd={handleDragEnd}
      sensors={sensors}
    >
      <div className="p-2">
        <table>
          <thead>
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
  )
}
