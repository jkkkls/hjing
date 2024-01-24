import React, { useState } from "react";
import { Tree, message, Form } from "antd";
import { DrawerForm } from "@ant-design/pro-form";
import { getMenu } from "@/services/ant-design-pro/api";
import { DataNode } from "antd/lib/tree";
import { ApiOutlined, ApiTwoTone, SmileOutlined } from "@ant-design/icons";

export const getMenus = async () => {
  const ret = await getMenu();
  return ret || [];
};

export type FormValueType = {
  target?: string;
  template?: string;
  type?: string;
  time?: string;
  frequency?: string;
} & Partial<API.RoleItem>;

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: FormValueType) => void;
  onSubmit: (values: FormValueType) => Promise<void>;
  updateModalVisible: boolean;
  values: Partial<API.RoleItem>;
};

const EditRole: React.FC<UpdateFormProps> = (props) => {
  const { values } = props;

  const [checkedKeys, setCheckedKeys] = useState<string[]>([]);
  const [selectedKeys, setSelectedKeys] = useState([]);
  const [autoExpandParent] = useState(true);
  const [treeNodes, setTreeNodes] = useState<DataNode[]>([]);

  const onCheck = (checkedKeysValue: any) => {
    setCheckedKeys(checkedKeysValue);
  };

  const onSelect = (selectedKeysValue: any, info: any) => {
    setSelectedKeys(selectedKeysValue);
  };

  const fillChildred = (node: DataNode, arr: API.MenuItem[]) => {
    arr.forEach((e) => {
      if (e.hide == 1) {
        return
      }
      let n = {
        key: e.key || "",
        title: e.title || "",
        children: [],
      } as DataNode;

      if ((n.key as string).includes("@")) {
        n.title = n.title+"[" + e.key +"]";
        n.icon = <ApiTwoTone twoToneColor="#eb2f96" />
      }

      if (e.children && e.children?.length > 0) {
        fillChildred(n, e.children);
      }

      node.children?.push(n);
    });
  };

  return (
    <>
      {!values.id ? null : (
        <DrawerForm
          request={async () => {
            const ret = await getMenu();
            const arr = ret.data || [];
            var data: DataNode[] = [];
            arr.forEach((e) => {
              let n = {
                key: e.key || "",
                title: e.title || "",
                children: [],
              } as DataNode;

              if (e.children && e.children?.length > 0) {
                fillChildred(n, e.children);
              }

              data.push(n);
            });
            setTreeNodes(data);
            setCheckedKeys(values.selected || []);

            let d: API.RoleItem = {};
            return d;
          }}
          visible={props.updateModalVisible}
          width={800}
          onFinish={async (values) => {
            message.success("提交成功");
            values.selected = (checkedKeys as any).checked;
            values.id = props.values.id;
            values.name = props.values.name;
            return props.onSubmit(values);
          }}
          drawerProps={{
            title: "修改角色[" + values.name + "]权限",
            destroyOnClose: true,
            onClose: () => {
              props.onCancel();
            },
          }}
        >
          <Tree<DataNode>
            checkStrictly //会附带比较复杂的数据结构
            showIcon
            defaultExpandAll
            checkable
            autoExpandParent={autoExpandParent}
            onCheck={onCheck}
            checkedKeys={checkedKeys}
            onSelect={onSelect}
            selectedKeys={selectedKeys}
            treeData={treeNodes}
          />
        </DrawerForm>
      )}
    </>
  );
};

export default EditRole;
