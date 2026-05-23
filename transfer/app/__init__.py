import os
import site
from pathlib import Path

_project_root = Path(__file__).resolve().parent.parent.parent
_dotenv_path = _project_root / ".env"
if _dotenv_path.exists():
    from dotenv import load_dotenv
    load_dotenv(_dotenv_path)

_venv_site = next((p for p in site.getsitepackages() if p.endswith("site-packages")), "")
if _venv_site:
    for lib_dir in [
        "nvidia/cublas/bin",
        "nvidia/cuda_nvrtc/bin",
    ]:
        _dll_dir = Path(_venv_site) / lib_dir
        if _dll_dir.is_dir():
            os.add_dll_directory(str(_dll_dir))
