import React from "react";
import { Layout, Image, Menu, theme } from "antd";
const { Header, Content, Footer } = Layout;
import { NavLink, Outlet } from "react-router-dom";
import {
  HomeOutlined,
  MessageOutlined,
  LinkOutlined,
  MobileOutlined,
} from "@ant-design/icons";
import { site } from "./config";
// const items = new Array(20).fill(null).map((_, index) => ({
//   key: index + 1,
//   label: `nav ${index + 1}`,
// }));

const menuItems = [
  {
    label: "首页",
    key: "/home",
    icon: <HomeOutlined />,
  },
  {
    label: "新闻公告",
    key: "/news",
    icon: <MessageOutlined />,
  },
  {
    label: "应用中心",
    key: "/apps",
    icon: <MobileOutlined />,
  },
  {
    label: "关于我们",
    key: "/about",
    icon: <LinkOutlined />,
  },
];

const App = () => {
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  const createFooter = () =>
    site.footerIcons.map((app, index) => (
      // <Image key={index} width={20} src="http://www.52gaoge.com/img/favicon2.png"></Image>粤B2-20191168号
      <a key={index} href={app.url}>
        <Image width={20} src={app.icon}></Image>
        {app.name}
      </a>
    ));

  return (
    <Layout
      style={{
        minHeight: "100vh",
        width: "100vw",
      }}
    >
      <Header
        style={{
          display: "flex",
          alignItems: "center",
          width: "100%",
          zIndex: 1000,
          background: "#fff",
        }}
      >
        {/* <div className="demo-logo" /> */}
        <Image preview={false} width={300} src="/icon.png"></Image>
        <Menu
          // theme="dark"
          mode="horizontal"
          // defaultSelectedKeys={['2']}
          items={menuItems.map((item) => {
            return {
              key: item.key,
              label: <NavLink to={item.key}>{item.label}</NavLink>,
              disabled: item.disabled,
              icon: item.icon,
            };
          })}
          style={{
            flex: 1,
            minWidth: 0,
          }}
        />
      </Header>
      <Content
        style={{
          margin: "0 auto",
        }}
      >
        <div
          style={{
            // background: colorBgContainer,
            // width: '400px',
            //默认宽度1200px，适配屏幕宽。
            width: "100%",

            // height: '100%',
            padding: "10px",
            borderRadius: borderRadiusLG,
            background: "#f2f2f2",
          }}
        >
          <Outlet />
        </div>
      </Content>
      <Footer
        style={{
          textAlign: "center",
        }}
      >
        <p>
          CopyRight ©{new Date().getFullYear()} ALL Right Reserved {site.name}{" "}
          版权所有
        </p>
        <p>{createFooter()}</p>
      </Footer>
    </Layout>
  );
};
export default App;
