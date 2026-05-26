import os
import logging
import time
from pathlib import Path
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Optional

logger = logging.getLogger(__name__)


@dataclass
class Segment:
    start: float
    end: float
    text: str


@dataclass
class TranscriptResult:
    language: str = ""
    full_text: str = ""
    segments: list[Segment] = field(default_factory=list)


class BaseTranscriber(ABC):
    @abstractmethod
    def transcribe(self, audio_path: str, language: Optional[str] = None) -> TranscriptResult:
        """转录音频文件。language 为 None 时自动检测语言。"""
        pass

    @abstractmethod
    def name(self) -> str:
        pass


def detect_transcriber() -> str:
    env_type = os.getenv("TRANSCRIBER_TYPE", "").strip().lower()
    if env_type in ("fast-whisper", "faster-whisper", "whisper", ""):
        return "fast-whisper"
    if env_type == "groq":
        return "groq"
    if env_type == "mlx-whisper":
        return "mlx-whisper"
    return "fast-whisper"


def create_transcriber() -> BaseTranscriber:
    t = detect_transcriber()

    if t == "fast-whisper":
        from app.transcriber.whisper import WhisperTranscriber
        model_size = os.getenv("WHISPER_MODEL", "base")
        device = os.getenv("WHISPER_DEVICE", "auto")
        compute_type = os.getenv("WHISPER_COMPUTE", "auto")
        logger.info(f"使用 fast-whisper: model={model_size} device={device} compute_type={compute_type}")
        return WhisperTranscriber(model_size=model_size, device=device, compute_type=compute_type)

    if t == "groq":
        from app.transcriber.groq import GroqTranscriber
        logger.info("使用 Groq API 转录")
        return GroqTranscriber()

    if t == "mlx-whisper":
        from app.transcriber.mlx_whisper import MLXWhisperTranscriber
        model_size = os.getenv("WHISPER_MODEL", "base")
        logger.info(f"使用 mlx-whisper: model={model_size}")
        return MLXWhisperTranscriber(model_size=model_size)

    logger.warning(f"未知转写器类型: {t}，回退到 fast-whisper")
    from app.transcriber.whisper import WhisperTranscriber
    return WhisperTranscriber()
