from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.transcriber.router import router as transcriber_router
from app.vision.router import router as vision_router
from app.embedding.router import router as embedding_router

app = FastAPI(title="AIVideoNote Python Service", version="1.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(transcriber_router, prefix="/api/v1")
app.include_router(vision_router, prefix="/api/v1")
app.include_router(embedding_router, prefix="/api/v1")


@app.get("/health")
def health():
    return {"status": "ok"}
