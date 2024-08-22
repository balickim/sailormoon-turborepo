import { fetchWrapper } from '@renderer/utils/fetchWrapper'

export default class SlipsApi {
  private readonly basePath = `${import.meta.env.VITE_API_URL}/slips`

  getSlips = async () => {
    return fetchWrapper(`${this.basePath}`)
  }
}
