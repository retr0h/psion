// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// Package cli owns terminal presentation, independently of Cobra and resource logic.
package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

const (
	bannerTop    = "█▀█ █▀ █ █▀█ █▄░█"
	bannerBottom = "█▀▀ ▄█ █ █▄█ █░▀█"
)

// ColorMode controls whether terminal output uses color.
type ColorMode string

func (mode *ColorMode) String() string { return string(*mode) }

// Type is the value type displayed in flag help.
func (*ColorMode) Type() string { return "mode" }

// Set validates the flag value before changing the mode.
func (mode *ColorMode) Set(
	value string,
) error {
	switch value {
	case "auto", "always", "never":
		*mode = ColorMode(value)
		return nil
	default:
		return fmt.Errorf("--color must be auto, always, or never")
	}
}

// Enabled detects color on the destination, honoring explicit flags and NO_COLOR.
func (mode ColorMode) Enabled(
	out io.Writer,
) bool {
	if mode == "always" {
		return true
	}
	if mode == "never" || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	file, ok := out.(*os.File)
	return ok && term.IsTerminal(int(file.Fd()))
}

// Theme gives presentation roles a shared Lip Gloss palette.
type Theme struct {
	Mute, Accent, OK, Err, Info, BannerTop, BannerBot lipgloss.Style
}

// NewTheme binds the palette to the actual terminal or output sink.
func NewTheme(
	out io.Writer,
	mode ColorMode,
) Theme {
	renderer := lipgloss.NewRenderer(out)
	profile := termenv.Ascii
	if mode.Enabled(out) {
		profile = termenv.TrueColor
	}
	renderer.SetColorProfile(profile)
	fg := func(value string) lipgloss.Style {
		return renderer.NewStyle().Foreground(lipgloss.Color(value))
	}
	return Theme{
		Mute:      fg("#8b949e"),
		Accent:    fg("#bc8cff"),
		OK:        fg("#3fb950"),
		Err:       fg("#f85149").Bold(true),
		Info:      fg("#d29922"),
		BannerTop: fg("#8b949e"),
		BannerBot: fg("#bc8cff"),
	}
}

// Title renders a heading in the project's accent color.
func (theme Theme) Title(
	value string,
) string {
	return theme.Accent.Bold(true).Render(value)
}

// Banner renders the same block wordmark as the SVG assets.
func (theme Theme) Banner() string {
	return theme.BannerTop.Render(bannerTop) + "\n" + theme.BannerBot.Render(bannerBottom)
}

// Phase renders a resource condition with the corresponding status color.
func (theme Theme) Phase(
	value string,
) string {
	switch value {
	case "Succeeded":
		return theme.OK.Render(value)
	case "Failed":
		return theme.Err.Render(value)
	case "Pending", "Unknown":
		return theme.Info.Render(value)
	default:
		return theme.Mute.Render(value)
	}
}
