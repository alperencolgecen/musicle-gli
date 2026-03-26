package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// Tcell color values
var (
	ColorBackground = tcell.NewRGBColor(18, 18, 18)    // #121212
	ColorSurface    = tcell.NewRGBColor(24, 24, 24)    // #181818
	ColorBorder     = tcell.NewRGBColor(40, 40, 40)    // #282828
	ColorAccent     = tcell.NewRGBColor(29, 185, 84)   // #1DB954
	ColorPrimary    = tcell.ColorWhite                  // #FFFFFF
	ColorSecondary  = tcell.NewRGBColor(179, 179, 179) // #B3B3B3
	ColorError      = tcell.NewRGBColor(255, 68, 68)   // #FF4444
	ColorOrange     = tcell.NewRGBColor(255, 165, 0)   // #FFA500
	ColorRowHover   = tcell.NewRGBColor(30, 215, 96)   // lighter green
	ColorBlack      = tcell.ColorBlack
)

// Tview inline color tags (used inside SetText / fmt strings)
const (
	TagAccent    = "[#1DB954]"
	TagWhite     = "[white]"
	TagSecondary = "[#B3B3B3]"
	TagError     = "[#FF4444]"
	TagOrange    = "[#FFA500]"
	TagBold      = "[::b]"
	TagReset     = "[-]"
	TagAttrReset = "[::-]"
)

// VolumeColor returns the tcell color for a given volume level (0.0–1.0)
func VolumeColor(vol float64) tcell.Color {
	switch {
	case vol <= 0.33:
		return ColorAccent
	case vol <= 0.66:
		return ColorOrange
	default:
		return ColorError
	}
}

// VolumeBar renders a progress bar string for display in text views
// width = total characters, filled = ratio 0.0–1.0
func VolumeBar(filled float64, width int) string {
	if width <= 0 {
		width = 10
	}
	n := int(filled * float64(width))
	if n < 0 {
		n = 0
	}
	if n > width {
		n = width
	}
	bar := ""
	for i := 0; i < n; i++ {
		bar += "█"
	}
	for i := n; i < width; i++ {
		bar += "░"
	}
	return bar
}

// ProgressBar returns a progress string: '──●───────'
func ProgressBar(pos, dur float64, width int) string {
	if width <= 0 {
		width = 20
	}
	ratio := 0.0
	if dur > 0 {
		ratio = pos / dur
	}
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	n := int(ratio * float64(width))
	bar := ""
	for i := 0; i < n; i++ {
		bar += "─"
	}
	bar += "●"
	for i := n + 1; i < width; i++ {
		bar += "─"
	}
	return bar
}

// FormatDuration converts seconds to MM:SS string
func FormatDuration(secs float64) string {
	s := int(secs)
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}
