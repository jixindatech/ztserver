import request from '@/utils/request'

export function add(data) {
  return request({
    url: `/api/v1/resource/`,
    method: 'post',
    data
  })
}

export function put(id, data) {
  return request({
    url: `/api/v1/resource/${id}`,
    method: 'put',
    data
  })
}

export function get(id) {
  return request({
    url: `/api/v1/resource/${id}`,
    method: 'get'
  })
}

export function getList(query, current = 1, size = 20) {
  return request({
    url: `/api/v1/resource/`,
    method: 'get',
    params: { ...query, current, size }
  })
}

export function deleteById(id) {
  return request({
    url: `/api/v1/resource/${id}`,
    method: 'delete'
  })
}

export function getCount() {
  return request({
    url: `/api/v1/info/resource/`,
    method: 'get'
  })
}
