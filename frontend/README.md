# AIVideoNote 前端 — Vue 3 SPA

## 技术栈

| 组件 | 库 | 用途 |
|------|-----|------|
| 框架 | Vue 3 | 响应式 SPA |
| 构建 | Vite 5 | 开发/构建 |
| 语言 | TypeScript | 类型安全 |
| 路由 | vue-router 4 | SPA 路由 |
| 状态 | Pinia 2 | 全局状态管理 |
| HTTP | axios | API 调用 |
| 样式 | 原生 CSS (scoped) | 组件级样式隔离 |
| Markdown | marked + highlight.js | 笔记渲染 |

## 目录结构

```
frontend/
├── index.html
├── package.json
├── vite.config.ts
├── tsconfig.json
├── src/
│   ├── main.ts                    # 应用入口
│   ├── App.vue                    # 根组件
│   ├── router/index.ts            # 路由定义
│   ├── api/                       # 后端 API 封装
│   │   ├── request.ts             # Axios 实例 + 拦截器
│   │   ├── task.ts                # 任务相关 API
│   │   ├── group.ts               # 分组相关 API
│   │   ├── provider.ts            # 提供商 API
│   │   └── upload.ts              # 上传 API
│   ├── stores/                    # Pinia 状态管理
│   │   ├── task.ts                # 任务状态 + 3s 轮询
│   │   ├── group.ts               # 分组列表
│   │   └── provider.ts            # 提供商列表
│   ├── views/                     # 页面
│   │   ├── HomeView.vue           # 首页（三栏布局）
│   │   ├── SettingsView.vue       # 设置页
│   │   └── NotFoundView.vue       # 404
│   ├── components/                # 组件
│   │   ├── layout/AppLayout.vue   # 全局布局
│   │   ├── note/
│   │   │   ├── NoteForm.vue       # 表单（链接 + 模型 + 分组）
│   │   │   ├── MarkdownViewer.vue # Markdown 渲染
│   │   │   └── StepBar.vue        # 任务进度条
│   │   ├── task/TaskHistory.vue   # 任务历史（分组树形）
│   │   └── provider/              # 提供商管理组件
│   └── types/index.ts             # 全局类型声明
└── dist/                          # 构建输出
```

## 路由

| 路径 | 组件 | 说明 |
|------|------|------|
| `/` | HomeView.vue | 主页：左侧表单 + 中间笔记 + 右侧任务历史 |
| `/settings` | SettingsView.vue | 设置：LLM 提供商管理 |
| `/:pathMatch(.*)*` | NotFoundView.vue | 404 |

## 三栏布局

```
┌──────────────────────────────────────────────────────────┐
│  AppLayout (顶部导航)                                     │
├────────────┬────────────────────────┬────────────────────┤
│  Sidebar   │       Content          │   History Panel    │
│  (320px)   │       (flex:1)         │   (260px)          │
│            │                        │                    │
│  NoteForm  │   MarkdownViewer       │   分组树            │
│  StepBar   │   [导出笔记]           │   └ 任务列表        │
│            │                        │                    │
└────────────┴────────────────────────┴────────────────────┘
```

## 状态管理

### taskStore
- `submitTask()` — 提交笔记生成任务
- `loadTask()` — 查看已有任务结果
- `currentStatus` — 当前任务流水线状态
- `currentMarkdown` — AI 生成的 Markdown 笔记
- 自动 3 秒轮询 `GET /task/:id/status`，完成/失败时停止

### groupStore
- `fetchGroups()` — 加载分组列表
- `addGroup()` — 创建分组
- `updateGroupName()` / `removeGroup()` — 管理分组

### providerStore
- `fetchProviders()` — 加载 LLM 提供商列表
- CRUD 操作提供商配置

## 启动

```bash
cd frontend
npm install

# 开发模式
npm run dev          # http://localhost:3015

# 生产构建
npm run build        # 输出到 dist/
npm run preview      # 预览构建结果
```
