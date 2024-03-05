import { PlusOutlined } from '@ant-design/icons';
import React, { useState, useRef } from 'react';
import { Button, message, Modal, Space } from 'antd';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import { PageContainer } from '@ant-design/pro-layout';
import { getUserList, addUser, delUser } from '@/services/ant-design-pro/api';
import AddUserForm, { getRoles } from './components/AddUser';
import moment from 'moment';

const handleAdd = async (fields: API.UserItem) => {
  const hide = message.loading('正在添加');
  try {
    await addUser(fields);
    hide();
    return true;
  } catch (error) {
    hide();
    return false;
  }
};

const handleDel = async (id: string) => {
  const hide = message.loading('正在添加');
  try {
    await delUser(id);
    hide();
    return true;
  } catch (error) {
    hide();
    return false;
  }
};

const UserList: React.FC = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [row, setCurrentRow] = useState<API.UserItem>();
  const actionRef = useRef<ActionType>();
  const columns: ProColumns<API.UserItem>[] = [
    {
      title: '用户名',
      dataIndex: 'username',
    },
    {
      title: '姓名',
      dataIndex: 'name',
    },
    {
      title: '角色',
      dataIndex: 'role',
      valueType: 'select',
      params: {},
      request: getRoles,
    },
    {
      title: '创建时间',
      dataIndex: 'createTs',
      hideInSearch: true,
      render: (_, record) =>
        record.createTs ? moment(record.createTs).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
    {
      title: '最近登录时间',
      dataIndex: 'lastLoginTs',
      hideInSearch: true,
      render: (_, record) =>
        record.lastLoginTs ? moment(record.lastLoginTs).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      width: 240,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setCurrentRow(record);
            console.log(record);
            handleModalVisible(true);
          }}
        >
          修改用户
        </a>,
        <a
          key="delete"
          onClick={() => {
            Modal.confirm({
              title: '删除任务',
              content: '确定删除该用户吗？',
              okText: '确认',
              cancelText: '取消',
              onOk: async () => {
                const success = await handleDel(record.id || '');
                if (success) {
                  if (actionRef.current) {
                    actionRef.current.reload();
                  }
                }
              },
            });
          }}
        >
          删除
        </a>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.UserItem, API.PageParams>
        search={{
          optionRender: false,
          collapsed: false,
        }}
        headerTitle={
          <Space>
            <Button
              type="default"
              key="default"
              onClick={() => {
                setCurrentRow(undefined);
                handleModalVisible(true);
              }}
            >
              <PlusOutlined /> 新增用户
            </Button>
          </Space>
        }
        columns={columns}
        request={getUserList}
        rowKey="id"
        actionRef={actionRef}
      />
      <AddUserForm
        onSubmit={async (value) => {
          console.log('commit', value);
          const success = await handleAdd(value as API.UserItem);
          if (success) {
            handleModalVisible(false);
            setCurrentRow(undefined);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
        onCancel={() => {
          handleModalVisible(false);
          setCurrentRow(undefined);
        }}
        updateModalVisible={createModalVisible}
        values={row || {}}
      />
    </PageContainer>
  );
};

export default UserList;
