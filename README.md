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
make build

# 启动项目
cd build
./<AppName>

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

## 示例

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

#修改data.yaml的节点端口和id
app:
  id: 2
  name: data
  desc: data
  host: 127.0.0.1
  port: 10002
  type: data

# 编译data和gate
➜  test_app make data
fatal: not a git repository (or any of the parent directories): .git
编译data完成
➜  test_app make gate
fatal: not a git repository (or any of the parent directories): .git
编译gate完成

#1. 启动etcd
cd etcd
./etcd
#2. 启动data应用
Version      : 2000.01.01.release
Git SHA      : xxx
Build PC     : anyone
Build Time   : 2000-01-01T00:00:00+0800
Build Area   :
2024-02-06 13:07:33.558153 INFO rpc_node.go:207 初始化服务节点 [id=1 name=data host=127.0.0.1 port=10001 region=0]
2024/02/06 13:07:33 RegisterNode data 127.0.0.1:10001 []
2024/02/06 13:07:33 Get All global Node data 127.0.0.1:10001
#3. 启动gate应用
Version      : 2000.01.01.release
Git SHA      : xxx
Build PC     : anyone
Build Time   : 2000-01-01T00:00:00+0800
Build Area   :
2024-02-06 13:08:48.031718 INFO rpc_node.go:207 初始化服务节点 [id=1 name=gate host=127.0.0.1 port=10001 region=0]
2024/02/06 13:08:48 RegisterNode gate 127.0.0.1:10001 []
2024/02/06 13:08:48 Get All global Node data 127.0.0.1:10002
2024/02/06 13:08:48 Get All global Node gate 127.0.0.1:10001
2024-02-06 13:08:48.054966 INFO rpc_node.go:133 连接节点成功 [name=data address=127.0.0.1:10002 region=0 isClose=false]
4. 启动后data的日志
2024-02-06 13:08:48.054986 INFO rpc_node.go:133 连接节点成功 [name=gate address=127.0.0.1:10001 region=0 isClose=false]

```