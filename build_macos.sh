wails build \
  -clean \
  -ldflags "-s -w" \
  -trimpath

# 1. 创建目标目录（注意引号）
# mkdir -p "./build/bin/Zpic Client.app/Contents/MacOS/bin/darwin"

# 2. 复制 darwin 目录下的所有内容（不包含目录本身）
# cp -r ./bin/darwin/. "./build/bin/Zpic Client.app/Contents/MacOS/bin/darwin/"