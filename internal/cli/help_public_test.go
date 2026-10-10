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
	"io"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/retr0h/psion/internal/cli"
)

type HelpPublicTestSuite struct {
	suite.Suite
}

func (s *HelpPublicTestSuite) TestRender() {
	tests := []struct {
		name     string
		help     cli.Help
		mode     cli.ColorMode
		fail     bool
		contains []string
		absent   []string
		wantANSI bool
	}{
		{
			name: "root banner, commands, flags, and footer",
			help: cli.Help{
				Name:        "name replaced by banner",
				Banner:      true,
				Description: "Inspect the input.\nKeep the original data.",
				Usage:       "tool <command> [flags]",
				Commands:    []cli.Item{{Name: "inspect", Description: "inspect the input"}},
				Flags:       []cli.Item{{Name: "--color mode", Description: "terminal color"}},
				Footer:      "Run tool <command> --help for more.",
			},
			mode: "never",
			contains: []string{
				"\n  █", "  Inspect the input.\n  Keep the original data.",
				"USAGE\n    tool <command> [flags]", "COMMANDS\n    inspect",
				"FLAGS\n    --color mode", "terminal color",
				"Run tool <command> --help for more.",
			},
			absent: []string{"name replaced by banner"},
		},
		{
			name: "subcommand title and unicode column widths",
			help: cli.Help{
				Name:  "tool inspect",
				Usage: "tool inspect [flags]",
				Flags: []cli.Item{
					{Name: "--路径", Description: "input path"},
					{Name: "--color", Description: "terminal color"},
				},
			},
			mode:     "auto",
			contains: []string{"tool inspect", "--路径    input path", "--color   terminal color"},
			absent:   []string{"COMMANDS", "█"},
		},
		{
			name:     "optional sections are omitted",
			help:     cli.Help{Usage: "tool"},
			mode:     "never",
			contains: []string{"USAGE\n    tool"},
			absent:   []string{"COMMANDS", "FLAGS", "█"},
		},
		{
			name: "forced color is honored by the help renderer",
			help: cli.Help{
				Name:     "tool inspect",
				Usage:    "tool inspect [flags]",
				Commands: []cli.Item{{Name: "inspect", Description: "inspect the input"}},
				Footer:   "More help is available.",
			},
			mode:     "always",
			contains: []string{"tool inspect", "USAGE", "COMMANDS", "More help is available."},
			wantANSI: true,
		},
		{
			name: "write failure is returned to the caller",
			help: cli.Help{Usage: "tool"},
			mode: "never",
			fail: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var out bytes.Buffer
			var destination io.Writer = &out
			if tc.fail {
				destination = failingWriter{}
			}
			err := tc.help.Render(destination, cli.NewTheme(destination, tc.mode))
			if tc.fail {
				s.Require().ErrorIs(err, io.ErrClosedPipe)
				return
			}
			s.Require().NoError(err)
			for _, want := range tc.contains {
				s.Contains(out.String(), want)
			}
			for _, unwanted := range tc.absent {
				s.NotContains(out.String(), unwanted)
			}
			if tc.wantANSI {
				s.Contains(out.String(), "\x1b[")
			} else {
				s.NotContains(out.String(), "\x1b")
			}
		})
	}
}

// A handwritten double is appropriate for the standard library's io.Writer.
type failingWriter struct{}

func (failingWriter) Write(
	[]byte,
) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestHelpPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(HelpPublicTestSuite))
}
