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

package cli_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/retr0h/psion/internal/cli"
)

func TestHelpOutput(t *testing.T) {
	var out bytes.Buffer
	help := cli.Help{
		Name:        "tool inspect",
		Description: "Inspect the input.\nKeep the original data.",
		Usage:       "tool inspect [flags]",
		Flags: []cli.Item{
			{Name: "--路径", Description: "input path"},
			{Name: "--color", Description: "terminal color"},
		},
	}
	if err := help.Render(&out, cli.NewTheme(&out, "auto")); err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, want := range []string{
		"tool inspect", "  Keep the original data.", "USAGE", "FLAGS",
		"--路径    input path", "--color   terminal color",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("help missing %q: %s", want, text)
		}
	}
	if strings.Contains(text, "COMMANDS") || strings.ContainsRune(text, '\x1b') {
		t.Errorf("plain help contains an empty section or ANSI escapes: %q", text)
	}
	if err := help.Render(failingWriter{}, cli.NewTheme(&out, "never")); !errors.Is(
		err,
		io.ErrClosedPipe,
	) {
		t.Fatalf("write error = %v, want %v", err, io.ErrClosedPipe)
	}
}

type failingWriter struct{}

func (failingWriter) Write(_ []byte) (int, error) { return 0, io.ErrClosedPipe }

func TestColorModes(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "xterm-256color")
	var out bytes.Buffer
	for _, tc := range []struct {
		mode cli.ColorMode
		ansi bool
	}{
		{mode: "auto", ansi: false},
		{mode: "never", ansi: false},
		{mode: "always", ansi: true},
	} {
		t.Run(string(tc.mode), func(t *testing.T) {
			theme := cli.NewTheme(&out, tc.mode)
			text := theme.Title("example") + theme.Banner()
			if got := strings.ContainsRune(text, '\x1b'); got != tc.ansi {
				t.Errorf("ANSI output = %v, want %v", got, tc.ansi)
			}
		})
	}
	var mode cli.ColorMode = "auto"
	if err := mode.Set("never"); err != nil || mode.String() != "never" {
		t.Fatalf("valid mode rejected: %v", err)
	}
	if err := mode.Set("invalid"); err == nil || mode.String() != "never" {
		t.Fatal("invalid mode should fail without changing the selected mode")
	}
	if mode.Type() != "mode" {
		t.Fatalf("flag type = %q", mode.Type())
	}
}
