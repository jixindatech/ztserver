import request from '@/utils/request'

export function add(data) {
  return request({
    url: `/api/v1/user/`,
    method: 'post',
    data
  })
}

export function put(id, data) {
  return request({
    url: `/api/v1/user/${id}`,
    method: 'put',
    data
  })
}

export function get(id) {
  return request({
    url: `/api/v1/user/${id}`,
    method: 'get'
  })
}

export function getList(query, current = 1, size = 20) {
  return request({
    url: `/api/v1/user/`,
    method: 'get',
    params: { ...query, current, size }
  })
}

export function deleteById(id) {
  return request({
    url: `/api/v1/user/${id}`,
    method: 'delete'
  })
}

export function saveUserResource(id, ids) {
  return request({
    url: `/api/v1/user/${id}/resource/`,
    method: 'post',
    data: ids
  })
}

export function getUserResource(id, ids) {
  return request({
    url: `/api/v1/user/${id}/resource/`,
    method: 'get',
    data: ids
  })
}

export function sendMail(id, data) {
  return request({
    url: `/api/v1/user/${id}/email/`,
    method: 'post',
    data
  })
}

export function getCount() {
  return request({
    url: `/api/v1/info/user/`,
    method: 'get'
  })
}

export function getOnlineCount() {
  return request({
    url: `/api/v1/info/user/online/`,
    method: 'get'
  })
}

export function getUserGwInfo(query, current = 1, size = 20) {
  return request({
    url: `/api/v1/info/user/gw/`,
    method: 'get',
    params: { ...query }
  })
}
