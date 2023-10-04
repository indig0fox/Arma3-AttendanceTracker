package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var ConfigJSON = viper.New()

func LoadConfig(modulePathDir string) (string, error) {
	ConfigJSON.SetConfigName("AttendanceTracker.config")
	ConfigJSON.SetConfigType("json")
	ConfigJSON.AddConfigPath(".")
	ConfigJSON.AddConfigPath(modulePathDir)

	ConfigJSON.SetDefault("armaConfig.dbUpdateInterval", "90s")
	ConfigJSON.SetDefault("armaConfig.debug", true)
	ConfigJSON.SetDefault("sqlConfig", map[string]interface{}{
		"mysqlHost":     "localhost",
		"mysqlPort":     3306,
		"mysqlUser":     "root",
		"mysqlPassword": "password",
		"mysqlDatabase": "a3attendance",
	})
	ConfigJSON.SetDefault("armaConfig", map[string]interface{}{
		"debug":            true,
		"traceLogToFile":   false,
		"dbUpdateInterval": "90s",
	})

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if err := ConfigJSON.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return "", fmt.Errorf(
				"config file not found, using defaults! searched in %s",
				[]string{
					ConfigJSON.ConfigFileUsed(),
					modulePathDir,
					wd,
				},
			)
		} else {
			// Config file was found but another error was produced
			return "", err
		}
	}

	return "Config loaded successfully!", nil
}

func ConfigArmaFormat() (string, error) {
	armaConfig := ConfigJSON.GetStringMap("armaConfig")
	bytes, err := json.Marshal(armaConfig)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
