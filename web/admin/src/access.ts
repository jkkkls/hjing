/**
 * @see https://umijs.org/zh-CN/plugins/plugin-access
 * */
export default function (initialState: {
  currentUser?: API.CurrentUser | undefined;
}) {
  const { currentUser } = initialState || {};
  return {
    normalRouteFilter: (route: any) => {
      if (!currentUser || !currentUser.hasRoutes) {
        return false;
      }
      for (const e of currentUser.hasRoutes) {
        if (e == route.path) {
          return true
        }
      }

      return false;
    },
  };
}
