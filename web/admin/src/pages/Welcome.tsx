import React from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import { Card, Alert, Image } from 'antd';
// import styles from './Welcome.less';

// const CodePreview: React.FC = ({ children }) => (
//   <pre className={styles.pre}>
//     <code>
//       <Typography.Text copyable>{children}</Typography.Text>
//     </code>
//   </pre>
// );

const Welcome: React.FC = () => {
  return (
    <PageContainer>
      <Card>
        <Alert
          message='欢迎欢迎'
          type="success"
          showIcon
          banner
          style={{
            margin: -12,
            marginBottom: 24,
          }}
        />
      </Card>
    </PageContainer>
  );
};

export default Welcome;
