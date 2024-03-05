import { getLogList } from "@/services/ant-design-pro/api";
import { PageContainer } from "@ant-design/pro-layout";
import ProTable, { ProColumns } from "@ant-design/pro-table";
import moment from "moment";

const LogList: React.FC = () => {
  const columns: ProColumns[] = [
    {
      title: "用户名",
      width: 140,
      dataIndex: "name",
    },
    {
      title: "操作",
      width: 140,
      dataIndex: "operate",
    },
    {
      title: "时间",
      width: 180,
      dataIndex: "time",
      hideInSearch: true,
      render: (_, record) => moment(record.time).format("YYYY-MM-DD HH:mm:ss"),
    },
    {
      title: "操作内容",
      ellipsis: true,
      hideInSearch: true,
      dataIndex: "data",
    },
  ];

  return (
    <PageContainer>
      <ProTable columns={columns} request={getLogList} rowKey="name"></ProTable>
    </PageContainer>
  );
};

export default LogList;
