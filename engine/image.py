"""
MusicLe Engine — image.py
Image processing for CLI display.
- Center-crop resize for avatar/playlist/album art slots
- ANSI block-character rendering for terminal display
"""
import os

try:
    from PIL import Image
    PIL_AVAILABLE = True
except ImportError:
    PIL_AVAILABLE = False


def center_crop(image_path: str, target_w: int, target_h: int, output_path: str) -> dict:
    """
    Crop and resize an image to target_w × target_h (pixels) with center-crop.
    Saves result to output_path. Returns dict with status.
    """
    if not PIL_AVAILABLE:
        return {"status": "error", "error": "Pillow not installed. Run: pip install Pillow"}

    if not os.path.isfile(image_path):
        return {"status": "error", "error": f"Image not found: {image_path}"}

    try:
        os.makedirs(os.path.dirname(output_path) or ".", exist_ok=True)
        with Image.open(image_path) as img:
            img = img.convert("RGB")
            img_w, img_h = img.size
            target_ratio = target_w / target_h
            img_ratio = img_w / img_h

            # Scale so that the shortest dimension fills the target
            if img_ratio > target_ratio:
                # Too wide — fit by height
                new_h = target_h
                new_w = int(new_h * img_ratio)
            else:
                # Too tall — fit by width
                new_w = target_w
                new_h = int(new_w / img_ratio)

            img = img.resize((new_w, new_h), Image.LANCZOS)
            left = (new_w - target_w) // 2
            top = (new_h - target_h) // 2
            img = img.crop((left, top, left + target_w, top + target_h))
            img.save(output_path, "JPEG", quality=85)

        return {"status": "ok", "path": output_path}
    except Exception as e:
        return {"status": "error", "error": str(e)}


def image_to_ansi(image_path: str, width: int = 20) -> dict:
    """
    Convert an image to ANSI block characters for terminal rendering.
    Uses the upper-half-block (▀) trick: each character cell represents
    2 vertical pixels, with fg=top pixel color, bg=bottom pixel color.
    Returns a dict with 'ansi' key containing the renderable string.
    """
    if not PIL_AVAILABLE:
        return _placeholder_art(width)

    if not os.path.isfile(image_path):
        return _placeholder_art(width)

    try:
        with Image.open(image_path) as img:
            img = img.convert("RGB")
            # Maintain aspect ratio; height = width/2 * 2 (2 pixels per row)
            h = max(4, width // 2)
            render_h = h * 2  # 2 pixels per terminal row
            img = img.resize((width, render_h), Image.LANCZOS)
            pixels = list(img.getdata())

        lines = []
        for row in range(0, render_h, 2):
            line = ""
            for col in range(width):
                r1, g1, b1 = pixels[row * width + col]
                r2, g2, b2 = pixels[(row + 1) * width + col] if row + 1 < render_h else (0, 0, 0)
                # fg = top, bg = bottom → ▀
                line += f"\x1b[38;2;{r1};{g1};{b1}m\x1b[48;2;{r2};{g2};{b2}m▀"
            line += "\x1b[0m"
            lines.append(line)

        return {"status": "ok", "ansi": "\n".join(lines), "width": width, "height": h}

    except Exception as e:
        return _placeholder_art(width, error=str(e))


def _placeholder_art(width: int, error: str = "") -> dict:
    """Return a simple block-character placeholder when no image is available."""
    h = max(4, width // 2)
    top_border = "  ╔" + "═" * (width - 2) + "╗"
    mid_line = "  ║" + " " * (width - 2) + "║"
    note_line = "  ║" + " " * ((width - 4) // 2) + "♫♪♫" + " " * ((width - 3) // 2) + "║"
    bot_border = "  ╚" + "═" * (width - 2) + "╝"

    inner_count = h - 2
    mid_lines = [note_line if i == inner_count // 2 else mid_line for i in range(inner_count)]
    art_lines = [top_border] + mid_lines + [bot_border]

    return {
        "status": "ok" if not error else "error",
        "ansi": "\n".join(art_lines),
        "width": width,
        "height": h,
        "error": error,
    }
