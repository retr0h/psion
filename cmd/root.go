package cmd

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	debug   bool
	logger  *logrus.Logger
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
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "psion.yaml", "config file")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "set log level to debug")

	// bind flags to config
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		return
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("psion.yaml")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("psion")

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Fatal("no valid config found")
	}

	// config = geo.Config{}
	// if err := viper.Unmarshal(&config); err != nil {
	// 	logrus.WithError(err).Fatal("failed to unmarshal config")
	// }
}

func initLogger() {
	logLevel := logrus.InfoLevel
	if viper.GetBool("debug") {
		logLevel = logrus.TraceLevel
	}

	logger = &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			FullTimestamp:             true,
			DisableLevelTruncation:    true,
			PadLevelText:              true,
			EnvironmentOverrideColors: true,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logLevel,
	}
}
