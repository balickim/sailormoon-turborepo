import { fetchWrapper } from '@renderer/utils/fetchWrapper'

interface IListParams {
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  page?: number
  pageSize?: number
}

export default class SlipsApi {
  private readonly basePath = `${import.meta.env.VITE_API_URL}/slips`

  getSlips = async (params: IListParams = {}) => {
    console.log(import.meta.env.VITE_API_URL)
    const queryString = new URLSearchParams(params).toString()
    const url = queryString ? `${this.basePath}?${queryString}` : `${this.basePath}`
    return fetchWrapper(url)
  }
}
