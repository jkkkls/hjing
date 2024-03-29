# hjing

***了解go即可使用的全栈微服务框架***

框架特点：

* 代码自动生成
* 集成可用的管理后台
* 集成http api网关
* 集成长链接网关

# 计划列表

- [ ] 集成后台，支持监控，日志，配置等
- [ ] http api网关，支持监控，统计，熔断降载等
- [ ] 优化长链接网关，支持游戏等业务快速开发

## 安装依赖

```bash
# 安装protoc, 解压后复制到PATH中，在https://github.com/protocolbuffers/protobuf/releases找到最新版本即可
wget https://github.com/protocolbuffers/protobuf/releases/download/v24.4/protoc-24.4-osx-x86_64.zip

# 安装protoc-gen-gogu
# https://gitee.com/jkkkls/protobuf/releases/tag/v1.32.0下载自己需要的版本

# goimports 用于格式化生成代码
go install golang.org/x/tools/cmd/goimports@latest

# 安装项目工具hjing
go install github.com/jkkkls/hjing/cmds/hjing@latest

# 安装配置表工具xlsx2proto，用于配置生成代码，不需要可不安装
go install github.com/jkkkls/hjing/cmds/xlsx2proto@latest

# 安装nodejs yarn, 用于后台的开发和打包
# brew install node yarn
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

- Landing [落地页前端](https://github.com/jkkkls/Landing.git)

- Axjmin [后台前后端](https://github.com/jkkkls/axjmin.git), 已经集成到admin模块中

## 示例

1. 创建项目

``` shell
# 创建项目
➜  hjing_space -> hjing new github.com/jkkkls/test_app
create project[github.com/jkkkls/test_app] success

➜  hjing_space -> cd test_app
# 添加网关应用, 之后记住修改gate.yaml的节点端口
$hjing add-app gate
create app[gate] success
# 添加数据应用, 之后记住修改data.yaml的节点端口
$hjing add-app data
create app[data] success

#修改data.yaml的节点端口和id
app:
  id: 2
  name: data
  desc: data
  host: 127.0.0.1
  port: 10002
  type: data
#修改gate.yaml的节点端口和id，增加内网http端口
app:
  id: 3
  name: gate
  desc: gate
  host: 127.0.0.1
  port: 10003
  type: gate
  httpPort: 8081

# 编译admin,data和gate
$go mod tidy
$make pb
$make admin
$make data
fatal: not a git repository (or any of the parent directories): .git
编译data完成
$make gate
fatal: not a git repository (or any of the parent directories): .git
编译gate完成

#1. 启动etcd
cd etcd
./etcd

#2. 启动data应用
$cd build
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
$cd build
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
$hjing add-svc data db
create service[db] success

# 添加两个接口到db服务中
$hjing  add-itf db get
create interface[get] for db success
$hjing  add-itf db set
create interface[set] for db success
```

3. 完善协议文件[pb/db.proto]

```
syntax = "proto3";
package pb;  // 声明所在包
option go_package = "github.com/jkkkls/test_app/pb";

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
$make pb
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
# 通过gate的rpcapi接口调用data的db服务
$curl -i -X POST -d '{"key":"a", "value":"b"}' 'http://127.0.0.1:8081/rpcapi/db/set'
HTTP/1.1 200 OK
Server: fasthttp
Date: Tue, 06 Feb 2024 09:48:56 GMT
Content-Type: application/json
Content-Length: 2

{}

$curl -i -X POST -d '{"key":"a"}' 'http://127.0.0.1:8081/rpcapi/db/get'
HTTP/1.1 200 OK
Server: fasthttp
Date: Tue, 06 Feb 2024 09:48:32 GMT
Content-Type: application/json
Content-Length: 13

{"value":"b"}
```


### admin 管理后台 默认账号密码 admin123/admin123

![avatar](png/login.png)
![avatar](png/welcome.png)
![avatar](png/monitor.png)