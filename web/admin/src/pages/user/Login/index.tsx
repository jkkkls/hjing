import { LockOutlined, UserOutlined } from "@ant-design/icons";
import { Alert, message, Space } from "antd";
import React, { useState } from "react";
import { ProFormText, LoginForm } from "@ant-design/pro-form";
import { history, useModel } from "@umijs/max";
import Footer from "@/components/Footer";
import { login } from "@/services/ant-design-pro/api";

import styles from "./index.less";
import { flushSync } from "react-dom";

const LoginMessage: React.FC<{
  content: string;
}> = ({ content }) => (
  <Alert
    style={{
      marginBottom: 24,
    }}
    message={content}
    type="error"
    showIcon
  />
);

const waitTime = (time: number = 100) => {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(true);
    }, time);
  });
};

const Login: React.FC = () => {
  const [userLoginState, setUserLoginState] = useState<API.LoginResult>({});
  const [type] = useState<string>("account");
  const { initialState, setInitialState } = useModel("@@initialState");

  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      // await setInitialState((s) => ({
      //   ...s,
      //   currentUser: userInfo,
      // }));
      flushSync(() => {
        setInitialState((s) => ({
          ...s,
          currentUser: userInfo,
        }));
      });
    }
  };

  const handleSubmit = async (values: API.LoginParams) => {
    try {
      // 登录
      await waitTime(1000);
      const msg = await login({ ...values, type });
      if (msg.status == "ok") {
        localStorage.setItem("token", msg.token || "");
        message.success("登录成功！");
        await fetchUserInfo();
        /** 此方法会跳转到 redirect 参数所在的位置 */
        const urlParams = new URL(window.location.href).searchParams;
        history.push(urlParams.get("redirect") || "/");
        return;
      }
      // 如果失败去设置用户错误信息
      message.error(msg.status);
      setUserLoginState(msg);
    } catch (error) {
      const defaultLoginFailureMessage = "登录失败，请重试！";
      console.log(error);
      message.error(defaultLoginFailureMessage);
    }
  };
  const { status } = userLoginState;

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <LoginForm
          logo={<img alt="logo" src="/logo.svg" />}
          title="高歌中台管理系统"
          subTitle="(❁´‿`❁)*✲ﾟ*"
          initialValues={{
            autoLogin: true,
          }}
          onFinish={async (values) => {
            await handleSubmit(values as API.LoginParams);
          }}
        >
          {status && status != "ok" && <LoginMessage content={status} />}
          {type === "account" && (
            <>
              <ProFormText
                name="username"
                fieldProps={{
                  size: "large",
                  prefix: <UserOutlined className={styles.prefixIcon} />,
                }}
                placeholder="用户名:"
                rules={[
                  {
                    required: true,
                    message: "请输入用户名",
                  },
                ]}
              />
              <ProFormText.Password
                name="password"
                fieldProps={{
                  size: "large",
                  prefix: <LockOutlined className={styles.prefixIcon} />,
                }}
                placeholder="密码:"
                rules={[
                  {
                    required: true,
                    message: "请输入密码！",
                  },
                ]}
              />
            </>
          )}
        </LoginForm>
      </div>
      <Footer />
    </div>
  );
};

export default Login;
