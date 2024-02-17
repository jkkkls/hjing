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

- web/front [落地页前端](https://github.com/jkkkls/Landing.git)

- web/axjmin [后台前后端](https://github.com/jkkkls/axjmin.git)

## 示例

1. 创建项目

``` shell
# 创建项目
➜  hjing_space -> hjing new github.com/jkkkls/test_app
create project[github.com/jkkkls/test_app] success

➜  hjing_space -> cd test_app
# 添加网关应用, 之后记住修改gate.yaml的节点端口
➜  test_app -> hjing add-app gate
create app[gate] success
# 添加数据应用, 之后记住修改data.yaml的节点端口
➜  test_app -> hjing add-app data
create app[data] success

#修改data.yaml的节点端口和id
app:
  id: 2
  name: data
  desc: data
  host: 127.0.0.1
  port: 10002
  type: data
#修改gate.yaml，增加内网http端口


# 编译data和gate
➜  test_app -> make data
fatal: not a git repository (or any of the parent directories): .git
编译data完成
➜  test_app -> make gate
fatal: not a git repository (or any of the parent directories): .git
编译gate完成

#1. 启动etcd
cd etcd
./etcd

#2. 启动data应用
➜  test_app -> cd build
➜  build -> ./data
Version      : 2000.01.01.release
Git SHA      : xxx
Build PC     : anyone
Build Time   : 2000-01-01T00:00:00+0800
Build Area   :
2024-02-06 13:07:33.558153 INFO rpc_node.go:207 初始化服务节点 [id=1 name=data host=127.0.0.1 port=10001 region=0]
2024/02/06 13:07:33 RegisterNode data 127.0.0.1:10001 []
2024/02/06 13:07:33 Get All global Node data 127.0.0.1:10001
2024-02-06 13:08:48.054986 INFO rpc_node.go:133 连接节点成功 [name=gate address=127.0.0.1:10001 region=0 isClose=false]

#3. 启动gate应用
➜  test_app -> cd build
➜  build -> ./gate
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

```

2. 添加服务

``` shell
# 添加数据库服务到data应用中
➜  test_app -> hjing add-svc data db
create service[db] success

# 添加两个接口到db服务中
➜  test_app -> hjing  add-itf db get --open
create interface[get] for db success
➜  test_app -> hjing  add-itf db set --open
create interface[set] for db success
```

3. 完善协议文件[pb/db.proto]

```
syntax = "proto3";
package pb;  // 声明所在包
option go_package = "github.com/jkkkls/test_app/pb";
//import "libs/pb/cli_common.proto";

message GetReq {
	string key = 1;
}
message GetRsp {
	string value = 1;
}

message SetReq {
	string key = 1;
	string value = 2;
}
message SetRsp {
	bool ok = 1;
}

```

4. 重新编译协议文件
``` shell
➜  test_app -> make pb
```

5. 添加接口实现[services/db/service.go]

``` go
package db

import (
	"github.com/jkkkls/hjing/rpc"
	"github.com/jkkkls/test_app/pb"
)

// DbService 服务
type DbService struct {
	kvs map[string]string
}

func (service *DbService) NodeConn(name string)                  {}
func (service *DbService) NodeClose(name string)                 {}
func (service *DbService) OnEvent(eventName string, args ...any) {}

// Exit 退出处理
func (service *DbService) Exit() {}

// Run 服务启动函数
func (service *DbService) Run() error {
	service.kvs = make(map[string]string)
	return nil
}

func (service *DbService) Get(context *rpc.Context, req *pb.GetReq, rsp *pb.GetRsp) (ret uint16, err error) {
	rsp.Value = service.kvs[req.Key]
	return
}
func (service *DbService) Set(context *rpc.Context, req *pb.SetReq, rsp *pb.SetRsp) (ret uint16, err error) {
	service.kvs[req.Key] = req.Value
	return
}

```

6. 重新编译gate和data。依次启动etcd、gate和data. 通过http请求验证接口是否正常。
``` shell
➜  test_app -> curl -i -X POST -d '{"key":"a", "value":"b"}' 'http://127.0.0.1:8081/rpcapi/db/set'
HTTP/1.1 200 OK
Server: fasthttp
Date: Tue, 06 Feb 2024 09:48:56 GMT
Content-Type: application/json
Content-Length: 2

{}

➜  test_app -> curl -i -X POST -d '{"key":"a"}' 'http://127.0.0.1:8081/rpcapi/db/get'
HTTP/1.1 200 OK
Server: fasthttp
Date: Tue, 06 Feb 2024 09:48:32 GMT
Content-Type: application/json
Content-Length: 13

{"value":"b"}
```