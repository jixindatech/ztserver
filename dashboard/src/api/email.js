import request from '@/utils/request'

export function get(data) {
  return request({
    url: `/api/v1/email/`,
    method: 'get'
  })
}

export function add(data) {
  return request({
    url: `/api/v1/email/`,
    method: 'post',
    data
  })
}

export function put(id, data) {
  return request({
    url: `/api/v1/email/${id}`,
    method: 'put',
    data
  })
}
