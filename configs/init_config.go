package configs

import (
	"HongJungWan-spaceIPX-backend-api/pkg/helper"
	"HongJungWan-spaceIPX-backend-api/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

func InitConfig(configFile string) error {

	viper.AutomaticEnv()
	configFileDirectory := os.Getenv(CONFIG_DIRECTORY)
	if configFileDirectory != "" {
		viper.AddConfigPath(configFileDirectory)
	}
	configFilePath := os.Getenv(CONFIG_FILE)
	if configFilePath != "" {
		os.Setenv(CONFIG_FILE, strings.Trim(configFilePath, " "))
		configFilePath = os.Getenv(CONFIG_FILE)
		log.Printf("InitConfig read from ENV CONFIG_FILE=%s", configFilePath)
		if exist, _ := helper.Exists(configFilePath); exist {
			configFile = configFilePath
		}
	}
	if configFile != "" {
		if exist, _ := helper.Exists(configFile); exist {
			viper.SetConfigFile(configFile)
		}
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")

		configDirectory, err := FindConfigDirectoryPath()
		viper.AddConfigPath(configDirectory)

		if runtime.GOOS == "windows" {

		} else if runtime.GOOS == "linux" {
			viper.AddConfigPath("/etc")
		}

		configFile := GetConfigFileName()
		viper.SetConfigName(configFile)
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return nil
}

func InitLogger() logger.Logger {
	// logger default setting
	loggerConfig := logger.Config{
		Console: logger.Console{
			Enable:     true,
			JsonFormat: false,
			Level:      logger.INFO,
		},
	}

	// Load config from viper
	if viper.IsSet("logger") {
		data, err := json.Marshal(viper.GetStringMap("logger"))
		if err != nil {
			panic(fmt.Sprintf("error marshaling logger config: %v", err))
		}

		err = json.Unmarshal(data, &loggerConfig)
		if err != nil {
			panic(fmt.Sprintf("error unmarshaling logger config: %v", err))
		}
	}

	return logger.Init(loggerConfig)
}

func getEnvironment() string {
	env := strings.Trim(os.Getenv(ENV_MODE), " ")
	if env == "" {
		env = "dev"
	}
	return env
}

func FindConfigDirectoryPath() (string, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	configDirectory := path.Join(currentDirectory, "/configs")
	return configDirectory, nil
}

func GetConfigFileName() string {
	return fmt.Sprintf("..%s", getEnvironment())
}
