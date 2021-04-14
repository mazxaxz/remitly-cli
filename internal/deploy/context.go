package deploy

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

const (
	contextsPath     = "$HOME/.remitly"
	contextsFileName = "contexts"

	contextDefaultProfile = "default"
)

const contextDefaultYml = `
contexts:
  - name: default
    http:
      url: https://XXXX/
      username: XXXX
`

func loadContexts() {
	viper.SetConfigName(contextsFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(contextsPath)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			if err = viper.ReadConfig(bytes.NewBuffer([]byte(contextDefaultYml))); err != nil {
				fmt.Println(err)
				// log
			}
			if err = touchContexts(); err != nil {
				fmt.Println(err)
				// log
			} else {
				if err = viper.WriteConfig(); err != nil {
					fmt.Println(err)
					// log
				}
			}
		} else {
			//fmt.Println(err)
			os.Exit(1)
		}
	}
	//fmt.Println("Config loaded, path: " + viper.ConfigFileUsed())

	viper.AutomaticEnv()
	viper.SetEnvPrefix("REMITLY")

	viper.SetDefault("PROFILE", contextDefaultProfile)
	if err := viper.BindEnv("CONTEXT"); err != nil {
		//fmt.Println(err)
	}
}

func touchContexts() error {
	fullContextsPath := strings.Replace(contextsPath, "$HOME", os.Getenv("HOME"), -1)
	fullFilePath := fmt.Sprintf("%s/%s.yml", fullContextsPath, contextsFileName)

	if _, err := os.Stat(fullContextsPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(fullContextsPath, os.FileMode(0755)); err != nil {
				return errors.Wrapf(err, "could not create directory: '%s'", fullContextsPath)
			}
			f, err := os.Create(fullFilePath)
			if err != nil {
				return errors.Wrapf(err, "could not create file: '%s'", fullFilePath)
			}
			defer f.Close()
			return nil
		} else {
			return errors.Wrapf(err, "an error occurred while checking if '%s' exists", fullContextsPath)
		}
	}

	if _, err := os.Stat(fullFilePath); err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "an error occurred while checking if '%s' exists", fullFilePath)
	}
	f, err := os.OpenFile(fullFilePath, os.O_CREATE, os.FileMode(0755))
	if err != nil {
		return errors.Wrapf(err, "could not create file: '%s'", fullFilePath)
	}
	defer f.Close()

	return nil
}
