import { useEffect, useState } from "react";
import { Typography, Image, Flex, Card, Space } from "antd";
import { GetApp } from "../data/interface";
const { Paragraph } = Typography;

const divStyle = {
  height: "100%",
  maxWidth: "1200px",
  // display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
  background: "#f2f2f2",
};

const AppPage = () => {
  //获取url参数
  const [app, setApp] = useState({ game: [], desc: [] });
  useEffect(() => {
    const url = window.location.href;
    const params = url.split("?")[1];
    const paramsArr = params.split("&");
    const obj = {};
    paramsArr.forEach((item) => {
      const arr = item.split("=");
      obj[arr[0]] = arr[1];
    });
    let data = GetApp("xxoo", obj.id);
    setApp(data);
    console.log(data);
  }, []);

  const createAppDesc = () =>
    app?.desc.map((str, index) => <Paragraph key={index}>{str}</Paragraph>);

  const createAppGame = () =>
    app?.game.map((str, index) => <Image key={index} src={str}></Image>);

  return (
    <Flex vertical gap="middle" style={divStyle}>
      <a href={app.downlink}>
        <Image preview={false} src={app.big}></Image>
      </a>
      <Card
        style={{ width: "100%" }}
        title="应用介绍"
        extra={
          <Space>
            <p>下载方式：</p>
            <a href={app.downlink}>
              <Image width={"40px"} preview={false} src="/android.png"></Image>
            </a>
          </Space>
        }
      >
        <Typography>{createAppDesc()}</Typography>
      </Card>
      <Card title="应用截图" style={{ width: "100%" }}>
        <Flex wrap>{createAppGame()}</Flex>
      </Card>
    </Flex>
  );
};

export default AppPage;
