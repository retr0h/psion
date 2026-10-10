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

package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

type colorMode string

func (mode *colorMode) String() string { return string(*mode) }

func (*colorMode) Type() string { return "mode" }

func (mode *colorMode) Set(value string) error {
	switch value {
	case "auto", "always", "never":
		*mode = colorMode(value)
		return nil
	default:
		return fmt.Errorf("--color must be auto, always, or never")
	}
}

type cliStyle struct {
	accent, title, muted, success, warning, failure lipgloss.Style
}

func colorEnabled(out io.Writer) bool {
	if color == "always" {
		return true
	}
	if color == "never" || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	file, ok := out.(*os.File)
	return ok && term.IsTerminal(int(file.Fd()))
}

func newCLIStyle(out io.Writer) cliStyle {
	renderer := lipgloss.NewRenderer(out)
	profile := termenv.Ascii
	if colorEnabled(out) {
		profile = termenv.TrueColor
	}
	renderer.SetColorProfile(profile)
	fg := func(value string) lipgloss.Style {
		return renderer.NewStyle().Foreground(lipgloss.Color(value))
	}
	return cliStyle{
		accent:  fg("#bc8cff"),
		title:   fg("#bc8cff").Bold(true),
		muted:   fg("#8b949e"),
		success: fg("#3fb950"),
		warning: fg("#d29922"),
		failure: fg("#f85149").Bold(true),
	}
}

func (style cliStyle) phase(value string) string {
	switch value {
	case "Succeeded":
		return style.success.Render(value)
	case "Failed":
		return style.failure.Render(value)
	case "Pending", "Unknown":
		return style.warning.Render(value)
	default:
		return style.muted.Render(value)
	}
}

func styleHelp(root *cobra.Command) {
	root.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		out := cmd.OutOrStdout()
		// Cobra's help hook cannot return write errors.
		_ = renderHelp(cmd, out, out, cmd == root)
	})
	root.SetUsageFunc(func(cmd *cobra.Command) error {
		// UsageString captures output in a buffer. Detect color on the actual
		// diagnostic destination, not on that temporary buffer.
		return renderHelp(cmd, cmd.OutOrStdout(), cmd.ErrOrStderr(), false)
	})
}

type helpItem struct {
	name, description string
}

func renderHelp(cmd *cobra.Command, out, terminal io.Writer, banner bool) error {
	style := newCLIStyle(terminal)
	var page strings.Builder
	if banner {
		fmt.Fprintf(
			&page,
			"\n  %s\n  %s\n",
			style.muted.Render(bannerTop),
			style.accent.Render(bannerBottom),
		)
	} else {
		fmt.Fprintf(&page, "\n  %s\n", style.title.Render(cmd.CommandPath()))
	}
	description := cmd.Long
	if description == "" {
		description = cmd.Short
	}
	fmt.Fprintf(&page, "\n  %s\n", strings.ReplaceAll(strings.TrimSpace(description), "\n", "\n  "))
	usage := cmd.UseLine()
	if cmd.HasAvailableSubCommands() {
		usage = cmd.CommandPath() + " <command> [flags]"
	}
	fmt.Fprintf(&page, "\n  %s\n    %s\n", style.title.Render("USAGE"), usage)
	var commands []helpItem
	for _, child := range cmd.Commands() {
		if child.IsAvailableCommand() {
			commands = append(commands, helpItem{child.Name(), child.Short})
		}
	}
	writeHelpItems(&page, style, "COMMANDS", commands)
	cmd.InitDefaultHelpFlag()
	var flags []helpItem
	for _, set := range []*pflag.FlagSet{cmd.LocalFlags(), cmd.InheritedFlags()} {
		set.VisitAll(func(flag *pflag.Flag) {
			if flag.Hidden {
				return
			}
			name := "--" + flag.Name
			if flag.Shorthand != "" {
				name = "-" + flag.Shorthand + ", " + name
			}
			if flag.Value.Type() != "bool" {
				name += " " + flag.Value.Type()
			}
			description := flag.Usage
			if flag.DefValue != "" && flag.DefValue != "false" {
				description += " (default " + flag.DefValue + ")"
			}
			flags = append(flags, helpItem{name, description})
		})
	}
	sort.Slice(flags, func(i, j int) bool { return flags[i].name < flags[j].name })
	writeHelpItems(&page, style, "FLAGS", flags)
	if len(commands) > 0 {
		fmt.Fprintf(
			&page,
			"\n  %s\n",
			style.muted.Render("Run \""+cmd.CommandPath()+" <command> --help\" for more."),
		)
	}
	page.WriteByte('\n')
	_, err := io.WriteString(out, page.String())
	return err
}

func writeHelpItems(page *strings.Builder, style cliStyle, title string, items []helpItem) {
	if len(items) == 0 {
		return
	}
	width := 0
	for _, item := range items {
		if len(item.name) > width {
			width = len(item.name)
		}
	}
	fmt.Fprintf(page, "\n  %s\n", style.title.Render(title))
	for _, item := range items {
		fmt.Fprintf(page, "    %s%s%s\n", style.accent.Render(item.name),
			strings.Repeat(" ", width-len(item.name)+3), style.muted.Render(item.description))
	}
}
