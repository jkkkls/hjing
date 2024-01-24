import { useEffect, useState } from "react";
import { Image, Flex, List, Skeleton, Typography, Button} from "antd";
const { Paragraph } = Typography;
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

const AppsPage = () => {
  const [list, setList] = useState([]);
  useEffect(() => {
    setList([]);
    const apps = GetApps("xxoo");
    for (let i = 0; i < apps.length; i++) {
      if (!apps[i].downlink) {
        continue;
      }
      const element = { ...apps[i] };
      setList((list) => [...list, element]);
    }
  }, []);

  const createAppDesc = (desc) =>
    desc.map((str, index) => <Paragraph key={index}>{str}</Paragraph>);

  return (
    <Flex vertical gap="middle" style={divStyle}>
      <Image preview={false} src={"banner.png"}></Image>
      <List
        //   className="demo-loadmore-list"
        style={{ minHeight: "400px", width: "100%" }}
        itemLayout="horizontal"
        dataSource={list}
        renderItem={(item) => (
          <List.Item
            style={{
              background: "#f2f2f2",
              paddingLeft: "20px",
              paddingRight: "20px",
            }}
            extra={<Image width={"400px"} alt="logo" src={item.small} />}
          >
            <Skeleton title={false} loading={item.loading} active>
              <List.Item.Meta
                title={
                  <a href={item.app}>
                    <h2>{item.name}</h2>
                  </a>
                }
                description={
                  <Typography>
                    {createAppDesc(item.desc)}
                    <Paragraph>
                      下载应用:{" "}
                      <a rel="noopener noreferrer" href={item.downlink}>
                        <Image
                          width={"40px"}
                          preview={false}
                          src="/android.png"
                        ></Image>
                      </a>
                      <Button style={{ marginLeft: "10px" }} href={item.app}>进入官网</Button>
                    </Paragraph>
                  </Typography>
                }
              />
            </Skeleton>
          </List.Item>
        )}
      />
    </Flex>
  );
};

export default AppsPage;
