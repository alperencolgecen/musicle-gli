"""
MusicLe Engine — download.py
Downloads music from YouTube (yt-dlp) and Spotify (spotdl).
After downloading, extracts metadata and returns structured result.
"""
import os
import subprocess
import sys
import json


def _find_exe(names: list) -> str:
    """Return the first executable name found in PATH."""
    import shutil
    for name in names:
        if shutil.which(name):
            return name
    return names[0]  # fallback, will fail with a clear error


def download_youtube(url: str, output_dir: str) -> dict:
    """Download a YouTube URL using yt-dlp, extract audio as mp3."""
    if not url.startswith("http"):
        return {"status": "error", "error": "Invalid URL"}

    os.makedirs(output_dir, exist_ok=True)
    ytdlp = _find_exe(["yt-dlp", "yt_dlp"])

    # Output template: title.mp3
    out_template = os.path.join(output_dir, "%(title)s.%(ext)s")

    cmd = [
        sys.executable, "-m", "yt_dlp" if ytdlp == "yt_dlp" else ytdlp,
        url,
        "--extract-audio",
        "--audio-format", "mp3",
        "--audio-quality", "192K",
        "--output", out_template,
        "--no-playlist",
        "--print", "after_move:filepath",
        "--quiet",
        "--no-warnings",
    ]
    # Prefer module invocation for reliability on Windows
    cmd = [sys.executable, "-m", "yt_dlp",
           url,
           "--extract-audio",
           "--audio-format", "mp3",
           "--audio-quality", "192K",
           "--output", out_template,
           "--no-playlist",
           "--print", "after_move:filepath",
           "--quiet",
           "--no-warnings",
           ]

    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=300)
        if result.returncode != 0:
            err = result.stderr.strip() or "yt-dlp failed"
            return {"status": "error", "error": err}

        # The last line of stdout is the output filepath
        filepath = result.stdout.strip().splitlines()[-1] if result.stdout.strip() else ""
        if not filepath or not os.path.isfile(filepath):
            # Try to find the latest mp3 in output dir
            filepath = _latest_file(output_dir, ".mp3")

        if not filepath:
            return {"status": "error", "error": "Downloaded file not found"}

        return _finalize_download(filepath, output_dir)

    except FileNotFoundError:
        return {"status": "error", "error": "yt-dlp not installed. Run: pip install yt-dlp"}
    except subprocess.TimeoutExpired:
        return {"status": "error", "error": "Download timed out"}
    except Exception as e:
        return {"status": "error", "error": str(e)}


def download_spotify(url: str, output_dir: str) -> dict:
    """Download a Spotify URL using spotdl."""
    if not url.startswith("http"):
        return {"status": "error", "error": "Invalid URL"}

    os.makedirs(output_dir, exist_ok=True)

    cmd = [
        sys.executable, "-m", "spotdl",
        url,
        "--output", output_dir,
        "--format", "mp3",
        "--bitrate", "192k",
    ]
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=300)
        if result.returncode != 0:
            err = result.stderr.strip() or "spotdl failed"
            return {"status": "error", "error": err}

        # Find newest mp3
        filepath = _latest_file(output_dir, ".mp3")
        if not filepath:
            return {"status": "error", "error": "Downloaded file not found"}

        return _finalize_download(filepath, output_dir)

    except FileNotFoundError:
        return {"status": "error", "error": "spotdl not installed. Run: pip install spotdl"}
    except subprocess.TimeoutExpired:
        return {"status": "error", "error": "Download timed out"}
    except Exception as e:
        return {"status": "error", "error": str(e)}


def _finalize_download(filepath: str, output_dir: str) -> dict:
    """Extract metadata after download and return result dict."""
    try:
        from metadata import extract_metadata
        meta = extract_metadata(filepath)
    except Exception:
        meta = {"title": os.path.splitext(os.path.basename(filepath))[0], "artist": "Unknown", "duration": 0.0}

    filename = os.path.basename(filepath)
    duration = meta.get("duration", 0.0)
    dur_str = _fmt_dur(duration)

    # Append to song_list.txt in output_dir
    try:
        from playlist import append_song
        list_path = os.path.join(output_dir, "song_list.txt")
        append_song(list_path, filename, meta.get("title", filename), meta.get("artist", "Unknown"), dur_str)
    except Exception:
        pass

    return {
        "status": "ok",
        "filename": filename,
        "title": meta.get("title", filename),
        "artist": meta.get("artist", "Unknown"),
        "duration": duration,
        "art_path": meta.get("art_path", ""),
    }


def _latest_file(directory: str, ext: str) -> str:
    """Return the most recently modified file with given extension in directory."""
    try:
        files = [
            os.path.join(directory, f)
            for f in os.listdir(directory)
            if f.lower().endswith(ext)
        ]
        if not files:
            return ""
        return max(files, key=os.path.getmtime)
    except Exception:
        return ""


def _fmt_dur(seconds: float) -> str:
    s = int(seconds)
    return f"{s // 60:02d}:{s % 60:02d}"
