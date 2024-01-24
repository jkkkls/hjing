import { Typography, Image, Space } from "antd";
const { Paragraph } = Typography;
import { site } from "../config";

const IntroPage = () => {
  return (
    <div style={{ width: "90%", margin: "20px" }}>
      <Typography>
        <Paragraph>
          {site.name}成立于2018年，是一家软件运营服务提供商，拥有一流的技术团队。
        </Paragraph>
        <Paragraph>
          公司以“数据智能，驱动业务增长”为使命，基于卓越的技术与算法能力，结合实时更新的全域数据资源，通过Al赋能的一站式互联网数据产品与运营服务体系，帮助游戏企业实现深度用户洞察、实时业务决策和持续业务增长。
        </Paragraph>
        <Paragraph>
          未来，将技术和业务能力充分拉通，实现产业数字化，不仅要做数据中台，更要把数据中台的灵活使用渗透到实际业务中。
        </Paragraph>
      </Typography>
      <Space>
        <Image src={"https://www.52gaoge.com/img/1637660137570.jpg"}></Image>
        <Image src={"https://www.52gaoge.com/img/1637660265512.jpg"}></Image>
        <Image src={"https://www.52gaoge.com/img/1637660218414.jpg"}></Image>
      </Space>
    </div>
  );
};

export default IntroPage;
