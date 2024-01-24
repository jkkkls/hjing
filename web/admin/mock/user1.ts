import { Request, Response } from 'express';
import { getRoleById } from './role';


const waitTime = (time: number = 100) => {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(true);
    }, time);
  });
};
let users = [
  {
    id: 'guangbo',
    avatar: 'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
    name: '广波',
    password: '123123',
    role: "admin",
    hasRoutes:['welcome'],
  },
  {
    id: 'xxoo',
    avatar: 'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
    name: 'xx11',
    password: '123123',
    role: "user",
    hasRoutes:['welcome'],
  },
];

let loginUser = {
  id: '',
  avatar: 'g',
  name: '',
  password: '',
  role: "",
  hasRoutes:['welcome'],
};

const getUser = (req: Request, res: Response) => {
  res.json({
    data: users,
  });
};

export default {
  'GET /api/users': getUser,

  'GET /api/currentUser': (req: Request, res: Response) => {
    console.log("---", loginUser)
    if (loginUser.id == "") {
      res.status(401).send({
        data: {
          isLogin: false,
        },
        errorCode: '401',
        errorMessage: '请先登录！',
        success: true,
      });
      return;
    }
    res.send({
      success: true,
      data: loginUser,
    });
  },
  'POST /api/login/account': async (req: Request, res: Response) => {
    const { password, username, type } = req.body;
    await waitTime(2000);
    for (let i = 0; i < users.length; i++) {
      const e = users[i];
      if (e.id == username) {
        if (e.password == password) {
          res.send({
            status: 'ok',
            type,
            currentAuthority: 'admin',
          });
          loginUser = users[i];
          loginUser.hasRoutes = getRoleById(loginUser.role) || [];
          console.log("-loginUser--", loginUser)
          return;
        }
      }

    }
    res.send({
      status: 'error',
      type,
      currentAuthority: 'guest',
    });
  },
  'POST /api/login/outLogin': (req: Request, res: Response) => {
    loginUser.id = "";
    loginUser.hasRoutes = [];
    res.send({ data: {}, success: true });
  },
};
