# AIVideoNote 后端 — Go API 服务

## 技术栈

| 组件 | 库 | 用途 |
|------|-----|------|
| HTTP 框架 | `gin-gonic/gin` | 路由、中间件、请求处理 |
| ORM | `gorm.io/gorm` + `gorm.io/driver/sqlite` | 数据库操作 |
| 配置 | `spf13/viper` + `joho/godotenv` | YAML + .env 双重配置 |
| 日志 | `go.uber.org/zap` | 结构化日志 |
| SQLite | `glebarez/sqlite` | 纯 Go SQLite 驱动（WAL 模式） |
| UUID | `google/uuid` | 资源 ID 生成 |

## 目录结构

```
backend/
├── cmd/server/main.go              # 入口：初始化 → 注册路由 → 启动
├── internal/
│   ├── config/config.go            # 配置加载（Viper + godotenv）
│   ├── database/database.go        # SQLite 连接 + AutoMigrate
│   ├── model/                      # GORM 数据模型
│   │   ├── task.go                 # 任务状态 & 常量
│   │   ├── provider.go             # LLM 提供商
│   │   ├── note.go                 # 笔记记录
│   │   └── group.go                # 任务分组
│   ├── repository/                 # 数据访问层
│   │   ├── task_repo.go
│   │   ├── provider_repo.go
│   │   ├── note_repo.go
│   │   └── group_repo.go
│   ├── handler/                    # HTTP Handler（薄层）
│   │   ├── task.go                 # 任务生成/状态/列表/重命名
│   │   ├── provider.go             # 提供商 CRUD
│   │   ├── group.go                # 分组 CRUD
│   │   ├── chat.go                 # AI 问答（后端保留，前端已移除）
│   │   ├── upload.go               # 文件上传
│   │   └── export.go               # 笔记导出
│   ├── middleware/                  # 中间件（CORS / Logger）
│   ├── service/                    # 业务逻辑
│   │   ├── task/pipeline.go        # 任务管道（解析→下载→转录→LLM→后处理）
│   │   ├── task/executor.go        # 任务调度器
│   │   ├── downloader/             # 下载器（B站/YouTube/抖音/快手/本地）
│   │   ├── transcriber/client.go   # HTTP 调用 Python 转录服务
│   │   ├── llm/                    # LLM 客户端 + Prompt 模板
│   │   └── media/ffmpeg.go         # FFmpeg 封装
│   └── pkg/                        # 工具包
│       ├── response/               # 统一 JSON 响应格式
│       ├── errcode/                # 错误码
│       └── logger/                 # Zap 初始化
├── data/                           # SQLite 数据库文件
│   └── aivideonote.db                # 主数据库（WAL 模式）
├── config.yaml                     # 默认配置
└── go.mod
```

## API 路由

所有接口前缀 `/api/v1`。

### 任务

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/note/generate` | 提交笔记生成任务 |
| `GET` | `/task/:id/status` | 查询任务状态 |
| `GET` | `/task/list` | 任务列表 |
| `DELETE` | `/task/:id` | 删除任务 |
| `PUT` | `/task/:id/name` | 重命名任务 |

### 分组

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/group/list` | 分组列表 |
| `POST` | `/group/create` | 创建分组 |
| `PUT` | `/group/:id/name` | 重命名分组 |
| `DELETE` | `/group/:id` | 删除分组 |

### 提供商 & 导出

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/provider` | 添加 LLM 提供商 |
| `GET` | `/provider` | 提供商列表 |
| `GET` | `/provider/:id` | 提供商详情 |
| `PUT` | `/provider/:id` | 更新提供商 |
| `DELETE` | `/provider/:id` | 删除提供商 |
| `GET` | `/model/list` | 获取可用模型列表 |
| `GET` | `/export/:id` | 导出笔记文件 |

## 任务流水线

```
PARSING → DOWNLOADING → TRANSCRIBING → GENERATING → POST_PROCESSING → SUCCESS
   ↓           ↓              ↓             ↓              ↓
   _FAILED     _FAILED        _FAILED       _FAILED        _FAILED
```

每个步骤出错时记录具体失败原因和阶段标记（如 `TRANSCRIBING_FAILED`），前端可精确展示失败位置。

## 数据存储分离

- **数据库**：`backend/data/aivideonote.db` — SQLite WAL 模式，后端独享
- **运行时产出**：`../data/{video_id}/` — 音频 `.mp3`、转录中间文件 `_transcript.json`、笔记 `.md`
- 二者分离，避免用户文件污染后端工作目录

## 配置

两套配置，分工明确：

**`config.yaml`** — 内部不可见配置（仅 `server`、`database`、`task`）：

```yaml
server:
  port: 8080          # 可被 .env BACKEND_PORT 覆盖
  host: "0.0.0.0"

database:
  file_path: "./data/aivideonote.db"

task:
  max_concurrent: 3
  max_retry: 3
  retry_delay: 30
```

**`.env`** — 用户可配置项（项目根目录）：

```bash
BACKEND_PORT=8080
TRANSCRIBER_PORT=9090
FFMPEG_PATH=ffmpeg
YT_DLP_PATH=yt-dlp
DATA_DIR=../data
LLM_DEFAULT_BASE_URL=https://api.openai.com/v1
```

`config.go` 启动时通过 `godotenv` 加载 `.env`，`applyEnvSections()` 将环境变量注入 `Tools`/`Storage`/`LLM`/`PythonService` 结构体。完整的变量列表见项目根目录 `README.md`。

## 启动

```bash
cd backend

# 开发模式直接运行
make run

# 编译到 bin/ 目录后运行
make build
./bin/aivideonote.exe
```
