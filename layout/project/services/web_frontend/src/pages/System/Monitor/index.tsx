import { getMonitorList } from "@/services/ant-design-pro/api";
import { PageContainer } from "@ant-design/pro-layout";
import ProTable, { ActionType, ProColumns } from "@ant-design/pro-table";
import Text from "antd/lib/typography/Text";
import moment from "moment";
import { useRef } from "react";


const EventNodeModal: React.FC = () => {
  const actionRef = useRef<ActionType>();

  const columns: ProColumns<any>[] = [
    {
      title: "节点名称",
      width: 70,
      dataIndex: "name",
      hideInSearch: true,
    },
    {
      title: "进程ID",
      width: 80,
      render: (_, record) => <Text>{record.pid}</Text>,
      hideInSearch: true,
    },
    {
      title: "CPU占用率",
      width: 80,
      render: (_, record) => <Text>{record.CPUPercent}</Text>,
      hideInSearch: true,
    },
    {
      title: "内存",
      width: 80,
      render: (_, record) => <Text>{record.memory}KB</Text>,
      hideInSearch: true,
    },
    {
      title: "协程数",
      width: 80,
      render: (_, record) => <Text>{record.numGoroutine}</Text>,
      hideInSearch: true,
    },
    {
      title: "Git标签",
      ellipsis: true,
      width: 80,
      render: (_, record) => <Text>{record.GitTag}</Text>,
      hideInSearch: true,
    },
    {
      title: "编译机器",
      ellipsis: true,
      width: 100,
      render: (_, record) => <Text>{record.PcName}</Text>,
      hideInSearch: true,
    },
    {
      title: "编译时间",
      ellipsis: true,
      width: 120,
      render: (_, record) => <Text>{record.BuildTime}</Text>,
      hideInSearch: true,
    },
    {
      title: "版本",
      ellipsis: true,
      width: 100,
      render: (_, record) => <Text>{record.GitSHA}</Text>,
      hideInSearch: true,
    },
    {
      title: "运行时间",
      width: 120,
      render: (_, record) => moment(record.runTime*1000).format("YYYY-MM-DD HH:mm:ss"),
      hideInSearch: true,
    },
    {
      title: "更新时间",
      width: 120,
      render: (_, record) => moment(record.time*1000).format("YYYY-MM-DD HH:mm:ss"),
      hideInSearch: true,
    },
  ];

  return (
    <PageContainer>
      <ProTable<any, API.PageParams>
        search={false}
        pagination={false}
        columns={columns}
        request={getMonitorList}
        actionRef={actionRef}
        rowKey="Name"
      ></ProTable>
    </PageContainer>
  );
};

export default EventNodeModal;
