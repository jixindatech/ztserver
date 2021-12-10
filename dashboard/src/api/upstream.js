import request from '@/utils/request'

export function add(data) {
  return request({
    url: `/api/v1/upstream/`,
    method: 'post',
    data
  })
}

export function put(id, data) {
  return request({
    url: `/api/v1/upstream/${id}`,
    method: 'put',
    data
  })
}

export function get(id) {
  return request({
    url: `/api/v1/upstream/${id}`,
    method: 'get'
  })
}

export function getList(query, current = 1, size = 20) {
  return request({
    url: `/api/v1/upstream/`,
    method: 'get',
    params: { ...query, current, size }
  })
}

export function deleteById(id) {
  return request({
    url: `/api/v1/upstream/${id}`,
    method: 'delete'
  })
}
