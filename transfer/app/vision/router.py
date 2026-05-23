from fastapi import APIRouter

router = APIRouter(prefix="/vision", tags=["vision"])


@router.get("/health")
def health():
    return {"available": False, "message": "视觉模块开发中"}
