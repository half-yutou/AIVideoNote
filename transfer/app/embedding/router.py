from fastapi import APIRouter

router = APIRouter(prefix="/embedding", tags=["embedding"])


@router.get("/health")
def health():
    return {"available": False, "message": "嵌入模块开发中"}
