#!/bin/bash
# 如果没有tutorial文件夹,则创建
if [ ! -d "tutorial" ]; then
    mkdir tutorial
fi
protoc -I . person.proto --go_out=plugins=grpc:tutorial
# 编译客户端
go build personclient.go
# 编译服务端
go build personserver.go
