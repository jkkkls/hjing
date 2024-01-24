import { PlusOutlined, UploadOutlined } from "@ant-design/icons";
import ProForm, {
  ModalForm,
  ProFormGroup,
  ProFormSelect,
  ProFormText,
} from "@ant-design/pro-form";
import { PageContainer } from "@ant-design/pro-layout";
import ProTable, { ActionType, ProColumns } from "@ant-design/pro-table";
import { Alert, Button, Card, Form, message, Modal, Space, Tag, Upload } from "antd";
import { useRef, useState } from "react";
import {
  deleteShared,
  getSharedPageList,
  updateShared,
} from "@/services/ant-design-pro/shared";
import { UploadChangeParam, UploadFile } from "antd/lib/upload/interface";
import Paragraph from "antd/lib/typography/Paragraph";

const handleAdd = async (fields: any) => {
  const hide = message.loading("正在添加");
  try {
    await updateShared(fields);
    hide();
    return true;
  } catch (error) {
    hide();
    return false;
  }
};
const handleDel = async (id: number) => {
  const hide = message.loading("正在删除");
  try {
  await deleteShared(id);
    hide();
    return true;
  } catch (error) {
    hide();
    return false;
  }
};

const SharedModal: React.FC = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [form] = Form.useForm();
  const [disabledAppId, setDisabledAppId] = useState<boolean>(false);

  const [id, setId] = useState<number>();
  const [gameId, setGameId] = useState<string>();
  const [channel, setChannel] = useState<string>();
  const [bgUrl, setBgUrl] = useState<string>();
  const [bannerImgUrl, setBannerImgUrl] = useState<string>();
  const [couponImgUrl, setCouponImgUrl] = useState<string>();
  const [copyImgUrl, setCopyImgUrl] = useState<string>();
  const [tempId, setTempId] = useState<number>();

  const fileList: UploadFile[] = [];

  const columns: ProColumns<any>[] = [
    {
      title: "游戏标示",
      width: 130,
      dataIndex: "GameId",
    },
    {
      title: "渠道标示",
      width: 140,
      dataIndex: "Channel",
    },
    {
      title: "备注",
      width: 200,
      dataIndex: "Desc",
    },
    {
      title: "分享页名",
      width: 200,
      dataIndex: "Title",
    },
    {
      title: "复制链接",
      width: 140,
      hideInSearch: true,
      render: (_, record) => (
        <Space>
          <Paragraph style={{ margin: 4 }} copyable={{ text: record.Domain }}>
            链接
          </Paragraph>
        </Space>
      ),
    },
    {
      title: "操作",
      dataIndex: "option",
      width: 140,
      valueType: "option",
      fixed: "right",
      render: (_: any, record: any) => [
        <a
          key="config"
          onClick={() => {
            setDisabledAppId(true);
            setId(record.ID);
            form.setFieldsValue({
              ...record,
            });
            handleModalVisible(true);
            setGameId(record.GameId);
            setChannel(record.Channel);
          }}
        >
          修改
        </a>,
        <a
          key="channel"
          onClick={() => {
            Modal.confirm({
              title: "删除分享页",
              content: "确定删除该分享页吗？",
              okText: "确认",
              cancelText: "取消",
              onOk: async () => {
                const success = await handleDel(record.ID || 0);
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
  const onChange = (value: number) => {
    setTempId(value);
  };
  return (
    <PageContainer>
      <ProTable
        headerTitle={
          <Space>
            <Button
              type="default"
              key="default"
              onClick={() => {
                form.resetFields();
                setDisabledAppId(false);
                handleModalVisible(true);
              }}
            >
              <PlusOutlined /> 新增数据分享页
            </Button>
          </Space>
        }
        columns={columns}
        request={getSharedPageList}
        actionRef={actionRef}
        rowKey="ID"
        search={false}
        scroll={{ x: 800 }}
      ></ProTable>
      <ModalForm
        title="新增分享页"
        width="1200px"
        visible={createModalVisible}
        onVisibleChange={handleModalVisible}
        form={form}
        preserve={false}
        modalProps={{
          destroyOnClose: true,
        }}
        onFinish={async (value) => {
          value.BgUrl = bgUrl;
          value.BannerImgUrl = bannerImgUrl;
          value.CouponImgUrl = couponImgUrl;
          value.CopyImgUrl = copyImgUrl;
          value.BannerImgUrl = bannerImgUrl;
          value.ID = id
          const success = await handleAdd(value);
          if (success) {
            form.resetFields();
            handleModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormGroup>
          <ProFormText
            rules={[
              {
                required: true,
                message: "游戏标示不能为空",
              },
            ]}
            width="md"
            name="GameId"
            label="游戏标示"
            disabled={disabledAppId}
            fieldProps={{
              onChange: (e) => {
                setGameId(e.target.value);
              },
            }}
          />
          <ProFormText
            rules={[
              {
                required: true,
                message: "渠道标示是必须的",
              },
            ]}
            width="md"
            name="Channel"
            disabled={disabledAppId}
            label="渠道标示"
            fieldProps={{
              onChange: (e) => {
                setChannel(e.target.value);
              },
            }}
          />
        </ProFormGroup>
        <ProFormSelect
          rules={[
            {
              required: true,
              message: "模板是必须的",
            },
          ]}
          label="模板"
          name="TempId"
          width="md"
          fieldProps={{
            options: [
              { label: "无邀请码", value: 0 },
              { label: "带邀请码", value: 1 },
            ],
            onChange: onChange,
          }}
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: "描述是必须的",
            },
          ]}
          width="md"
          name="Desc"
          label="描述"
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: "分享页标题是必须的",
            },
          ]}
          width="md"
          name="Title"
          label="分享页标题"
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: "安卓包下载地址是必须的",
            },
          ]}
          width="lg"
          name="DownloadUrl"
          label="安卓包下载地址"
        />
        <ProForm.Item>
          <Upload
            action={() =>
              "/api/shared/img?gameId=" + gameId + "&channel=" + channel
            }
            maxCount={1}
            listType="picture"
            defaultFileList={[...fileList]}
            onChange={(info: UploadChangeParam) => {
              if (
                info.file &&
                info.file.status == "done" &&
                info.file.response.success == true
              ) {
              setBgUrl(info.file.response.addr);
              }
            }}
          >
            <Button icon={<UploadOutlined />} disabled={tempId != 0 && tempId != 1}>上传背景图(649*1380/750*4106)</Button>
          </Upload>
          </ProForm.Item>
          <ProForm.Item>
          <Upload
            action={() =>
              "/api/shared/img?gameId=" + gameId + "&channel=" + channel
            }
            maxCount={1}
            listType="picture"
            defaultFileList={[...fileList]}
            onChange={(info: UploadChangeParam) => {
              if (
                info.file &&
                info.file.status == "done" &&
                info.file.response.success == true
              ) {
                setCopyImgUrl(info.file.response.addr);
              }
            }}
          >
            <Button icon={<UploadOutlined />} disabled={tempId != 1}>上传复制按钮图(77*161)</Button>
          </Upload>
          </ProForm.Item>
          <ProForm.Item>
          <Upload
            action={() =>
              "/api/shared/img?gameId=" + gameId + "&channel=" + channel
            }
            maxCount={1}
            listType="picture"
            defaultFileList={[...fileList]}
            onChange={(info: UploadChangeParam) => {
              if (
                info.file &&
                info.file.status == "done" &&
                info.file.response.success == true
              ) {
                setCouponImgUrl(info.file.response.addr);
              }
            }}
          >
            <Button icon={<UploadOutlined />} disabled={tempId != 1}>上传红包按钮图(109*344)</Button>
          </Upload>
          </ProForm.Item>
          <ProForm.Item>
          <Upload
            action={() =>
              "/api/shared/img?gameId=" + gameId + "&channel=" + channel
            }
            maxCount={1}
            listType="picture"
            defaultFileList={[...fileList]}
            onChange={(info: UploadChangeParam) => {
              if (
                info.file &&
                info.file.status == "done" &&
                info.file.response.success == true
              ) {
                setBannerImgUrl(info.file.response.addr);
              }
            }}
          >
            <Button icon={<UploadOutlined />} disabled={tempId != 1}>上传Banner按钮图(93*234)</Button>
          </Upload>
          </ProForm.Item>
      </ModalForm>
    </PageContainer>
  );
};

export default SharedModal;
