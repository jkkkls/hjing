import React, { useState, useRef } from 'react';
import { message } from 'antd';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import EditableProTable from '@ant-design/pro-table';
import { PageContainer } from '@ant-design/pro-layout';
import { getPermissions, updatePermission } from '@/services/ant-design-pro/api';

const handleUpdate = async (fields: any) => {
  const hide = message.loading('正在添加');
  try {
    let res = await updatePermission(fields);
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

const getData = async (options?: { [key: string]: any }) => {
  try {
    const ret = await getPermissions(options);

    if (ret.errorMessage && ret.errorMessage !== '') {
      message.warning(ret.errorMessage);
    }
    return ret;
  } catch (error) {
    return [];
  }
};

const PermissionPage: React.FC = () => {
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
      title: '路径',
      dataIndex: 'path',
      width: 200,
      readonly: true,
    },
    {
      title: '图标',
      dataIndex: 'icon',
      width: 200,
    },
    {
      title: '排序等级',
      dataIndex: 'level',
      width: 200,
      valueType: 'digit',
    },
    {
      title: '状态',
      key: 'status',
      dataIndex: 'status',
      width: 200,
      valueType: 'select',
      valueEnum: {
        normal: {
          text: '正常',
          status: 'Success',
        },
        hide: {
          text: '隐藏',
          status: 'Error',
        },
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
    <PageContainer>
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
    </PageContainer>
  );
};

export default PermissionPage;
