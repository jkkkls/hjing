import type { Settings as LayoutSettings } from '@ant-design/pro-layout';
import { PageLoading } from '@ant-design/pro-layout';
import type { RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import { history } from '@umijs/max';
import RightContent from '@/components/RightContent';
import Footer from '@/components/Footer';
import { currentUser as queryCurrentUser, getMenu } from './services/ant-design-pro/api';
import defaultSettings from '../config/defaultSettings';
import { ConfigProvider } from 'antd';


const loginPath = '/user/login';

/**
 * @see  https://umijs.org/zh-CN/plugins/plugin-initial-state
 * */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: API.CurrentUser;
  loading?: boolean;
  fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
}> {
  const fetchUserInfo = async () => {
    try {
      const msg = await queryCurrentUser();
      return msg.data;
    } catch (error) {
      history.push(loginPath);
    }
    return undefined;
  };
  // 如果是登录页面，不执行
  if (history.location.pathname !== loginPath) {
    const currentUser = await fetchUserInfo();
    return {
      fetchUserInfo,
      currentUser,
      settings: defaultSettings,
    };
  }
  return {
    fetchUserInfo,
    settings: defaultSettings,
  };
}

const authHeaderInterceptor = (url: string, options: any) => {
  const authHeader = { Authorization: 'Bearer ' + localStorage.getItem('token') };
  return {
    url: `${url}`,
    options: { ...options, interceptors: true, headers: authHeader },
  };
};
const responseInterceptors = (response: Response) => {
  return response;
};
export const request: RequestConfig = {
  // credentials: 'include',
  requestInterceptors: [authHeaderInterceptor],
  responseInterceptors: [(response) => {return response;}],
  // responseInterceptors: [responseInterceptors],
};


// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({ initialState, setInitialState }) => {
  return {
    rightContentRender: () => <RightContent />,
    disableContentMargin: false,
    // waterMarkProps: {
    //   content: initialState?.currentUser?.name,
    // },
    footerRender: () => <Footer />,
    onPageChange: () => {
      const { location } = history;
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.push(loginPath);
      }
    },
    menuHeaderRender: undefined,
    ...initialState?.settings,
  };
};


// let authRoute :API.MenuItem[]= [];

// function parseRoutes(as: API.MenuItem[]) {
//   if (as) {
//     return as.map(item => {
//       let n = {
//         path : item.key,
//         name: item.title,
//         routes :[
//           {}
//         ],
//       }

//       if (item.children) {
//         // item.component.leng
//         item.children.map(e => ({
//           path : e.key,
//           name: e.title,
//           // component: e.component,
//           component: require(e.component || ''),
//         })).forEach(e1 => {
//           if (e1.path != "") {
//             n.routes.push(e1);
//           }
//         });
//       }

//       return n;
//     })
//   }

//   return [];
// }


// export function patchRoutes({routes}) {
//   console.log("------patchRoutes-----", authRoute);
//   console.log("------patchRoutes2-----", routes);
//   parseRoutes(authRoute).forEach(item=>routes[1].routes.push(item))
//   console.log("------patchRoutes3-----", routes);
// }

// export function render(oldRoutes: Function) {
//   const call =  getMenu();
//   call.then(res => {
//     authRoute = res.data || [];
//     console.log("------render-----", authRoute);
//     oldRoutes();

//   })
// }

ConfigProvider.config({
  theme: {
    primaryColor: '#FA541C',
  },
});