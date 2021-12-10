import request from '@/utils/request'

export function add(data) {
  return request({
    url: `/api/v1/ssl/`,
    method: 'post',
    data
  })
}

export function put(id, data) {
  return request({
    url: `/api/v1/ssl/${id}`,
    method: 'put',
    data
  })
}

export function get(id) {
  return request({
    url: `/api/v1/ssl/${id}`,
    method: 'get'
  })
}

export function getList(query, current = 1, size = 20) {
  return request({
    url: `/api/v1/ssl/`,
    method: 'get',
    params: { ...query, current, size }
  })
}

export function deleteById(id) {
  return request({
    url: `/api/v1/ssl/${id}`,
    method: 'delete'
  })
}
