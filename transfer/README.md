# AIVideoNote Transfer — Python GPU 转录服务

独立 GPU 微服务，提供语音转文字（Whisper）、视觉理解、向量嵌入等能力。

## 技术栈

| 组件 | 库 | 用途 |
|------|-----|------|
| Web 框架 | FastAPI | HTTP API |
| 语音转录 | faster-whisper | CTranslate2 加速的 Whisper |
| GPU 加速 | nvidia-cublas-cu12 | CUDA 12 cuBLAS 支持 |
| 包管理 | uv | Python 环境 & 依赖管理 |
| 配置 | python-dotenv | 从项目根 .env 读取配置 |
| 服务 | uvicorn | ASGI 服务器 |

## 目录结构

```
transfer/
├── main.py                         # 开发入口（uvicorn server）
├── pyproject.toml                  # uv 项目配置 & 依赖
├── app/
│   ├── __init__.py                 # .env 加载 & CUDA DLL 注册
│   ├── main.py                     # FastAPI 应用定义 & 路由注册
│   ├── transcriber/
│   │   ├── __init__.py
│   │   ├── router.py               # 转录 API 路由
│   │   ├── base.py                 # BaseTranscriber 抽象基类
│   │   ├── whisper.py              # faster-whisper 实现
│   │   ├── groq.py                 # Groq API 转录（可选）
│   │   └── mlx_whisper.py          # Apple MLX 转录（可选）
│   ├── vision/
│   │   ├── __init__.py
│   │   └── router.py               # 视觉理解 API
│   └── embedding/
│       ├── __init__.py
│       └── router.py               # 向量嵌入 API
└── uv.lock                          # 依赖锁定文件
```

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/health` | 健康检查 |
| `GET` | `/api/v1/transcribe/health` | 转录服务健康检查 |
| `POST` | `/api/v1/transcribe` | 提交音频转录任务 |
| `POST` | `/api/v1/vision/*` | 视觉理解（预留） |
| `POST` | `/api/v1/embedding/*` | 向量嵌入（预留） |

### 转录请求格式

```json
{
  "audio_path": "/absolute/path/to/audio.mp3",
  "model": "large-v3"
}
```

### 转录响应格式

```json
{
  "language": "zh",
  "full_text": "...",
  "segments": [
    { "start": 0.0, "end": 3.5, "text": "..." }
  ]
}
```

## 模型配置

通过项目根目录 `.env` 配置：

| 变量 | 说明 | 可选值 |
|------|------|--------|
| `WHISPER_MODEL` | 模型大小 | tiny / base / small / medium / large-v3 |
| `WHISPER_DEVICE` | 推理设备 | auto / cpu / cuda |
| `WHISPER_COMPUTE` | 计算精度 | auto / float16 / int8 |
| `WHISPER_BEAM_SIZE` | 束搜索宽度 (1=greedy) | 1~5 |
| `TRANSCRIBER_TYPE` | 转录引擎 | whisper / groq / mlx |

## 模型大小参考

| 模型 | 显存需求 | 速度 | 精度 |
|------|---------|------|------|
| tiny | ~1 GB | 极快 | 较低 |
| base | ~1 GB | 快 | 一般 |
| small | ~2 GB | 较快 | 较好 |
| medium | ~5 GB | 中等 | 好 |
| large-v3 | ~10 GB | 慢 | 最好 |

## 启动

```bash
cd transfer

# 安装依赖
uv sync

# 开发模式启动（端口由 .env TRANSCRIBER_PORT 控制，默认 9090）
uv run main.py

# 生产模式（端口同样由 .env 控制）
uv run uvicorn app.main:app --host 0.0.0.0 --port $TRANSCRIBER_PORT
```
