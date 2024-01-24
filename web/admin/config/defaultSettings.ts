import { Settings as LayoutSettings } from '@ant-design/pro-layout';

const Settings: LayoutSettings & {
  pwa?: boolean;
  logo?: string;
} = {
  navTheme: "light",
  "layout": "mix",
  contentWidth: "Fluid",
  fixedHeader: false,
  fixSiderbar: true,
  pwa: false,
  logo: "/logo.svg",
  splitMenus: false,
  title: "中台系统",
};

export default Settings;
