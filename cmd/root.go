package cmd

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/retr0h/psion/internal/cli"
)

var (
	debug  bool
	logger *slog.Logger
	color  cli.ColorMode = "auto"
	//go:embed resources/*.yaml
	eFs       embed.FS
	appFs     afero.Fs
	stateFile string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "psion",
	Short: "Declare system state. Ship a binary. Apply it.",
	Long:  "Preview and apply the system state embedded in this binary.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	appFs = afero.NewOsFs()
	styleHelp(rootCmd)
	rootCmd.SilenceErrors = true
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		out := rootCmd.ErrOrStderr()
		style := cli.NewTheme(out, color)
		_, _ = fmt.Fprintf(out, "\n  %s %s\n\n", style.Err.Render("Error:"), err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger)

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "set log level to debug")
	rootCmd.PersistentFlags().Var(&color, "color", "terminal color: auto, always, or never")
	rootCmd.Flags().
		StringVarP(&stateFile, "state-file", "s", ".state", "path to the state file.")

	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		return
	}
}

func initLogger() {
	logLevel := slog.LevelInfo
	if viper.GetBool("debug") {
		logLevel = slog.LevelDebug
	}

	logger = slog.New(
		tint.NewTextHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.Kitchen,
			NoColor:    !color.Enabled(os.Stderr),
		}),
	)
}
