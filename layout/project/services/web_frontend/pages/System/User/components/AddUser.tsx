import React from 'react';
import { message, Form } from 'antd';
import { ProForm, ModalForm, ProFormCheckbox, ProFormText } from '@ant-design/pro-form';
import { getAllRoles } from '@/services/ant-design-pro/api';

export type FormValueType = {
  target?: string;
  template?: string;
  type?: string;
  time?: string;
  frequency?: string;
} & Partial<API.UserItem>;

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: FormValueType) => void;
  onSubmit: (values: FormValueType) => Promise<void>;
  updateModalVisible: boolean;
  values: Partial<API.UserItem>;
};

export const getRoles = async () => {
  try {
    const ret = await getAllRoles({});
    const arr = ret.data || [];
    let data = [];
    for (let i = 0; i < arr.length; i++) {
      const e = arr[i];
      data.push({ label: e.name || '', value: e.id || '' });
    }
    return data;
  } catch (error) {
    return [];
  }
};

const AddUserForm: React.FC<UpdateFormProps> = (props) => {
  const [form] = Form.useForm();
  return (
    <ModalForm<{
      id: string;
    }>
      title={props.values?.id ? '修改用户' : '新建用户'}
      autoFocusFirstInput
      form={form}
      modalProps={{
        destroyOnClose: true,
        onCancel: () => {
          props.onCancel();
        },
      }}
      visible={props.updateModalVisible}
      onVisibleChange={async (visible) => {
        if (visible) {
          form.setFieldsValue(props.values);
        } else {
          form.resetFields();
        }
      }}
      onFinish={async (values) => {
        message.success('提交成功');

        return props.onSubmit(values);
      }}
    >
      <ProForm.Group>
      <ProFormText
          name="id"
          hidden= {true}
        />
        <ProFormText
          width="md"
          name="username"
          label="用户名"
          placeholder="请输入名称"
          // initialValue= {props.values?.id}
        />
      </ProForm.Group>
      <ProForm.Group>
        <ProFormText width="md" name="name" label="姓名" placeholder="请输入姓名" />
      </ProForm.Group>
      <ProForm.Group>
        <ProFormText.Password width="md" name="password" label="密码" />
      </ProForm.Group>
      <ProForm.Group>
        <ProFormCheckbox.Group
          width="xs"
          params={{}}
          request={getRoles}
          name="sysRoleIDs"
          label="角色"
          fieldProps={{
            // defaultValue:props.values?.role,
            style: { width: 240 },
          }}
        />
      </ProForm.Group>
    </ModalForm>
  );
};

export default AddUserForm;
