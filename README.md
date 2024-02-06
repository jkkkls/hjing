# hjing

微服务框架

# 计划列表

- [ ] hjing cli

- [ ] 支持配置和ETCD的微服务注册

- [ ] RPC

- []

## 安装依赖

```bash
# 安装protoc, 解压后复制到PATH中
wget https://github.com/protocolbuffers/protobuf/releases/download/v24.4/protoc-24.4-osx-x86_64.zip

# 安装protoc-gen-gogu
go install gitee.com/jkkkls/protobuf/cmd/protoc-gen-gogu@latest

# goimports
go install golang.org/x/tools/cmd/goimports@latest

# 安装项目工具hjing
go install github.com/jkkkls/hjing/cmds/hjing@latest
```

## 创建一个项目

```bash

# 创建空项目
hjing new github.com/<authorName>/<projectName>

cd <projectName>
# 添加应用
hjing add-app <AppName>

# 添加模版服务到指定应用中
hjing add-svc <AppName> <ServiceName>

# 添加模版接口到指定服务中
hjing add-itf <ServiceName> <interfaceName> --open

# 添加HTTP网关模版到应用中
hjing add-gate <AppName>


# 添加网关模版到应用中
hjing add-kcp-gate <AppName>
hjing add-ws-gate <AppName>
hjing add-tcp-gate <AppName>

# 添加模版数据库模块到应用中，postgres+grom
hjing add-db <DbName>

# 编译项目
make <projectName>

# 启动项目
cd build
./<projectName>

```

## 项目目录结构

- web/front 落地页前端
```
yarn install
yarn dev
```


- web/axjmin 后台前后端
```
# 启动后端
# 先配置server.yaml的mysql
# 默认账号密码admin123/admin123
cd web/axjmin/backend
go build
./webserver

# 启动前端
cd web/axjmin/frontend
yarn install
yarn start:no-mock
```

## 实例

``` shell
# 创建项目
➜  hjing_space hjing/cmds/hjing/hjing new github.com/jkkkls/test_app
create project[github.com/jkkkls/test_app] success

➜  hjing_space cd test_app
# 添加网关应用, 之后记住修改gate.yaml的节点端口
➜  test_app ../hjing/cmds/hjing/hjing add-app gate
create app[gate] success
# 添加数据应用, 之后记住修改data.yaml的节点端口
➜  test_app ../hjing/cmds/hjing/hjing add-app data
create app[data] success




```