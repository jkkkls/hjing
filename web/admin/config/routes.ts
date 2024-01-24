export default [
  {
    path: "/user",
    layout: false,
    routes: [
      {
        path: "/user",
        routes: [
          { name: "登录", path: "/user/login", component: "./user/Login" },
        ],
      },
      { component: "./404" },
    ],
  },
  { path: "/welcome", icon: "smile", name: "欢迎界面", component: "./Welcome" },
  {
    path: "/system",
    name: "系统管理",
    icon: "setting",
    access: "normalRouteFilter",
    routes: [
      {
        path: "/system/user",
        name: "用户管理",
        component: "./System/User",
        access: "normalRouteFilter",
      },
      {
        path: "/system/role",
        name: "角色管理",
        component: "./System/Role",
        access: "normalRouteFilter",
      },
      {
        path: "/system/log",
        name: "后台日志",
        component: "./System/Log",
        access: "normalRouteFilter",
      },
      { component: "./404" },
    ],
  },
  { path: "/", redirect: "/welcome" },
  { component: "./404" },
];
