#!/usr/bin/env python3
"""Render scripts/og-card.html to web/static/img/og.png (1200×630)."""

from __future__ import annotations

import shutil
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
HTML = ROOT / "scripts/og-card.html"
OUT = ROOT / "web/static/img/og.png"

CHROME_CANDIDATES = [
    "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
    "/Applications/Chromium.app/Contents/MacOS/Chromium",
    "google-chrome",
    "chromium",
    "chromium-browser",
]


def chrome_bin() -> str:
    for path in CHROME_CANDIDATES:
        if "/" in path and Path(path).exists():
            return path
        found = shutil.which(path)
        if found:
            return found
    raise SystemExit("Chrome/Chromium not found; install Chrome to render the Open Graph card")


def main() -> None:
    chrome = chrome_bin()
    OUT.parent.mkdir(parents=True, exist_ok=True)
    cmd = [
        chrome,
        "--headless=new",
        "--disable-gpu",
        "--hide-scrollbars",
        "--force-device-scale-factor=1",
        "--window-size=1200,630",
        "--virtual-time-budget=8000",
        f"--screenshot={OUT}",
        HTML.resolve().as_uri(),
    ]
    subprocess.run(cmd, check=True)
    if not OUT.exists() or OUT.stat().st_size < 1000:
        raise SystemExit(f"og screenshot missing or too small: {OUT}")
    print(f"wrote {OUT.relative_to(ROOT)} ({OUT.stat().st_size} bytes)", file=sys.stderr)


if __name__ == "__main__":
    main()
