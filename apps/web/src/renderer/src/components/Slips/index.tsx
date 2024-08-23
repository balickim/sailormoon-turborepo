import SlipsApi from '@renderer/api/slips/routes'
import { useQuery } from '@tanstack/react-query'

export function Slips() {
  const slipsApi = new SlipsApi()

  const { data } = useQuery({
    queryKey: ['getSlips'],
    queryFn: () => slipsApi.getSlips(),
    refetchInterval: 5000
  })

  return <pre>{JSON.stringify(data, null, 2)}</pre>
}
