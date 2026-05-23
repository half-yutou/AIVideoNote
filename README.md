# AIVideoNote

AI 视频笔记生成工具 — 输入视频链接，自动下载音频、语音转文字、AI 生成结构化 Markdown 笔记。

## 技术栈

| 层级     | 技术                                | 默认端口 |
| ------ | --------------------------------- | ---- |
| 前端     | Vue 3 + Vite + TypeScript + Pinia | 3015 |
| 后端     | Go (Gin + GORM + Viper + Zap)     | 8080 |
| 转录服务 | Python (FastAPI + faster-whisper) | 9090 |
| 数据库    | SQLite (WAL 模式)                   | —    |

## 快速开始

### 前置条件

- Go 1.24+
- Node.js 20+
- Python 3.11+（推荐使用 [uv](https://docs.astral.sh/uv/) 管理）
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)（视频下载）
- [ffmpeg](https://ffmpeg.org/)（音频处理）
- NVIDIA GPU + CUDA 12.x（用于 GPU 加速转录，否则会退化为 CPU 版本，速度会慢很多）

### 安装与运行

```bash
# 1. 克隆仓库
git clone <repo-url>
cd AIVideoNote

# 2. 复制环境配置
cp .env.example .env
# 编辑 .env 配置模型和设备参数

# 3. 启动转录服务（Python）
cd transfer
uv sync
uv run main.py

# 4. 启动后端（Go）
cd ../backend
go run ./cmd/server/

# 5. 启动前端（Vue）
cd ../frontend
npm install
npm run dev
```

### 环境变量

所有配置集中在项目根目录 `.env` 文件中，参考 `.env.example`：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `BACKEND_PORT` | Go 后端端口 | `8080` |
| `TRANSCRIBER_PORT` | Python 转录服务端口 | `9090` |
| `FRONTEND_PORT` | Vue 前端开发端口 | `3015` |
| `FFMPEG_PATH` | ffmpeg 路径 | `ffmpeg` |
| `FFPROBE_PATH` | ffprobe 路径 | `ffprobe` |
| `YT_DLP_PATH` | yt-dlp 路径 | `yt-dlp` |
| `TRANSCRIBE_TIMEOUT` | 转录超时（秒） | `600` |
| `WHISPER_MODEL` | Whisper 模型大小 | `base` |
| `WHISPER_DEVICE` | 推理设备 | `auto` |
| `WHISPER_COMPUTE` | 计算精度 | `auto` |
| `WHISPER_BEAM_SIZE` | 束搜索宽度 | `1` |
| `TRANSCRIBER_TYPE` | 转录引擎 | `fast-whisper` |
| `DATA_DIR` | 运行时产出目录 | `../data` |
| `UPLOAD_DIR` | 上传目录 | `../data/uploads` |
| `LLM_DEFAULT_BASE_URL` | LLM 默认 API 地址 | `https://api.openai.com/v1` |

## 项目结构

```
AIVideoNote/
├── backend/          # Go 后端服务
│   └── data/         # SQLite 数据库文件
├── frontend/         # Vue 3 前端
├── transfer/         # Python GPU 转录服务
├── data/             # 运行时产出（音频/转录/笔记，按视频分目录）
├── Docs/             # 项目文档
├── .env              # 环境变量（不提交）
├── .env.example      # 环境变量模板
└── README.md
```

## 文档

- [后端架构](backend/README.md)
- [前端架构](frontend/README.md)
- [GPU 服务架构](transfer/README.md)
- [详细文档](Docs/)

