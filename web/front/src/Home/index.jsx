import { Carousel, Flex, Card, Space, Image } from "antd";
const { Meta } = Card;
import { GetApps } from "../data/interface";
import { site } from "../config";

const contentStyle = {
  height: "400px",
  color: "#ffffff",
  lineHeight: "400px",
  textAlign: "center",
  background: "#364d79",
};

const carouselStyle = {
  height: "400px",
  width: "1200px",
  //最大宽度1200px，适配屏幕的宽
  maxWidth: "1200px",
  // width: "100%",
  // 居中
  margin: "0 auto",
};

const divStyle = {
  height: "100%",
  width: "100%",
  // display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
  background: "#f2f2f2",
};

const HomePage = () => {
  const apps = GetApps("xxoo");

  const createCarousel = () =>
    apps.map(
      (app, index) =>
        app.big && (
          <div key={index}>
            <h3 style={contentStyle}>
              <a href={app.app}>
                <img
                  src={app.big}
                  style={{ width: "100%", height: "100%" }}
                ></img>
              </a>
            </h3>
          </div>
        )
    );

  const createPrivacy = (app) =>
    app.privacy && (
      <p style={{ width: "calc(100% - 80px)" }}>
        <a href={app.privacy}>{app.name}隐私文件</a>
      </p>
    );

  const createUserfile = (app) =>
    app.userfile && (
      <p style={{ width: "calc(100% - 80px)" }}>
        <a href={app.userfile}>{app.name}用户文件</a>
      </p>
    );

  const createNews = () =>
    apps.map(
      (app, index) =>
        index <= 1 && (
          <div key={index}>
            {createPrivacy(app)} {createUserfile(app)}
          </div>
        )
    );

  const createHotApps = () =>
    apps.map(
      (app, index) =>
        app.small && (
          <a key={index} href={app.app}>
            <Card
              hoverable
              style={{
                width: 240,
              }}
              cover={<Image preview={false} alt="example" src={app.small} />}
            >
              <Meta title={app.name} description={app.title} />
            </Card>
          </a>
        )
    );

  return (
    <Flex vertical gap="middle" style={divStyle}>
      <Carousel autoplay style={carouselStyle}>
        {createCarousel()}
      </Carousel>
      <Flex
        gap="middle"
        wrap
        style={{
          width: "100%",
          marginTop: "20px",
        }}
      >
        <Card
          title="新闻公告"
          bordered={false}
          extra={<a href="/news">更多</a>}
          style={{ width: "33.3%" }}
          bodyStyle={{ marginTop: "-20px" }}
        >
          {createNews()}
        </Card>
        <Card
          title="招兵买马"
          bordered={false}
          extra={<a href="/about">更多</a>}
          style={{ width: "33.3%" }}
          bodyStyle={{ marginTop: "-20px" }}
        >
          <p>线上运营</p>
          <p>服务端开发主程</p>
        </Card>
        <Card
          title="公司简介"
          bordered={false}
          extra={<a href="/about">详细</a>}
          style={{ width: "33.3%" }}
          bodyStyle={{ marginTop: "-20px" }}
        >
          <p>{site.name}成立于2018...</p>
        </Card>
      </Flex>
      <Card title="热门应用" style={{ width: "100%" }} bordered={false}>
        <Space wrap={true}>{createHotApps()}</Space>
      </Card>
    </Flex>
  );
};

export default HomePage;
