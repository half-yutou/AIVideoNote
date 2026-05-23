"""
fast-whisper 语音转录器

模型选择 (环境变量 WHISPER_MODEL):
  tiny           — 模型最小，速度最快，适合实时性要求高但容忍低精度的场景
  base           — 体积小、速度快，精度比 tiny 略高 (默认)
  small          — 平衡精度与性能，适用于大多数普通转写任务
  medium         — 精度更高，适合需要较强识别能力的场景
  large-v1       — 精度极高
  large-v2       — 在 v1 基础上改进，识别更稳定
  large-v3       — 当前主流大模型，推荐
  large-v3-turbo — 速度优化版，保证精度的同时提升处理效率

 ⚠️ 模型越大，VRAM 需求越高。RTX 5060 8GB 推荐 base/small/medium。

GPU 配置 (环境变量):
  WHISPER_DEVICE=auto   — 自动检测 CUDA/CPU (默认)
  WHISPER_COMPUTE=auto  — float16 或 int8 推理 (默认)

beam_size 说明:
  beam_size=1 (贪心搜索) — 每步取概率最高的一个候选，速度快，中文日常对话精度足够
  beam_size=5 (束搜索)   — 保留 5 个候选序列推理，质量略好但慢约 5 倍
  选择 1 是因为中文转录 beam_size=1 vs 5 差异极小，但速度差几倍，不值得
"""

import logging
import os
from pathlib import Path

from faster_whisper import WhisperModel

from app.transcriber.base import BaseTranscriber, TranscriptResult, Segment

logger = logging.getLogger(__name__)

SUPPORTED_MODELS = [
    "tiny", "base", "small", "medium",
    "large-v1", "large-v2", "large-v3", "large-v3-turbo",
]


class WhisperTranscriber(BaseTranscriber):
    def __init__(self, model_size: str = "base", device: str = "auto", compute_type: str = "auto"):
        if model_size not in SUPPORTED_MODELS:
            logger.warning(f"未知模型 {model_size}，回退到 base")
            model_size = "base"
        self.model_size = model_size
        self.device = device
        self.compute_type = compute_type
        self._model: WhisperModel | None = None

    def name(self) -> str:
        return f"fast-whisper-{self.model_size}"

    def _get_model(self) -> WhisperModel:
        if self._model is None:
            logger.info(f"加载 Whisper 模型: {self.model_size} (device={self.device}, compute={self.compute_type})")
            self._model = WhisperModel(
                self.model_size,
                device=self.device,
                compute_type=self.compute_type,
            )
        return self._model

    def transcribe(self, audio_path: str) -> TranscriptResult:
        path = Path(audio_path)
        if not path.exists():
            raise FileNotFoundError(f"音频文件不存在: {audio_path}")

        logger.info(f"开始转录: {audio_path} (GPU)")
        model = self._get_model()

        beam_size = int(os.getenv("WHISPER_BEAM_SIZE", "1"))
        # beam_size=1 贪心搜索：每步只取概率最高的候选，速度快且中文日常对话精度足够
        # 与 beam_size=5 相比质量差异极小，但速度差几倍，因此默认 1
        segments_raw, info = model.transcribe(str(path), beam_size=beam_size, language="zh")

        segments = []
        full_text_parts = []

        for seg in segments_raw:
            segments.append(Segment(start=seg.start, end=seg.end, text=seg.text.strip()))
            full_text_parts.append(seg.text.strip())

        full_text = " ".join(full_text_parts)
        language = info.language if info else ""

        logger.info(f"转录完成: {len(segments)} 段, 语言={language}, 长度={len(full_text)}")

        return TranscriptResult(
            language=language,
            full_text=full_text,
            segments=segments,
        )
