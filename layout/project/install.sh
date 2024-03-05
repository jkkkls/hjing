#!/bin/bash

#
# 环境部署更新脚本
#

source env_proc

os="unknown"
arch="unknown"

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    os="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    os="darwin"
else
    exit 0
fi

if [[ "$(uname -m)" == "arm64" ]]; then
    arch="arm64"
elif [[ "$(uname -m)" == "x86_64" ]]; then
    arch="amd64"
else
    exit 0
fi


echo "os: $os arch: $arch"
echo "path: $dir"

supervisord=$(which supervisord)
if [[ ${supervisord} == "" ]]; then
    echo "supervisord not found"
    download_url="https://gitee.com/jkkkls/protobuf/releases/download/v1.32.0/supervisord.$os.$arch.zip"
    echo "download supervisord from ${download_url}"
    curl -L ${download_url} > /tmp/supervisord.zip
    sudo unzip /tmp/supervisord.zip -d /usr/local/bin
fi
