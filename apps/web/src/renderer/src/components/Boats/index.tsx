import { useQuery } from '@tanstack/react-query'
import BoatsApi from '@renderer/api/boats/routes'
import { useEffect, useMemo, useState } from 'react'
import { ColumnDef, SortingState, ColumnFiltersState } from '@tanstack/react-table'
import { TableBase } from '../TableBase'

export function Boats() {
  const boatsApi = new BoatsApi()

  const [sorting, setSorting] = useState<SortingState>([])
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10
  })
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [globalFilter, setGlobalFilter] = useState('')
  const [totalRows, setTotalRows] = useState(0)
  const [tableData, setTableData] = useState<unknown[]>([])

  const queryKey = useMemo(() => {
    return ['getBoats', sorting, pagination, columnFilters, globalFilter]
  }, [sorting, pagination, columnFilters, globalFilter])

  const { data, refetch } = useQuery({
    queryKey,
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

      return boatsApi.getBoats(params)
    },
    refetchOnWindowFocus: true
  })

  useEffect(() => {
    if (data) {
      setTableData(data.data)
      setTotalRows(data.meta.total)
    }
  }, [data])

  const columns = useMemo<ColumnDef<unknown>[]>(
    () => [
      {
        id: 'name',
        accessorKey: 'name',
        header: 'Name',
        cell: (info) => info.getValue(),
        filterFn: 'includesString'
      },
      {
        id: 'type',
        accessorKey: 'type',
        header: 'Type',
        cell: (info) => info.getValue(),
        filterFn: 'includesString'
      },
      {
        id: 'length',
        accessorKey: 'length',
        header: 'Length',
        cell: (info) => `${info.getValue()} meters`,
        filterFn: 'equals'
      },
      {
        id: 'width',
        accessorKey: 'width',
        header: 'Width',
        cell: (info) => `${info.getValue()} meters`,
        filterFn: 'equals'
      },
      {
        id: 'owners',
        accessorKey: 'owners',
        header: 'Owners',
        cell: ({ getValue }) => {
          const owners = getValue() as Array<{ firstName: string; lastName: string }> | undefined
          return (
            <ul>
              {owners?.map((owner, index) => (
                <li key={index}>
                  {owner.firstName} {owner.lastName}
                </li>
              )) || 'No owners'}
            </ul>
          )
        },
        filterFn: 'includesString'
      }
    ],
    []
  )

  const loadColumnOrder = () => {
    const stored = JSON.parse(localStorage.getItem('boatsColumnOrder') || 'null')
    if (stored && Array.isArray(stored)) {
      return stored
    }
    return columns.map((col) => col.id as string)
  }

  const [columnOrder, setColumnOrder] = useState<string[]>(loadColumnOrder())

  useEffect(() => {
    localStorage.setItem('boatsColumnOrder', JSON.stringify(columnOrder))
  }, [columnOrder])

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
    <TableBase
      columns={columns}
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
  )
}
