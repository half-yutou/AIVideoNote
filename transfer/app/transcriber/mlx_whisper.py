import logging

from app.transcriber.base import BaseTranscriber, TranscriptResult, Segment

logger = logging.getLogger(__name__)


class MLXWhisperTranscriber(BaseTranscriber):
    def __init__(self, model_size: str = "base"):
        self.model_size = model_size

    def name(self) -> str:
        return f"mlx-whisper-{self.model_size}"

    def transcribe(self, audio_path: str) -> TranscriptResult:
        try:
            import mlx_whisper
        except ImportError:
            raise ImportError("mlx-whisper 未安装，仅在 Apple Silicon 上可用")

        logger.info(f"MLX Whisper 转录: model={self.model_size}")
        result = mlx_whisper.transcribe(audio_path, path_or_hf_repo=f"mlx-community/whisper-{self.model_size}")

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
