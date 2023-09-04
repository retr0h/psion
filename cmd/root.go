package cmd

import (
	"embed"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	debug  bool
	logger *slog.Logger
	//go:embed resources/*.yaml
	eFs       embed.FS
	appFs     afero.Fs
	stateFile string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "psion",
	Short: "A simplistic Go based system automation tool.",
	Long: `A simple command-line tool to carry out system automation tasks.

      '
|~~\(~|/~\|/~\
|__/_)|\_/|   |
|
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	appFs = afero.NewOsFs()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger)

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "set log level to debug")
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
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.Kitchen,
		}),
	)
}
