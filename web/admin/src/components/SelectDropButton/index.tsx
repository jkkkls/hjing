import { Button, Divider, Space } from "antd";
import { FormInstance } from "antd/es/form";
import { ReactNode } from "react";

export function genSelectDropButton (menu: ReactNode, channels: { label: string; value: any }[], formRef :React.MutableRefObject<FormInstance<any> | undefined>)  {
    return (
        <>
            <Space style={{ margin: "2px 0 0 8px" }}>
              <Button
                onClick={(e) => {
                  let all = [] as any[];
                  channels?.forEach(e => {
                    all.push(e.value);
                  });
                  formRef.current?.setFieldsValue({Channel: all});
                }}
              >
                全选
              </Button>
              <Button
                onClick={(e) => {
                  formRef.current?.setFieldsValue({Channel: []});
                }}
              >
                清空
              </Button>
              <Button
                onClick={(e) => {
                  let all = [] as any[];
                  const selected = formRef.current?.getFieldValue('Channel') as any[];
                  channels?.forEach(e => {
                    if (selected.includes(e.value)) {
                      return;
                    }
                    all.push(e.value);
                  });
                  formRef.current?.setFieldsValue({Channel: all});
                }}
              >
                反选
              </Button>
            </Space>
            <Divider style={{ margin: "8px 0" }} />
            {menu}
          </>
    );
};
