import os
import logging

from app.transcriber.base import BaseTranscriber, TranscriptResult, Segment

logger = logging.getLogger(__name__)


class GroqTranscriber(BaseTranscriber):
    def __init__(self):
        self.api_key = os.getenv("GROQ_API_KEY", "")

    def name(self) -> str:
        return "groq"

    def transcribe(self, audio_path: str) -> TranscriptResult:
        import requests

        if not self.api_key:
            raise ValueError("GROQ_API_KEY 环境变量未设置")

        logger.info(f"使用 Groq API 转录: {audio_path}")

        with open(audio_path, "rb") as f:
            response = requests.post(
                "https://api.groq.com/openai/v1/audio/transcriptions",
                headers={"Authorization": f"Bearer {self.api_key}"},
                files={"file": f},
                data={
                    "model": "whisper-large-v3",
                    "response_format": "verbose_json",
                },
                timeout=600,
            )

        if response.status_code != 200:
            raise RuntimeError(f"Groq API 错误: {response.status_code} {response.text}")

        data = response.json()

        segments = []
        for seg in data.get("segments", []):
            segments.append(Segment(
                start=seg.get("start", 0),
                end=seg.get("end", 0),
                text=seg.get("text", "").strip(),
            ))

        logger.info(f"Groq 转录完成: {len(segments)} 段")
        return TranscriptResult(
            language=data.get("language", ""),
            full_text=data.get("text", ""),
            segments=segments,
        )
