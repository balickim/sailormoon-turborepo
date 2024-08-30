import { IApiListResponse } from '@renderer/types'
import { fetchWrapper } from '@renderer/utils/fetchWrapper'

export default class UsersApi {
  private readonly basePath = `${import.meta.env.VITE_API_URL}/users`

  login = async (body: {
    email: string
    password: string
  }): Promise<IApiListResponse<unknown[]>> => {
    return fetchWrapper(`${this.basePath}/login`, { method: 'POST', body: JSON.stringify(body) })
  }
}
