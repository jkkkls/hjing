import { DefaultFooter } from "@ant-design/pro-layout";

const Footer: React.FC = () => {
  const defaultMessage = "XXOO网络";
  const currentYear = new Date().getFullYear();
  return (
    <DefaultFooter
      style={{
        background: "none",
      }}
      copyright={`${currentYear} ${defaultMessage}`}
    />
  );
};

export default Footer;
