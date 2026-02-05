# AGENTS.md

本文档包含适用于此 Wails + Vue.js 图床客户端项目的智能编码助手的说明。

## 项目概述

这是一个结合了 Go 后端和 Vue.js 前端的 Wails 图床客户端。前端位于 `frontend/` 目录中，使用 Vue 3 和 Composition API (`<script setup>` 语法)。

## 项目要求

* 前端UI使用Naive UI，图标使用：@vicons/ionicons5 作为图标库
* 后端网络请求使用：resty/v2
* 如果前端涉及到路由，请使用Vue Router
* 如果需要状态管理，请使用Pinia

## 构建命令

### 开发

- `wails dev` - 以开发模式运行应用程序，支持热重载
- `cd frontend && pnpm run dev` - 仅运行前端开发服务器

### 构建

- `wails build` - 构建应用程序的生产版本
- `cd frontend && pnpm run build` - 仅构建前端资源

### 前端包管理

- `cd frontend && pnpm install` - 安装前端依赖

## 测试

目前该项目未配置测试框架。如需添加测试：

1. 对于 Vue 组件，考虑使用 Vue Test Utils 配合 Vitest
2. 对于 Go 代码，使用标准的 `go test` 命令
3. 运行单个 Go 测试：`go test -run TestName ./path/to/package`

## 代码检查和格式化

该项目目前未配置代码检查或格式化工具。如需添加：

1. 对于 Vue/JavaScript/TypeScript，考虑使用 ESLint 配合 Prettier
2. 对于 Go，使用 `golint` 和 `go fmt`

## 代码风格指南

### Vue 文件

1. 使用 `<script setup>` 语法（Composition API）而非 Options API
2. 将导入语句放在 `<script setup>` 部分的顶部
3. 使用 Vue 的 `reactive()` 或 `ref()` 处理响应式数据
4. 遵循单文件组件结构：`<script setup>`、`<template>`、`<style>`
5. 尽可能使用 scoped 样式以避免样式泄露

### JavaScript/TypeScript

1. 使用 ES6+ 特性（箭头函数、解构等）
2. 当变量值不改变时，优先使用 const 而非 let
3. 使用模板字面量进行字符串插值
4. 使用描述性的变量和函数名
5. 保持函数简短并专注于单一职责

### Go

1. 遵循标准 Go 约定（https://golang.org/doc/effective_go）
2. 导出函数/变量使用驼峰命名法
3. 为导出函数添加注释说明其用途
4. 显式处理错误 - 不要忽略它们
5. 适当使用 context.Context 处理请求范围的数据

### 导入

1. 在 Vue 文件中，将依赖导入放在 `<script setup>` 部分的顶部
2. 在 Go 文件中，将标准库导入与第三方库导入分组
3. 在每组内按字母顺序排列导入

### 命名约定

1. Vue 组件名和文件名使用 PascalCase（大驼峰命名法）
2. JavaScript/TypeScript 变量和函数使用 camelCase（小驼峰命名法）
3. Go 变量和函数使用 camelCase（小驼峰命名法）
4. Go 常量使用 UPPER_SNAKE_CASE（大写蛇形命名法）
5. CSS 类名和文件名使用 kebab-case（短横线隔开命名法）

### 错误处理

1. 在 Go 中，始终适当检查和处理错误
2. 返回错误而不是记录日志（让调用者决定如何处理）
3. 在 Vue/JavaScript 中，对异步操作使用 try/catch
4. 在 UI 中显示用户友好的错误消息

### 样式

1. 在 Vue 组件中使用 scoped 样式防止样式泄露
2. 使用 CSS 类而不是内联样式
3. 遵循一致的 CSS 类命名约定（例如 BEM）
4. 适当使用相对单位（rem, em）而非绝对单位（px）

## 项目结构

- `/` - 根目录，包含 Go 后端文件和 Wails 配置
- `/frontend` - Vue.js 前端应用程序
- `/frontend/src` - Vue 源文件
- `/frontend/src/components` - Vue 组件
- `/frontend/dist` - 构建的前端资源（生成的）

## Wails 特性

1. 绑定到前端的 Go 方法在 `frontend/wailsjs/go/` 中自动生成
2. 前端可以直接从生成的绑定中导入 Go 方法
3. 在 `main.go` 中使用 `wails.Run()` 配置应用程序
4. `wails.json` 文件包含项目配置，包括构建命令

## 添加新功能

1. 对于新 UI 组件，在 `frontend/src/components/` 中创建
2. 对于新 Go 功能，向 `app.go` 中的 App 结构体添加方法
3. 记得在 `main.go` 的 `Bind` 部分绑定新的 Go 方法
4. 添加新依赖时，Go 使用 `go get`，前端使用 `pnpm add`

## 常见任务

1. 添加新页面/视图时，创建新的 Vue 组件并添加路由
2. 添加新后端功能时，向 App 结构体添加方法
3. 前后端通信时，调用生成的 Go 绑定
4. 添加外部库时，Go 使用 `go get`，前端使用 `pnpm add`

## 特定技术要求

### 前端 UI 组件

1. 使用 Naive UI 组件库构建界面
2. 使用 xicons 作为图标库
3. 确保所有 UI 组件遵循 Naive UI 的设计规范

### 后端网络请求

1. 使用 resty/v2 库处理所有 HTTP 请求
2. 在 Go 代码中统一使用该库进行网络通信
3. 正确处理请求和响应的错误情况

本指南应能帮助智能编码助手在此项目中有效贡献代码，同时保持与现有代码的一致性。