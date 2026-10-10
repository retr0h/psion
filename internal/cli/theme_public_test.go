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
	"os"
	"regexp"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/psion/internal/cli"
)

type ThemePublicTestSuite struct {
	suite.Suite
}

func (s *ThemePublicTestSuite) TestSet() {
	tests := []struct {
		name    string
		value   string
		want    cli.ColorMode
		wantErr bool
	}{
		{name: "automatic color", value: "auto", want: "auto"},
		{name: "forced color", value: "always", want: "always"},
		{name: "disabled color", value: "never", want: "never"},
		{name: "invalid mode preserves selection", value: "invalid", want: "always", wantErr: true},
		{name: "empty mode is rejected", value: "", want: "always", wantErr: true},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			var mode cli.ColorMode = "always"
			err := mode.Set(tc.value)
			if tc.wantErr {
				s.Require().EqualError(err, "--color must be auto, always, or never")
			} else {
				s.Require().NoError(err)
			}
			s.Equal(tc.want, mode)
		})
	}
}

func (s *ThemePublicTestSuite) TestString() {
	tests := []struct {
		name string
		mode cli.ColorMode
		want string
	}{
		{name: "automatic", mode: "auto", want: "auto"},
		{name: "always", mode: "always", want: "always"},
		{name: "never", mode: "never", want: "never"},
		{name: "zero value"},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.Equal(tc.want, tc.mode.String())
		})
	}
}

func (s *ThemePublicTestSuite) TestType() {
	for _, value := range []cli.ColorMode{"auto", "always", "never"} {
		s.Run(string(value), func() {
			mode := value
			s.Equal("mode", mode.Type())
		})
	}
}

func (s *ThemePublicTestSuite) TestEnabled() {
	tests := []struct {
		name    string
		mode    cli.ColorMode
		noColor string
		term    string
		pipe    bool
		want    bool
	}{
		{
			name:    "explicit color overrides the environment",
			mode:    "always",
			noColor: "1",
			term:    "dumb",
			want:    true,
		},
		{name: "explicit disable", mode: "never", term: "xterm-256color"},
		{
			name:    "NO_COLOR disables automatic color",
			mode:    "auto",
			noColor: "1",
			term:    "xterm-256color",
		},
		{name: "dumb terminal disables automatic color", mode: "auto", term: "dumb"},
		{name: "buffer is not a terminal", mode: "auto", term: "xterm-256color"},
		{
			name: "pipe file descriptor is not a terminal",
			mode: "auto",
			term: "xterm-256color",
			pipe: true,
		},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.T().Setenv("NO_COLOR", tc.noColor)
			s.T().Setenv("TERM", tc.term)
			var out io.Writer = &bytes.Buffer{}
			if tc.pipe {
				reader, writer, err := os.Pipe()
				s.Require().NoError(err)
				s.T().Cleanup(func() {
					_ = reader.Close()
					_ = writer.Close()
				})
				out = writer
			}
			s.Equal(tc.want, tc.mode.Enabled(out))
		})
	}
}

func (s *ThemePublicTestSuite) TestNewTheme() {
	tests := []struct {
		name     string
		mode     cli.ColorMode
		noColor  string
		wantANSI bool
	}{
		{name: "automatic color on a buffer is plain", mode: "auto"},
		{name: "NO_COLOR keeps output plain", mode: "auto", noColor: "1"},
		{name: "explicit disable is plain", mode: "never"},
		{name: "explicit color overrides NO_COLOR", mode: "always", noColor: "1", wantANSI: true},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.T().Setenv("NO_COLOR", tc.noColor)
			s.T().Setenv("TERM", "xterm-256color")
			theme := cli.NewTheme(io.Discard, tc.mode)
			s.Equal(lipgloss.Color("#bc8cff"), theme.Accent.GetForeground())
			s.Equal(theme.Accent.GetForeground(), theme.BannerBot.GetForeground())
			got := theme.Accent.Render("example")
			if tc.wantANSI {
				s.Contains(got, "\x1b[")
				s.Contains(got, "example")
			} else {
				s.Equal("example", got)
			}
		})
	}
}

func (s *ThemePublicTestSuite) TestTitle() {
	escapes := regexp.MustCompile("\x1b\\[[0-9;]*m")
	tests := []struct {
		name     string
		mode     cli.ColorMode
		value    string
		wantBold bool
	}{
		{name: "plain title", mode: "never", value: "RESOURCE STATUS"},
		{name: "colored title is bold", mode: "always", value: "RESOURCE STATUS", wantBold: true},
		{name: "empty colored title has no visible text", mode: "always", wantBold: true},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := cli.NewTheme(io.Discard, tc.mode).Title(tc.value)
			s.Equal(tc.value, escapes.ReplaceAllString(got, ""))
			if tc.wantBold {
				s.Contains(got, "\x1b[1;")
				s.Contains(got, tc.value)
			} else {
				s.Equal(tc.value, got)
			}
		})
	}
}

func (s *ThemePublicTestSuite) TestBanner() {
	const want = "█▀█ █▀ █ █▀█ █▄░█\n█▀▀ ▄█ █ █▄█ █░▀█"
	escapes := regexp.MustCompile("\x1b\\[[0-9;]*m")
	tests := []struct {
		name     string
		mode     cli.ColorMode
		wantANSI bool
	}{
		{name: "plain wordmark", mode: "never"},
		{name: "colored wordmark keeps its text", mode: "always", wantANSI: true},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := cli.NewTheme(io.Discard, tc.mode).Banner()
			s.Equal(want, escapes.ReplaceAllString(got, ""))
			if tc.wantANSI {
				s.Contains(got, "\x1b[")
			} else {
				s.Equal(want, got)
			}
		})
	}
}

func TestThemePublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(ThemePublicTestSuite))
}

func (s *ThemePublicTestSuite) TestPhase() {
	theme := cli.NewTheme(io.Discard, "always")
	tests := []struct {
		name  string
		phase string
		want  lipgloss.Style
	}{
		{name: "success", phase: "Succeeded", want: theme.OK},
		{name: "failure", phase: "Failed", want: theme.Err},
		{name: "pending", phase: "Pending", want: theme.Info},
		{name: "unknown condition", phase: "Unknown", want: theme.Info},
		{name: "unrecognized condition stays visible", phase: "Custom", want: theme.Mute},
		{name: "empty condition stays empty", phase: "", want: theme.Mute},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.Equal(tc.want.Render(tc.phase), theme.Phase(tc.phase))
			plain := cli.NewTheme(io.Discard, "never")
			s.Equal(tc.phase, plain.Phase(tc.phase))
		})
	}
}
