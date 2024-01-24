// @ts-ignore
/* eslint-disable */

declare namespace API {
  type CurrentUser = {
    name?: string;
    avatar?: string;
    userid?: string;
    email?: string;
    signature?: string;
    title?: string;
    group?: string;
    tags?: { key?: string; label?: string }[];
    notifyCount?: number;
    unreadCount?: number;
    country?: string;
    access?: string;
    geographic?: {
      province?: { label?: string; key?: string };
      city?: { label?: string; key?: string };
    };
    address?: string;
    phone?: string;
    hasRoutes?: string[];
    role?: string;
  };

  type LoginResult = {
    status?: string;
    type?: string;
    currentAuthority?: string;
    token?: string;
  };

  type PageParams = {
    current?: number;
    pageSize?: number;
  };

  type RuleListItem = {
    key?: number;
    disabled?: boolean;
    href?: string;
    avatar?: string;
    name?: string;
    owner?: string;
    desc?: string;
    callNo?: number;
    status?: number;
    updatedAt?: string;
    createdAt?: string;
    progress?: number;
  };

  type GameItem = {
    appId?: string;
    name?: string;
    secret?: string;
    dbName?: string;

    wxAppId?: string;
    wxKey?: string;
    qttAppId?: string;
    qttSecret?: string;

    moneyCheck?: boolean;
    moneyCheckAmount?: number;
    bindNum?: number;
    submitLimit?: number;
    transLimit?: number;
    amountLimit?: number;

    moneyCheckAmount2?: number;
    amountLimit2?: number;
  };
  type GameItemList = {
    data?: GameItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type RedChannel = {
    ID?: number;
    appId?: string;
    name?: string;
    appName?: string;
    channelMin?: number;
    channelMax?: number;
    wxAppId?: string;
    wxKey?: string;
  };
  type RedChannelList = {
    data?: RedChannel[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type RuleList = {
    data?: RuleListItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type FakeCaptcha = {
    code?: number;
    status?: string;
  };

  type LoginParams = {
    username?: string;
    password?: string;
    autoLogin?: boolean;
    type?: string;
  };

  type Log = {
    name?: string;
    time?: string;
    operate?: string;
    data?: string;
  };
  type LogList = {
    data?: Log[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type ErrorResponse = {
    /** 业务约定的错误码 */
    errorCode: string;
    /** 业务上的错误信息 */
    errorMessage?: string;
    /** 业务上的请求是否成功 */
    success?: boolean;
  };

  type NoticeIconList = {
    data?: NoticeIconItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type NoticeIconItemType = "notification" | "message" | "event";

  type NoticeIconItem = {
    id?: string;
    extra?: string;
    key?: string;
    read?: boolean;
    avatar?: string;
    title?: string;
    status?: string;
    datetime?: string;
    description?: string;
    type?: NoticeIconItemType;
  };

  type UserItem = {
    id?: string;
    name?: string;
    password?: string;
    role?: string;
    roleName?: string;
    createTs?: stirng;
    lastLoginTs?: string;
  };

  type MenuItem = {
    key?: string;
    title?: string;
    hide?: number;
    children?: MenuItem[];
  };
  type MenuList = {
    data?: MenuItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };
  type RoleItem = {
    id?: string;
    name?: string;
    desc?: string;
    selected?: string[];
  };
  type RoleList = {
    data?: RoleItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };
  type UserList = {
    data?: UserItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type PublicSettings = {
    moneyCheck?: boolean;
    moneyCheckAmount?: number;
    bindNum?: number;
    rejectBlack?: boolean;
    black: string[];
    wxSubmitLimit?: number;
    wxRedLimit?: number;
    wxTransLimit?: number;
    wxAmountLimit?: number;
    aliTransLimit?: number;
    aliAmountLimit?: number;
    aliAccountCheck?: boolean;
  };

  type WithdrawRecord = {
    CreatedAt?: string;
    Uid?: number;
    AppId?: string;
    PlatformId: string;
    OrderId?: string;
    Channel?: number;
    WdType?: string;
    OpenId?: string;
    Name?: string;
    Money?: number;
    SubmitTime?: string;
    FinishTime?: string;
    Status?: number;
    StatusMsg?: string;
    SubmitRetCode?: number;
    Msg?: string;
    Ip?: string;
    CheckTime?: string;
    CheckUser?: string;
    PlatformCode?: string;
    PlatformMsg?: string;
    IsTest?: number;
  };
  type WithdrawRecordList = {
    data?: WithdrawRecord[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };
}
