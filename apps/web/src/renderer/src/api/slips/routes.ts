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
    const processedParams: Record<string, string> = {}

    Object.keys(params).forEach((key) => {
      const value = params[key]
      if (value !== undefined && value !== null) {
        if (typeof value === 'object') {
          processedParams[key] = JSON.stringify(value)
        } else {
          processedParams[key] = String(value)
        }
      }
    })

    const queryString = new URLSearchParams(processedParams).toString()
    const url = queryString ? `${this.basePath}?${queryString}` : `${this.basePath}`
    return fetchWrapper(url)
  }
}
