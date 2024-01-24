# Ant Design Pro

This project is initialized with [Ant Design Pro](https://pro.ant.design). Follow is the quick guide for how to use.

## Environment Prepare

Install `node_modules`:

```bash
npm install
```

or

```bash
yarn
```

## Provided Scripts

Ant Design Pro provides some useful script to help you quick start and build with web project, code style check and test.

Scripts provided in `package.json`. It's safe to modify or add additional script:

### Start project

```bash
npm start
```

### Build project

```bash
npm run build
```

### Check code style

```bash
npm run lint
```

You can also use script to auto fix some lint error:

```bash
npm run lint:fix
```

### Test code

```bash
npm test
```

## More



💥 feat(模块): 添加了个很棒的功能
🐛 fix(模块): 修复了一些 bug
📝 docs(模块): 更新了一下文档
🌷 UI(模块): 修改了一下样式
🏰 chore(模块): 对脚手架做了些更改
🌐 locale(模块): 为国际化做了微小的贡献


export NODE_OPTIONS=--openssl-legacy-provider
yarn upgrade-interactive --latest



1. 阿里云ALB网关
2. 自研evproc服务
  支持多节点运行，负载均衡
  备份功能，可以通过内网的http接口下载指定游戏和日期的备份文件
  写入数据库失败处理: 数据库异常的话
    - 写入出错备份，正常时从出错备份恢复恢复数据
  插件服务注册功能，由nsq实现
3. 插件服务模块，支持一些埋点实时周边功能的开发
4. 埋点系统后台
  应用管理，支持测试服, 支持game，ad埋点选择
  运行节点状态，cpu，内存。ALB网关转发状态，可以通过后台控制节点的负载均衡
  服务器状态
  报警处理
    支持邮件，短信
    数据库写入异常报警
    报错信息
    服务器cpu和内存状态，过高报警
    进程cpu，内存，协程状态，暴增状态，过高报警
5. 埋点类型
  -