import logging
import os
from pathlib import Path

import uvicorn

_PROJECT_ROOT = Path(__file__).resolve().parent.parent
_DOTENV = _PROJECT_ROOT / ".env"
if _DOTENV.exists():
    from dotenv import load_dotenv
    load_dotenv(_DOTENV)

_THIS_DIR = Path(__file__).resolve().parent


def _setup_logging():
    log_dir = _THIS_DIR / "log"
    log_dir.mkdir(parents=True, exist_ok=True)
    log_file = log_dir / "transcriber.log"

    fmt = logging.Formatter("%(asctime)s %(levelname)s [%(name)s] %(message)s")

    stream_handler = logging.StreamHandler()
    stream_handler.setFormatter(fmt)

    file_handler = logging.FileHandler(str(log_file), encoding="utf-8")
    file_handler.setFormatter(fmt)

    root_logger = logging.getLogger()
    root_logger.setLevel(logging.INFO)
    root_logger.addHandler(stream_handler)
    root_logger.addHandler(file_handler)

    root_logger.info("日志文件: %s", log_file)


def main():
    _setup_logging()
    port = int(os.getenv("TRANSCRIBER_PORT", "9090"))
    uvicorn.run("app.main:app", host="0.0.0.0", port=port, log_level="info")


if __name__ == "__main__":
    main()
