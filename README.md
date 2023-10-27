# hjing

微服务框架

# 计划列表

- [ ] hjing cli

- [ ] 支持配置和ETCD的微服务注册

- [ ] RPC

## 安装依赖

```bash
# 安装protoc, 解压后复制到PATH中
wget https://github.com/protocolbuffers/protobuf/releases/download/v24.4/protoc-24.4-osx-x86_64.zip

# 安装protoc-gen-gogu
go install gitee.com/jkkkls/protobuf/cmd/protoc-gen-gogu@latest

# 安装gonew
go install golang.org/x/tools/cmd/gonew@latest

# 安装项目工具hjing
go install github.com/jkkkls/hjing/cmds/hjing@latest
```

## 创建一个项目

```bash

# 创建空项目
gonew github.com/jkkkls/hjing/layout/project github.com/<authorName>/<projectName>

cd <projectName>
# 添加应用
hjing add-app <AppName>

# 添加模版服务到指定应用中
hjing add-svc <AppName> <ServiceName>

# 添加模版接口到指定服务中
hjing add-itf <ServiceName> <interfaceName>

# 添加网关模版到应用中
hjing add-kcp-gate <AppName>
hjing add-ws-gate <AppName>
hjing add-tcp-gate <AppName>

# 添加模版数据库模块到应用中，postgres+grom
hjing add-db <AppName>

# 编译项目
make <projectName>

# 启动项目
cd build
./<projectName>

```