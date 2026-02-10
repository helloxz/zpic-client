#!/bin/bash

ARG1=$1

# 根据ARG参数进行不同的编译
if [ "$ARG1" == "macos" ]; then
  wails build \
  -clean \
  -ldflags "-s -w" \
  -trimpath
  exit 0
elif [ "$ARG1" == "linux" ]; then
  wails build -tags webkit2_41 \
  -clean \
  -ldflags "-s -w" \
  -trimpath
  exit 0
else
  echo "无效的参数: $ARG1"
  exit 1
fi



# 1. 创建目标目录（注意引号）
# mkdir -p "./build/bin/Zpic Client.app/Contents/MacOS/bin/darwin"

# 2. 复制 darwin 目录下的所有内容（不包含目录本身）
# cp -r ./bin/darwin/. "./build/bin/Zpic Client.app/Contents/MacOS/bin/darwin/"