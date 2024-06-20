package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/getsops/sops/v3/decrypt"
	"github.com/spf13/viper"
	"github.com/steffakasid/eslog"
)

// Constants used in command flags
const (
	debugFlag           = "debug"
	debugFlagShorthand  = "d"
	configFlag          = "config"
	configFlagShorthand = "c"

	defaultLogLevel = "debug"

	configFileType = "yaml"
	configFileName = ".{ .Values.ProjectName }"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	eslog.LogIfErrorf(err, eslog.Warnf, "Error on getting user home dir: %v")

	if cfgFile := viper.GetString(configFlag); cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".{ .Values.ProjectName }")
	}

	viper.AutomaticEnv()
	// Early log level after env vars are loaded to get log level from ENV
	err = eslog.Logger.SetLogLevel(viper.GetString(debugFlag))
	eslog.LogIfErrorf(err, eslog.Warnf, "Error on setting log level: %v")

	usedConfigFile := getConfigFilename(home)
	if usedConfigFile != "" {
		cleartext, err := decrypt.File(usedConfigFile, configFileType)

		if err != nil {
			eslog.Warnf("Error decrypting. %s. Maybe you're not using an encrypted config?", err)
			if err := viper.ReadInConfig(); err != nil {
				eslog.Warnf("Error reading config. %s. Are you using a config?", err)
			} else {
				eslog.Debug("Using config file:", viper.ConfigFileUsed())
			}
		} else {
			if err := viper.ReadConfig(bytes.NewBuffer(cleartext)); err != nil {
				eslog.Fatal(err)
			} else {
				eslog.Debug("Using sops encrypted config file:", viper.ConfigFileUsed())
			}
		}
	} else {
		eslog.Debug("No config file used!")
	}

	// Late log level to get log level from config file
	err = eslog.Logger.SetLogLevel(viper.GetString(debugFlag))
	eslog.LogIfErrorf(err, eslog.Warnf, "Error on setting log level: %v")

}

func getConfigFilename(homedir string) string {
	pathWithoutExt := path.Join(homedir, configFileName)
	eslog.Debugf("Check if %s exists", pathWithoutExt)
	if _, err := os.Stat(pathWithoutExt); err == nil {
		return pathWithoutExt
	}

	pathWithExt := fmt.Sprintf("%s.%s", pathWithoutExt, configFileType)
	eslog.Debugf("Check if %s exists", pathWithExt)
	if _, err := os.Stat(pathWithExt); err == nil {
		return pathWithExt
	}
	pathWithExt = fmt.Sprintf("%s.%s", pathWithoutExt, "yml")
	eslog.Debugf("Check if %s exists", pathWithExt)
	if _, err := os.Stat(pathWithExt); err == nil {
		return pathWithExt
	}
	return ""
}
