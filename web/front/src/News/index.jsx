import { useEffect, useState } from "react";
import { Space, Tag, Image, Flex, List } from "antd";
import { GetApps } from "../data/interface";

const divStyle = {
  height: "100%",
  width: "100%",
  // display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
  background: "#f2f2f2",
};

const NewsPage = () => {
  const [list, setList] = useState([]);
  useEffect(() => {
    setList([]);
    const apps = GetApps("xxoo");
    for (let i = 0; i < apps.length; i++) {
      const element = { ...apps[i] };
      setList((list) => [...list, element]);
    }
  }, []);

  const createPrivacy = (item) =>
    item.privacy && (
      <List.Item
        style={{
          background: "#f2f2f2",
          paddingLeft: "20px",
          paddingRight: "20px",
        }}
        extra={<p>{item.date}</p>}
      >
        <Space>
          <Tag color="warning">公告</Tag>
          <a href={item.privacy}>【{item.name}】隐私文件</a>
        </Space>
      </List.Item>
    );

  const createUserfile = (item) =>
    item.userfile && (
      <List.Item
        style={{
          background: "#f2f2f2",
          paddingLeft: "20px",
          paddingRight: "20px",
        }}
        extra={<p>{item.date}</p>}
      >
        <Space>
          <Tag color="warning">公告</Tag>
          <a href={item.userfile}>【{item.name}】用户文件</a>
        </Space>
      </List.Item>
    );
  return (
    <Flex vertical gap="middle" style={divStyle}>
      <Image preview={false} src={"banner.png"}></Image>
      <List
        //   className="demo-loadmore-list"
        style={{ minHeight: "400px", width: "100%" }}
        itemLayout="horizontal"
        dataSource={list}
        renderItem={(item) => (
          <>
            {createPrivacy(item)}
            {createUserfile(item)}
          </>
        )}
      />
    </Flex>
  );
};

export default NewsPage;
