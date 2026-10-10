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

package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Item is a command or flag and its description.
type Item struct {
	Name, Description string
}

// Help describes a help page without depending on a command framework.
type Help struct {
	Name, Description, Usage, Footer string
	Commands, Flags                  []Item
	Banner                           bool
}

// Render writes a help page using a theme bound to its diagnostic destination.
func (help Help) Render(
	out io.Writer,
	theme Theme,
) error {
	var page strings.Builder
	if help.Banner {
		fmt.Fprintf(&page, "\n  %s\n", strings.ReplaceAll(theme.Banner(), "\n", "\n  "))
	} else if help.Name != "" {
		fmt.Fprintf(&page, "\n  %s\n", theme.Title(help.Name))
	}
	if help.Description != "" {
		fmt.Fprintf(
			&page,
			"\n  %s\n",
			strings.ReplaceAll(strings.TrimSpace(help.Description), "\n", "\n  "),
		)
	}
	if help.Usage != "" {
		fmt.Fprintf(&page, "\n  %s\n    %s\n", theme.Title("USAGE"), help.Usage)
	}
	writeHelpItems(&page, theme, "COMMANDS", help.Commands)
	writeHelpItems(&page, theme, "FLAGS", help.Flags)
	if help.Footer != "" {
		fmt.Fprintf(&page, "\n  %s\n", theme.Mute.Render(help.Footer))
	}
	page.WriteByte('\n')
	_, err := io.WriteString(out, page.String())
	return err
}

func writeHelpItems(
	page *strings.Builder,
	theme Theme,
	title string,
	items []Item,
) {
	if len(items) == 0 {
		return
	}
	width := 0
	for _, item := range items {
		width = max(width, lipgloss.Width(item.Name))
	}
	fmt.Fprintf(page, "\n  %s\n", theme.Title(title))
	for _, item := range items {
		fmt.Fprintf(
			page,
			"    %s%s%s\n",
			theme.Accent.Render(item.Name),
			strings.Repeat(
				" ",
				width-lipgloss.Width(item.Name)+3,
			),
			theme.Mute.Render(item.Description),
		)
	}
}
