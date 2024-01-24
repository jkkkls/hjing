import { useEffect, useState } from "react";
import { Flex } from "antd";
import { Typography } from "antd";
import GetPrivacyContent from "../data/privacy";
import { GetApp } from "../data/interface";

const { Title } = Typography;

const divStyle = {
  height: "100%",
  maxWidth: "1200px",
  // display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
  background: "#f1f1f1",
};

const PrivacyPage = () => {
  //获取url参数
  const [name, setName] = useState("");
  useEffect(() => {
    const url = window.location.href;
    const params = url.split("?")[1];
    const paramsArr = params.split("&");
    const obj = {};
    paramsArr.forEach((item) => {
      const arr = item.split("=");
      obj[arr[0]] = arr[1];
    });

    //
    // console.log(obj);
    // console.log(GetApp('xxoo', obj.id));
    setName(GetApp("xxoo", obj.id)?.name);
  }, []);

  return (
    <Flex vertical gap="middle" style={divStyle}>
      <Title>{name}隐私文件</Title>
      {/* <Paragraph>
        <Text>
          发布时间：<Text strong>2023年1月1日</Text>{"  "}
        </Text>
        <Text>
          作者：<Text strong>淘小宝 </Text>{"  "}
        </Text>
        <Text>
          阅读：<Text strong>123456 </Text>
        </Text>
      </Paragraph> */}
      {GetPrivacyContent(name)}
    </Flex>
  );
};

export default PrivacyPage;
