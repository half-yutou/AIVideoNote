import logging
from typing import Optional

from app.transcriber.base import BaseTranscriber, TranscriptResult, Segment

logger = logging.getLogger(__name__)


class MLXWhisperTranscriber(BaseTranscriber):
    # mlx-community 仓库命名不统一，需要映射
    MODEL_REPO_MAP = {
        "tiny": "mlx-community/whisper-tiny",
        "base": "mlx-community/whisper-base-mlx",
        "small": "mlx-community/whisper-small-mlx",
        "medium": "mlx-community/whisper-medium-mlx",
        "large-v3": "mlx-community/whisper-large-v3-mlx",
        "large-v3-turbo": "mlx-community/whisper-large-v3-turbo",
    }

    def __init__(self, model_size: str = "base"):
        self.model_size = model_size

    def name(self) -> str:
        return f"mlx-whisper-{self.model_size}"

    def _get_repo(self) -> str:
        return self.MODEL_REPO_MAP.get(self.model_size, f"mlx-community/whisper-{self.model_size}-mlx")

    def transcribe(self, audio_path: str, language: Optional[str] = None) -> TranscriptResult:
        try:
            from mlx_whisper import transcribe as mlx_transcribe
        except ImportError:
            raise ImportError("mlx-whisper 未安装，仅在 Apple Silicon 上可用")

        repo = self._get_repo()
        logger.info(f"MLX Whisper 转录: model={self.model_size}, repo={repo}, language={language or 'auto'}")
        kwargs = {"path_or_hf_repo": repo}
        if language:
            kwargs["language"] = language
        result = mlx_transcribe(audio_path, **kwargs)

        segments = []
        for seg in result.get("segments", []):
            segments.append(Segment(
                start=seg.get("start", 0),
                end=seg.get("end", 0),
                text=seg.get("text", "").strip(),
            ))

        logger.info(f"MLX 转录完成: {len(segments)} 段")
        return TranscriptResult(
            language=result.get("language", ""),
            full_text=result.get("text", ""),
            segments=segments,
        )
