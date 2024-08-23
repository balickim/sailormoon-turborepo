import { IApiListResponse } from '@renderer/types'
import { fetchWrapper } from '@renderer/utils/fetchWrapper'

interface IListParams {
  sort_by?: string
  sort_order?: 'asc' | 'desc'
  page?: number
  page_size?: number
}

export default class SlipsApi {
  private readonly basePath = `${import.meta.env.VITE_API_URL}/slips`

  getSlips = async (params: IListParams = {}): Promise<IApiListResponse<unknown[]>> => {
    console.log(import.meta.env.VITE_API_URL)
    const queryString = new URLSearchParams(params).toString()
    const url = queryString ? `${this.basePath}?${queryString}` : `${this.basePath}`
    return fetchWrapper(url)
  }
}
