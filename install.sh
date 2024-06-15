#!/bin/bash

# 获取用户的 HOME 目录和 GOPATH
USER_HOME=$(eval echo ~)
GOPATH=$(go env GOPATH)

# 目标目录
TEMPLATES_TARGET_DIR="$USER_HOME/.cpj/templates"
BIN_TARGET_DIR="$GOPATH/bin"

# 确保目标目录存在，不存在则创建
mkdir -p "$TEMPLATES_TARGET_DIR"
mkdir -p "$BIN_TARGET_DIR"

# Clone git repository and copy templates directory
git clone git@github.com:0226zy/cpj.git /tmp/cpj_repo    # 克隆仓库到临时目录
rsync -av /tmp/cpj_repo/templates/ "$TEMPLATES_TARGET_DIR"  # 将templates目录同步到目标目录

echo "Templates directory copied to $TEMPLATES_TARGET_DIR"

# 编译项目并安装可执行文件到 GOPATH/bin
cd /tmp/cpj_repo
go install

echo "cpj installed to $BIN_TARGET_DIR"

# 清理临时目录
rm -rf /tmp/cpj_repo
