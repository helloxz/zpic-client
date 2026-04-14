# Zpic Client

Zpic Client是一款开源、实用的跨平台图床客户端，适用于 [ImgURL](https://www.imgurl.org/)、 [图链](https://go.piclink.cc/) 及 [Zpic Pro](https://www.zpic.pro/)，使用Golang + Wails开发。

## 特点

* **跨平台：** 支持 macOS、Windows、Linux 三大主流平台。
* **批量上传：** 支持扫描目录上传，可一次性上传数千张图片。
* **URL 上传：** 支持通过图片 URL 直接上传，方便快捷。
* **粘贴上传：** 支持Ctrl + V粘贴上传，粘贴图片后自动上传并复制链接到剪贴板。
* **表格导出：** 可将上传结果导出为 CSV 表格，方便导入其它系统。
* **低占用：** 得益于 Golang + Wails 开发，资源占用极低，运行流畅高效。
* **图片压缩：** 扫描上传时自动对大于 300KB 的 PNG/JPG 图片进行无损压缩。
* **图片去重：** 智能检测重复图片，自动跳过已上传的图片，避免重复占用空间。

## 部分截图

![CleanShot 2026-02-10 at 09.32.32@2x.png](https://v.uuu.ovh/2026/02/10/7qHFNtIb.png)

![CleanShot 2026-02-10 at 09.33.53@2x.png](https://v.uuu.ovh/2026/02/10/0tjn8qvJ.png)

![CleanShot 2026-02-10 at 09.33.16@2x.png](https://v.uuu.ovh/2026/02/10/Srl3F3jo.png)

![CleanShot 2026-02-10 at 09.34.33@2x.png](https://v.uuu.ovh/2026/02/10/cyaqlkSC.png)

## 开发环境

### 后端依赖

- **Go**: 1.24.0 或更高版本
- **Wails**: v2.11.0


### 前端依赖

- **Node.js**: 18.0 或更高版本
- **pnpm**: 8.0 或更高版本
- **Vue**: 3.2+
- **Vite**: 3.0+


## 编译指令

### 开发模式

```bash
# 启动完整应用开发模式（前后端热重载）
wails dev
```

### 生产构建

```bash
# 构建完整应用
wails build
```

## 联系我们

* 博客：[https://blog.xiaoz.org/](https://blog.xiaoz.org/)
* ImgURL：[https://www.imgurl.org/](https://www.imgurl.org/)
* 图链：[https://go.piclink.cc/](https://go.piclink.cc/)

