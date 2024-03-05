import { PlusOutlined } from '@ant-design/icons';
import React, { useState, useRef } from 'react';
import { Button, Drawer, message, Modal, Space } from 'antd';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import { PageContainer } from '@ant-design/pro-layout';
import { getRoleList, addRole, delRole } from '@/services/ant-design-pro/api';
import { ModalForm, ProFormText } from '@ant-design/pro-form';
import RolePermissionPage from './components/RolePermisson';

const handleAdd = async (fields: API.RoleItem) => {
  const hide = message.loading('正在添加');
  try {
    await addRole(fields);
    hide();
    return true;
  } catch (error) {
    hide();
    return false;
  }
};

const handleDel = async (id: string) => {
  const hide = message.loading('正在删除');
  try {
    await delRole(id);
    hide();
    return true;
  } catch (error) {
    hide();
    return false;
  }
};

const RoleList: React.FC = () => {
  const [createModalVisible, handleCreateModalVisible] = useState<boolean>(false);
  const [row, setCurrentRow] = useState<API.RoleItem>();
  const [open, setOpen] = useState(false);

  const actionRef = useRef<ActionType>();

  const showDrawer = () => {
    setOpen(true);
  };

  const onClose = () => {
    setOpen(false);
  };

  const columns: ProColumns<API.RoleItem>[] = [
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '描述',
      dataIndex: 'desc',
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
            showDrawer();
          }}
        >
          修改权限
        </a>,
        <a
          key="delete"
          onClick={() => {
            Modal.confirm({
              title: '删除任务',
              content: '确定删除该任务吗？',
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
      <ProTable<API.RoleItem, API.PageParams>
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
                handleCreateModalVisible(true);
              }}
            >
              <PlusOutlined />
              新增角色
            </Button>
          </Space>
        }
        columns={columns}
        request={getRoleList}
        rowKey="id"
        actionRef={actionRef}
      />
      <ModalForm
        title="新建角色"
        width="400px"
        visible={createModalVisible}
        onVisibleChange={handleCreateModalVisible}
        onFinish={async (value) => {
          const success = await handleAdd(value as API.RoleItem);
          if (success) {
            handleCreateModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
        modalProps={{
          destroyOnClose: true,
        }}
      >
        <ProFormText
          rules={[
            {
              required: true,
              message: '角色名称不能为空',
            },
          ]}
          width="md"
          name="name"
          placeholder="请输入角色名称"
        />
        <ProFormText width="md" name="desc" placeholder="请输入角色描述" />
      </ModalForm>
      <Drawer title={"修改角色[" + row?.name + "]权限"} width={1024} onClose={onClose} open={open} destroyOnClose={true}>
        <RolePermissionPage values={row || {}}></RolePermissionPage>
      </Drawer>
    </PageContainer>
  );
};

export default RoleList;
