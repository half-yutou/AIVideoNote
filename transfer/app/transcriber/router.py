import logging

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.transcriber.base import create_transcriber

logger = logging.getLogger(__name__)
router = APIRouter(tags=["transcribe"])

_transcriber = None


def get_transcriber():
    global _transcriber
    if _transcriber is None:
        _transcriber = create_transcriber()
    return _transcriber


class TranscribeRequest(BaseModel):
    audio_path: str


class SegmentModel(BaseModel):
    start: float
    end: float
    text: str


class TranscribeResponse(BaseModel):
    language: str
    full_text: str
    segments: list[SegmentModel]


@router.get("/transcribe/health")
def health():
    t = get_transcriber()
    return {"available": True, "transcriber": t.name()}


@router.post("/transcribe", response_model=TranscribeResponse)
def transcribe(req: TranscribeRequest):
    logger.info(f"收到转录请求: audio_path={req.audio_path}")
    try:
        t = get_transcriber()
        result = t.transcribe(req.audio_path)
        return TranscribeResponse(
            language=result.language,
            full_text=result.full_text,
            segments=[
                SegmentModel(start=s.start, end=s.end, text=s.text)
                for s in result.segments
            ],
        )
    except FileNotFoundError as e:
        raise HTTPException(status_code=404, detail=str(e))
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logger.exception("转录异常")
        raise HTTPException(status_code=500, detail=str(e))
