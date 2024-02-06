// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

export async function getLogList(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.LogList>('/api/log', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}


/** 获取当前的用户 GET /api/currentUser */
export async function currentUser(options?: { [key: string]: any }) {
  return request<{
    data: API.CurrentUser;
  }>('/api/currentUser', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 退出登录接口 POST /api/login/outLogin */
export async function outLogin(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/login/outLogin', {
    method: 'POST',
    ...(options || {}),
  });
}

/** 登录接口 POST /api/login/account */
export async function login(body: API.LoginParams, options?: { [key: string]: any }) {
  return request<API.LoginResult>('/api/login/account', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 此处后端没有提供注释 GET /api/notices */
export async function getNotices(options?: { [key: string]: any }) {
  return request<API.NoticeIconList>('/api/notices', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取规则列表 GET /api/rule */
export async function rule(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.RuleList>('/api/rule', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 更新规则 PUT /api/rule */
export async function updateRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'POST',
    data:{
      method: 'update',
      ...(options || {}),
    }
  });
}

/** 新建规则 POST /api/rule */
export async function addRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'POST',
    data:{
      method: 'post',
      ...(options || {}),
    }
  });
}

/** 删除规则 DELETE /api/rule */
export async function removeRule(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/rule', {
    method: 'POST',
    data:{
      method: 'delete',
      ...(options || {}),
    }
  });
}



export async function getRoleList(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.RoleList>('/api/role', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
export async function addRole(body: any, options?: { [key: string]: any }) {
  return request<API.RoleItem>('/api/role', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
export async function delRole(id: string, options?: { [key: string]: any }) {
  return request<API.RoleItem>('/api/role', {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data: {id: id},
    ...(options || {}),
  });
}

export async function getMenu(
  options?: { [key: string]: any },
) {
  return request<API.MenuList>('/api/menu', {
    method: 'GET',
    ...(options || {}),
  });
}


export async function getAllRoles(
  options?: { [key: string]: any },
) {
  return request<API.UserList>('/api/all_roles', {
    method: 'GET',
    ...(options || {}),
  });
}

export async function getAllUsers(
  options?: { [key: string]: any },
) {
  return request<API.UserList>('/api/all_users', {
    method: 'GET',
    ...(options || {}),
  });
}

export async function getUserList(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.UserList>('/api/users', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
export async function addUser(body: any, options?: { [key: string]: any }) {
  return request<API.UserItem>('/api/user', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
export async function delUser(id: string, options?: { [key: string]: any }) {
  return request<API.UserItem>('/api/user', {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data: {id: id},
    ...(options || {}),
  });
}

export async function getPermissions(
  options?: { [key: string]: any },
) {
  return request('/api/permission', {
    method: 'GET',
    ...(options || {}),
  });
}
export async function updatePermission(body: any, options?: { [key: string]: any }) {
  return request('/api/permission', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function getRolePermissions(
  params: {
    current?: number;
    pageSize?: number;
    id?: number;
  },
  options?: { [key: string]: any },
) {
  return request('/api/role_permission', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
export async function updateRolePermission(body: any, options?: { [key: string]: any }) {
  return request('/api/role_permission', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}