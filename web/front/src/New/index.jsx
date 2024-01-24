import { useEffect, useState } from "react";
import { Flex } from "antd";
import { Typography } from "antd";

const { Title, Paragraph } = Typography;

const divStyle = {
  height: "100%",
  maxWidth: "1200px",
  // display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
  background: "#f1f1f1",
};

const NewPage = () => {
  //获取url参数
  const [urlParams, setUrlParams] = useState({});
  useEffect(() => {
    const url = window.location.href;
    const params = url.split("?")[1];
    const paramsArr = params.split("&");
    const obj = {};
    paramsArr.forEach((item) => {
      const arr = item.split("=");
      obj[arr[0]] = arr[1];
    });
    setUrlParams(obj);

    console.log(obj);
  }, []);

  return (
    <Flex vertical gap="middle" style={divStyle}>
      <Title>文章{urlParams?.id}的标题。。。。</Title>
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
      <Typography>
        <Paragraph>
          蚂蚁的企业级产品是一个庞大且复杂的体系。这类产品不仅量级巨大且功能复杂，而且变动和并发频繁，常常需要设计与开发能够快速的做出响应。同时这类产品中有存在很多类似的页面以及组件，可以通过抽象得到一些稳定且高复用性的内容。
        </Paragraph>

        <Paragraph>
          随着商业化的趋势，越来越多的企业级产品对更好的用户体验有了进一步的要求。带着这样的一个终极目标，我们（蚂蚁集团体验技术部）经过大量的项目实践和总结，逐步打磨出一个服务于企业级产品的设计体系
          Ant Design。基于『确定』和『自然』
          的设计价值观，通过模块化的解决方案，降低冗余的生产成本，让设计者专注于
          更好的用户体验。
        </Paragraph>
      </Typography>
    </Flex>
  );
};

export default NewPage;
