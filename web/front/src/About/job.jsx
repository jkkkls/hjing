import { Space, Collapse } from "antd";
import { Typography } from "antd";
const { Text } = Typography;

const text1 = [
  '工作职责：',
  '1、根据应用特点，利用微信、QQ、论坛以及各大社交平台运营推广我们的产品；',
  '2、收集市场信息，了解市场需求，通过数据分析和信息反馈做出合理有效的运营推广和实施方案；',
  '3、持续关注竞争对手，做竞品调查与分析；',
  '4、对于线上活动运营、新媒体运营有着深入的了解，与各渠道商和合作伙伴保持良好的沟通，建立良好的长期合作关系；',
  '5、及时反馈市场情况，定期总结产品状况，将问题和建议反馈给公司；',
  '6、做过金融、保险类者都可以。',

  '任职要求：',
  '1、口齿清晰，普通话流利，语言富有感染力；',
  '2、专科以上学历，如优秀者不限学历，2年以上网络行业相关经验；',
  '3、对于线上活动运营、新媒体运营有着深入的了解，有一些资源优先；',
  '4、认真踏实，吃苦耐劳，目标感强。',
];

const text2 = [
  '岗位职责：',
  '1，服务端构架设计、模块划分；',
  '2，开发、自测业务逻辑，交由测试验收反馈；',
  '3，协调部门业务链接；',
  '4，熟悉后台开发。',

  '任职要求：',
  '1，对项目有绝对的责任感和团队协作精神；',
  '2，3年以上服务端开发经验；',
  '3，熟悉网络协议HTTP、TCP；linux系统；',
  '4，精通C/C++和lua；以及常用数据库开发与优化；熟悉实时内存数据库如redis、MongoDB等；',
  '5，有良好的编程思想；',
];

const JobPage = () => {
  return (
    <div style={{ width: "100%", margin: "20px" }}>
      <Space direction={"vertical"} style={{ width: "70%"}}>
        <Collapse

          items={[
            {
              key: "1",
              label: <Text strong>线上运营</Text>,
              children: text1.map((item, i) => (<><Text key={i}>{item}</Text><br/></>)),
            },
          ]}
        />
        <Collapse
          items={[
            {
              key: "2",
              label: <Text strong>服务端开发主程</Text>,
              children: text2.map((item, i) => (<><Text key={i}>{item}</Text><br/></>)),
            },
          ]}
        />
      </Space>
    </div>
  );
};

export default JobPage;
