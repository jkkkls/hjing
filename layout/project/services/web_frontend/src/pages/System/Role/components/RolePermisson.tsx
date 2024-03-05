import React, { useState, useRef } from 'react';
import { message } from 'antd';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import EditableProTable from '@ant-design/pro-table';
import { getRolePermissions, updateRolePermission } from '@/services/ant-design-pro/api';


export type RolePermissionPageProps = {
  values: Partial<any>;
};


const handleUpdate = async (fields: any) => {
  const hide = message.loading('正在添加');
  try {
    let res = await updateRolePermission(fields);
    if (!res.success) {
      hide();
      message.warning(res.errorMessage);
      return false;
    }
    hide();
    return true;
  } catch (error) {
    hide();
    return false;
  }
};



const RolePermissionPage: React.FC<RolePermissionPageProps> = (props) => {
  const { values } = props;

  const getData = async (
    params: {
      current?: number;
      pageSize?: number;
      id?: number;
    },
    options?: { [key: string]: any }
  ) => {
    try {
      console.log(values);
      params.id = values?.id;
      const ret = await getRolePermissions(params, options);
      if (ret.errorMessage && ret.errorMessage !== '') {
        message.warning(ret.errorMessage);
      }
      return ret;
    } catch (error) {
      return [];
    }
  };

  const [editableKeys, setEditableRowKeys] = useState<React.Key[]>([]);
  const actionRef = useRef<ActionType>();
  const columns: ProColumns[] = [
    {
      title: '名称',
      dataIndex: 'name',
      width: 200,
      readonly: true,
    },
    {
      title: '状态',
      key: 'code',
      dataIndex: 'code',
      width: 200,
      valueType: 'checkbox',
      valueEnum: {
        R: "查询",
        U: "更新",
        D: "删除",
      },
    },
    {
      title: '操作',
      valueType: 'option',
      width: 200,
      render: (text, record, _, action) => [
        <a
          key="editable"
          onClick={() => {
            action?.startEditable?.(record.id);
          }}
        >
          编辑
        </a>,
      ],
    },
  ];

  return (
      <EditableProTable
        search={false}
        columns={columns}
        request={getData}
        rowKey="id"
        actionRef={actionRef}
        pagination={false}
        onChange={getData}
        scroll={{ x: 800 }}

        editable={{
          type: 'single',
          editableKeys,
          onSave: async (rowKey, data) => {
            data.sysRoleID = values.id;
            await handleUpdate(data);
            actionRef.current?.reload();
          },
          onChange: setEditableRowKeys,
          actionRender: (row, config, defaultDom) => [
            defaultDom.save,
            defaultDom.cancel,
          ],
        }}
      />
  );
};

export default RolePermissionPage;
