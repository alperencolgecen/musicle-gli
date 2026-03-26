"""
MusicLe Engine — playlist.py
Manages song_list.txt: append, remove, reorder, read.
Also handles local file import (copy + metadata).
"""
import os
import shutil
from datetime import date


def append_song(list_path: str, filename: str, title: str, artist: str, duration: str):
    """Append a song entry to song_list.txt."""
    today = date.today().isoformat()
    entry = f"{filename}|{title}|{artist}|{today}|{duration}\n"
    with open(list_path, "a", encoding="utf-8") as f:
        f.write(entry)


def read_songs(list_path: str) -> list:
    """Read all songs from song_list.txt. Returns list of dicts."""
    if not os.path.isfile(list_path):
        return []
    songs = []
    with open(list_path, "r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            parts = line.split("|", 4)
            if len(parts) == 5:
                songs.append({
                    "filename": parts[0],
                    "title": parts[1],
                    "artist": parts[2],
                    "date_added": parts[3],
                    "duration": parts[4],
                })
    return songs


def remove_song(list_path: str, filename: str) -> dict:
    """Remove a song entry from song_list.txt by filename."""
    if not os.path.isfile(list_path):
        return {"status": "error", "error": "song_list.txt not found"}
    songs = read_songs(list_path)
    original_count = len(songs)
    songs = [s for s in songs if s["filename"] != filename]
    if len(songs) == original_count:
        return {"status": "error", "error": f"Song not found: {filename}"}
    _write_songs(list_path, songs)
    return {"status": "ok", "removed": filename}


def _write_songs(list_path: str, songs: list):
    """Write all songs back to song_list.txt."""
    with open(list_path, "w", encoding="utf-8") as f:
        for s in songs:
            f.write(f"{s['filename']}|{s['title']}|{s['artist']}|{s['date_added']}|{s['duration']}\n")


def add_local_file(source_path: str, playlist_dir: str) -> dict:
    """
    Copy an audio file into the playlist directory, extract metadata,
    and append to song_list.txt.
    Returns structured result dict.
    """
    if not os.path.isfile(source_path):
        return {"status": "error", "error": f"File not found: {source_path}"}

    ext = os.path.splitext(source_path)[1].lower()
    allowed = {".mp3", ".flac", ".m4a", ".aac", ".ogg", ".wav", ".opus"}
    if ext not in allowed:
        return {"status": "error", "error": f"Unsupported format: {ext}"}

    os.makedirs(playlist_dir, exist_ok=True)
    filename = os.path.basename(source_path)
    dest_path = os.path.join(playlist_dir, filename)

    # Handle filename collision
    if os.path.abspath(source_path) != os.path.abspath(dest_path):
        base, e = os.path.splitext(filename)
        counter = 1
        while os.path.exists(dest_path):
            filename = f"{base}_{counter}{e}"
            dest_path = os.path.join(playlist_dir, filename)
            counter += 1
        shutil.copy2(source_path, dest_path)

    # Extract metadata
    try:
        from metadata import extract_metadata
        meta = extract_metadata(dest_path)
    except Exception:
        meta = {
            "title": os.path.splitext(filename)[0],
            "artist": "Unknown",
            "duration": 0.0,
        }

    duration = meta.get("duration", 0.0)
    s = int(duration)
    dur_str = f"{s // 60:02d}:{s % 60:02d}"

    list_path = os.path.join(playlist_dir, "song_list.txt")
    append_song(list_path, filename, meta.get("title", filename), meta.get("artist", "Unknown"), dur_str)

    return {
        "status": "ok",
        "filename": filename,
        "title": meta.get("title", filename),
        "artist": meta.get("artist", "Unknown"),
        "duration": duration,
        "art_path": meta.get("art_path", ""),
    }
