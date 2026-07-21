---
name: build-release
description: 编译打包 ZpicClient，根据操作系统执行不同流程（macOS→darwin_arm64, Windows→windows_amd64）
---

# 编译打包 ZpicClient

根据当前操作系统执行编译和打包操作。

## 前置准备

1. 确保在项目根目录下执行
2. 确保已安装 `wails` CLI 工具
3. 确保已安装 `zip` 命令（macOS 自带，Windows 需确认）

## 步骤

### 1. 获取版本号

读取 `core/app.go` 中的 `VERSION` 变量：

```bash
grep 'var VERSION' core/app.go | head -1
```

提取版本号，例如 `1.3.2`。

### 2. 检测操作系统

使用 Go 的 `runtime.GOOS` 概念判断当前系统，或者通过 `uname -s` 检测：

- macOS → `darwin`
- Windows → `windows`

### 3. macOS 编译打包

```bash
# 编译
bash build.sh macos

# 打包为 zip（arm64）
cd build/bin
zip -r ../../ZpicClient_v{version}_darwin_arm64.zip ZpicClient.app
cd ../..
```

### 4. Windows 编译打包

```powershell
# 编译
wails build -clean -ldflags "-s -w" -trimpath

# 打包为 zip（amd64）
Compress-Archive -Path build\bin\ZpicClient.exe -DestinationPath ZpicClient_v{version}_windows_amd64.zip
```

或在 Git Bash 中使用 zip：

```bash
# 编译
wails build -clean -ldflags "-s -w" -trimpath

# 打包为 zip
cd build/bin
zip ../../ZpicClient_v{version}_windows_amd64.zip ZpicClient.exe
cd ../..
```

## 输出

- macOS: `ZpicClient_v{version}_darwin_arm64.zip`（项目根目录）
- Windows: `ZpicClient_v{version}_windows_amd64.zip`（项目根目录）

## 验证

检查生成的 zip 文件：

```bash
ls -lh ZpicClient_v*.zip
```
