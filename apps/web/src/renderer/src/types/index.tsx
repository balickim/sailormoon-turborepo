export interface IApiResponse<T> {
  statusCode: number
  data: T
}

export interface IApiListResponse<T> extends IApiResponse<T> {
  meta: {
    total: number
    current_page: number
    page_size: number
    last_page: number
  }
}
