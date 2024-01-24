import { defaultApp, getDataApps } from "@/services/ant-design-pro/data";
import { PageContainer } from "@ant-design/pro-layout";
import { Select } from "antd";
import { useEffect, useState } from "react";

export type DataPageContainerProps = {
    onSubmit: (values: string) => Promise<void>;
    children: any;
  };

const DataPageContainer: React.FC<DataPageContainerProps> = (props) => {
  const [data, setData] = useState<any[]>([]);
  const [value, setValue] = useState<string>();

  const getApps = async () => {
    try {
      const ret = await getDataApps({});
      const arr = ret.data || [];
      var data = [];
      for (let i = 0; i < arr.length; i++) {
        const e = arr[i];
        data.push({ text: e.Name || "", value: e.AppId || "" });
      }
      setData(data);
      setValue(ret.appId);
    } catch (error) {
    }
  };

  useEffect(() => {
    getApps();
  }, []);

  const handleChange = async (newValue: string) => {
    defaultApp(newValue);
    setValue(newValue);

    props.onSubmit(newValue);
    //调用自己函数
  };

  const options = data.map((d) => (
    <Select.Option key={d.value}>{d.text}</Select.Option>
  ));

  return (
    <PageContainer
      header={{
        subTitle: (
          <Select style={{ width: 180 }} value={value} onChange={handleChange}>
            {options}
          </Select>
        ),
      }}
      { ...props }
    />
  );
};

export default DataPageContainer;
