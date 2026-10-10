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
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/retr0h/psion/internal/cli"
)

// styleHelp adapts Cobra's command tree to the CLI renderer.
func styleHelp(
	root *cobra.Command,
) {
	root.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		out := cmd.OutOrStdout()
		// Cobra's help hook cannot return write errors.
		_ = commandHelp(cmd, cmd == root).Render(out, cli.NewTheme(out, color))
	})
	root.SetUsageFunc(func(cmd *cobra.Command) error {
		// UsageString captures output in a buffer; color follows the real destination.
		theme := cli.NewTheme(cmd.ErrOrStderr(), color)
		return commandHelp(cmd, false).Render(cmd.OutOrStdout(), theme)
	})
}

func commandHelp(
	cmd *cobra.Command,
	banner bool,
) cli.Help {
	help := cli.Help{
		Name:        cmd.CommandPath(),
		Description: cmd.Long,
		Usage:       cmd.UseLine(),
		Banner:      banner,
	}
	if help.Description == "" {
		help.Description = cmd.Short
	}
	if cmd.HasAvailableSubCommands() {
		help.Usage = cmd.CommandPath() + " <command> [flags]"
		help.Footer = "Run \"" + cmd.CommandPath() + " <command> --help\" for more."
	}
	for _, child := range cmd.Commands() {
		if child.IsAvailableCommand() {
			help.Commands = append(
				help.Commands,
				cli.Item{Name: child.Name(), Description: child.Short},
			)
		}
	}
	cmd.InitDefaultHelpFlag()
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
			help.Flags = append(help.Flags, cli.Item{Name: name, Description: description})
		})
	}
	sort.Slice(help.Flags, func(i, j int) bool { return help.Flags[i].Name < help.Flags[j].Name })
	return help
}
