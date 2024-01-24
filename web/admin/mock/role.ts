import { Request, Response } from 'express';

let roles = [
  {
    id: 'admin',
    name: '管理员',
    selected: ['/system', '/system/user', '/system/role'],
  },
  {
    id: 'user',
    name: '普通用户',
    selected: ['/system', '/system/user', '/system/role'],
  },
];

const getRole = (req: Request, res: Response) => {
  res.json({
    data: roles,
  });
};

export let  getRoleById = (id: string) => {
  for (let i = 0; i < roles.length; i++) {
    const element = roles[i];
    if (element.id == id) {
      return element.selected;
    }
  }

  return [];
}

export default {
  'GET /api/role': getRole,
};

