/*
Copyright © 2024 steffakasid
*/
package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/steffakasid/eslog"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{ .Values.ProjectName }",
	Short: "{ .Values.ShortDescription }",
	Long:  `{ .Values.LongDescription }`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Implement me!")
	},
}

func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	err := eslog.Logger.SetLogLevel(defaultLogLevel)
	eslog.LogIfErrorf(err, eslog.Warnf, "Error on setting log level: %v")

	cobra.OnInitialize(initConfig)

	peristentFlags := rootCmd.PersistentFlags()
	peristentFlags.StringP(debugFlag,
		debugFlagShorthand,
		defaultLogLevel,
		"Set log level [\"debug\", \"info\", \"warn\", \"error\"]")

	peristentFlags.StringP(configFlag,
		configFlagShorthand,
		"",
		"Set custom log file with path")

	err = viper.BindPFlags(peristentFlags)
	eslog.LogIfErrorf(err, eslog.Warnf, "[root] error binding persistent flags: %v")
}
