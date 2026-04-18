# Agent Spec — Pix2Vec (Wails + Vue3 + Go)

## 🧠 项目目标

构建一个桌面应用（Windows 优先），用于：

* PNG → SVG（矢量化）
* SVG → PNG（渲染）
* 支持拖拽、预览、批量处理

技术栈：

* 后端：Go (Wails v2)
* 前端：Vue3 + Vite
* 执行引擎：CLI（vtracer / resvg / potrace）

---

## ⚙️ 架构设计

采用 **CLI 驱动 + 前后端分离（Wails桥接）**

```
Vue3 UI
   ↓ (Wails Bind)
Go Service Layer
   ↓
CLI Tools (vtracer / resvg / potrace)
```

---

## 📦 目录结构

```
pix2vec/
├── main.go
├── app.go                # Wails 绑定入口
├── internal/
│   ├── converter/
│   │   ├── convert.go    # 核心转换逻辑
│   │   ├── runner.go     # CLI 调用封装
│   │   └── types.go
│   ├── utils/
│   │   ├── path.go
│   │   └── exec.go
│
├── frontend/
│   ├── src/
│   │   ├── App.vue
│   │   ├── components/
│   │   ├── views/
│   │   │   └── Home.vue
│   │   ├── services/
│   │   │   └── api.ts   # 调用 Go
│   │   └── store/
│
├── bin/
│   ├── vtracer.exe
│   ├── resvg.exe
│   └── potrace.exe
│
└── wails.json
```

---

## 🔧 核心功能

### 1️⃣ PNG → SVG（彩色）

使用：

```
vtracer input.png -o output.svg
```

支持参数：

* color_mode（默认）
* max_colors（可选）
* filter_speckle（降噪）

---

### 2️⃣ PNG → SVG（黑白）

流程：

```
png → bmp → potrace → svg
```

命令：

```
potrace input.bmp -s -o output.svg
```

---

### 3️⃣ SVG → PNG

使用：

```
resvg input.svg output.png
```

支持：

* width / height
* background（可选）

---

## 🧩 Go API 设计（必须实现）

```go
type ConvertRequest struct {
    InputPath  string
    OutputPath string
    Mode       string // "color" | "bw" | "render"
}

type ConvertResponse struct {
    Success bool
    Error   string
}
```

---

### 核心方法

```go
func (c *Converter) Convert(req ConvertRequest) (*ConvertResponse, error)
```

---

## ⚙️ CLI 调用规范

必须使用：

```go
exec.Command(...)
```

要求：

* 禁止使用 shell
* 必须捕获 CombinedOutput
* 必须返回详细错误

---

## 📁 路径处理（强制）

必须实现：

```go
func GetExeDir() string
```

所有 CLI 路径必须：

```go
filepath.Join(GetExeDir(), "bin", "vtracer.exe")
```

禁止使用相对路径

---

## 🖥️ 前端要求（Vue3）

### UI 功能

* 拖拽上传
* 文件选择
* 转换按钮
* 实时状态（loading / success / error）
* 图片预览（前后对比）

---

### API 调用（Wails）

```ts
import { Convert } from "../../wailsjs/go/app/App"

await Convert({
  inputPath,
  outputPath,
  mode
})
```

---

## 🎯 模式设计（前端）

```ts
type Mode = "color" | "bw" | "render"
```

| 模式     | 功能               |
| ------ | ---------------- |
| color  | 彩色矢量（vtracer）    |
| bw     | 黑白描边（potrace）    |
| render | SVG → PNG（resvg） |

---

## 🚀 用户流程

1. 拖入图片
2. 选择模式
3. 点击转换
4. 显示结果
5. 支持打开输出目录

---

## ⚠️ 约束（必须遵守）

* 不使用 CGO
* 不使用 WASM
* 不依赖系统 PATH
* 所有 CLI 必须随程序打包
* Windows 优先（路径兼容）

---

## 🧪 错误处理

必须返回：

* CLI stderr
* 文件不存在
* 格式错误

---

## 📦 打包要求

输出：

```
/dist
  Pix2Vec.exe
  /bin
    *.exe
```

---

## ✨ 可选增强（如果模型有能力实现）

* 批量处理
* 参数调节 UI
* 拖拽文件夹
* 转换队列

---

## 🎯 代码风格

* Go：清晰结构，不写一堆 main.go
* Vue：Composition API
* 禁止过度抽象
* 优先可读性

---

## 界面设计风格

全中文，遵循 DESIGN.md

---

## 📌 最终目标

生成一个：

👉 **可运行的桌面应用（双击即用）**

支持：

* PNG ↔ SVG
* 稳定调用 CLI
* 无需额外安装依赖

外加 生成一个专业好看的README.md (需要产品图片的地方可以占位)
---

END
