/*
Copyright © 2024 steffakasid
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mozilla.org/sops/v3/decrypt"
)

const (
	configFileType = "yaml"
	configFileName = ".clinar"
)

// Constants used in command flags
const (
	dryrun    = "dry-run"
	olderthen = "older-then"
	debugFlag = "debug"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{ .Values.ProjectName }",
	Short: "{ .Values.ShortDescription }",
	Long:  `{ .Values.LongDescription }`,
}

func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	logger.SetLevel(logger.DebugLevel)
	cobra.OnInitialize(initConfig)

	peristentFlags := rootCmd.PersistentFlags()

	cobra.CheckErr(viper.BindPFlags(peristentFlags))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".{ .Values.ProjectName }")
	}

	viper.AutomaticEnv()
	// Early log level after env vars are loaded to get log level from ENV
	setLogLevel()

	usedConfigFile := getConfigFilename(home)
	if usedConfigFile != "" {
		cleartext, err := decrypt.File(usedConfigFile, configFileType)

		if err != nil {
			logger.Warnf("Error decrypting. %s. Maybe you're not using an encrypted config?", err)
			if err := viper.ReadInConfig(); err != nil {
				logger.Warnf("Error reading config. %s. Are you using a config?", err)
			} else {
				logger.Debug("Using config file:", viper.ConfigFileUsed())
			}
		} else {
			if err := viper.ReadConfig(bytes.NewBuffer(cleartext)); err != nil {
				logger.Fatal(err)
			} else {
				logger.Debug("Using sops encrypted config file:", viper.ConfigFileUsed())
			}
		}
	} else {
		logger.Debug("No config file used!")
	}

	// Late log level to get log level from config file
	setLogLevel()
}

func getConfigFilename(homedir string) string {
	pathWithoutExt := path.Join(homedir, configFileName)
	logger.Debugf("Check if %s exists", pathWithoutExt)
	if _, err := os.Stat(pathWithoutExt); err == nil {
		return pathWithoutExt
	}

	pathWithExt := fmt.Sprintf("%s.%s", pathWithoutExt, configFileType)
	logger.Debugf("Check if %s exists", pathWithExt)
	if _, err := os.Stat(pathWithExt); err == nil {
		return pathWithExt
	}
	pathWithExt = fmt.Sprintf("%s.%s", pathWithoutExt, "yml")
	logger.Debugf("Check if %s exists", pathWithExt)
	if _, err := os.Stat(pathWithExt); err == nil {
		return pathWithExt
	}
	return ""
}

func setLogLevel() {
	var level logger.Level

	switch strings.ToLower(viper.GetString(debugFlag)) {
	case "debug":
		level = logger.DebugLevel
	case "info":
		level = logger.InfoLevel
	case "warn":
		level = logger.WarnLevel
	case "error":
		level = logger.ErrorLevel
	case "fatal":
		level = logger.FatalLevel
	default:
		level = logger.InfoLevel
	}
	logger.SetLevel(level)
}
