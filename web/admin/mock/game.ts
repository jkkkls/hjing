import { Request, Response } from 'express';

const getGame = (req: Request, res: Response) => {
  res.json({
    data: [
      {
        id: 'xqhj',
        name: '象棋',
        appId: '123123',
        appSecret: '2017-08-09',
        dbName: '123',
        status: 1,
      },
    ],
  });
};

export default {
  'GET /api/game': getGame,
};
