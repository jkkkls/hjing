import { Image, Flex, Tabs} from "antd";
import IntroPage from "./intro";
import CulturePage from "./culture";
import ContactPage from "./contact";
import JobPage from "./job";

const divStyle = {
  height: "100%",
  width: "100%",
  // maxWidth: "1200px",

  // display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
  background: "#f2f2f2",
};

const AboutPage = () => {
  return (
    <Flex vertical gap="middle" style={divStyle}>
      <Image preview={false} src={"https://www.52gaoge.com/img/1704263380717.png"}></Image>
      <Tabs
        style={{width : '1280px', minHeight: '400px', border:'1px solid #e5ecff', background:'#ffffff'}}
        tabBarStyle={{background: '#ccccff', width: '15%'}}
        tabPosition={'left'}
        items={[
            {label: '公司简介', key: '1', children: <IntroPage></IntroPage>},
            {label: '企业文化', key: '2', children: <CulturePage></CulturePage>},
            {label: '人才招聘', key: '3', children: <JobPage></JobPage>},
            {label: '联系我们', key: '4', children: <ContactPage></ContactPage>},
        ]}
      />
    </Flex>
  );
};

export default AboutPage;
