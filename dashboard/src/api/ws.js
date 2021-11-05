import request from '@/utils/request'

export function getList(query, current = 1, size = 20) {
  return request({
    url: `/api/v1/event/ws/`,
    method: 'get',
    params: { ...query, current, size }
  })
}
