#!/bin/bash

# 获取用户的 HOME 目录
USER_HOME=$(eval echo ~)

# 目标目录
TARGET_DIR="$HOME/.cpj"

# 源目录
SOURCE_DIR="./templates"

# 确保目标目录存在，不存在则创建
mkdir -p "$TARGET_DIR"

# 复制源目录到目标目录
cp -r "$SOURCE_DIR" "$TARGET_DIR"

echo "cpj templates directory copied to $TARGET_DIR"

# 编译项目并安装到 GOPATH/bin
go install

echo "Executable installed to GOPATH/bin"
