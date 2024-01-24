import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
import './index.css'
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Home from './Home/index.jsx';
import About from './About/index.jsx';
import News from './News/index.jsx';
import New from './New/index.jsx';
import Apps from './Apps/index.jsx';
import Privacy from './Privacy/index.jsx';
import Userfile from './Userfile/index.jsx';
import OneApp from './App/index.jsx';

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        path: "/",
        element: <Home />,
      },
      {
        path: "/home",
        element: <Home />,
      },
      {
        path: "/news",
        element: <News />,
      },
      {
        path: "/apps",
        element: <Apps />,
      },
      {
        path: "/about",
        element: <About />,
      },
      {
        path: "/new",
        element: <New />,
      },
      // {
      //   path: "/privacy",
      //   element: <Privacy />,
      // },
      // {
      //   path: "/userfile",
      //   element: <Userfile />,
      // },
      {
        path: "/app",
        element: <OneApp />,
      }
    ],
  },
]);

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
)
