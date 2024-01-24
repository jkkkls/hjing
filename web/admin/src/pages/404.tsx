import { Button, Result } from 'antd';
import React from 'react';
import { history } from '@umijs/max';

const NoFoundPage: React.FC = () => (
  <Result
    status="403"
    // title="404"
    subTitle="消失了？？？功能开发中，请稍后"
    extra={
      <Button type="primary" onClick={() => history.push('/')}>
        返回首页
      </Button>
    }
  />
);

export default NoFoundPage;
