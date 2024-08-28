import { useState, useEffect, useMemo, useRef } from 'react'
import { useQuery } from '@tanstack/react-query'
import SlipsApi from '@renderer/api/slips/routes'
import { ColumnDef, ColumnFiltersState, SortingState } from '@tanstack/react-table'
import { TableBase } from '../TableBase'
import AddSlipDialog from './AddSlipDialog'

export function Slips() {
  const slipsApi = new SlipsApi()

  const [sorting, setSorting] = useState<SortingState>([])
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10
  })
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [globalFilter, setGlobalFilter] = useState('')
  const [totalRows, setTotalRows] = useState(0)
  const [tableData, setTableData] = useState<unknown[]>([])

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
      setTotalRows(data.meta.total)
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

  const handleSetColumnFilters = (filters: ColumnFiltersState) => {
    setColumnFilters(filters)
    setPagination((prev) => ({
      ...prev,
      pageIndex: 0
    }))
  }

  const handleSetGlobalFilter = (filter: string) => {
    setGlobalFilter(filter)
    setPagination((prev) => ({
      ...prev,
      pageIndex: 0
    }))
  }

  const pageCount = Math.ceil(totalRows / pagination.pageSize)

  return (
    <>
      <button onClick={() => resetToDefaultOrder()}> Reset to default </button>
      <div>
        <input
          value={globalFilter ?? ''}
          onChange={(e) => handleSetGlobalFilter(e.target.value)}
          placeholder="Search all columns..."
          className="p-2 font-lg shadow border border-block"
        />
      </div>

      <button className="btn btn-primary" onClick={() => dialogRef.current?.open()}>
        Add
      </button>

      <TableBase
        columns={defaultColumns}
        data={tableData}
        columnOrder={columnOrder}
        setColumnOrder={setColumnOrder}
        sorting={sorting}
        setSorting={setSorting}
        pagination={pagination}
        setPagination={setPagination}
        columnFilters={columnFilters}
        setColumnFilters={handleSetColumnFilters}
        globalFilter={globalFilter}
        setGlobalFilter={handleSetGlobalFilter}
        pageCount={pageCount}
        refetch={refetch}
      />

      <AddSlipDialog ref={dialogRef} />
    </>
  )
}
