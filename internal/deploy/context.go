package deploy

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func loadContexts(_ *cobra.Command, _ []string) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("REMITLY")
	viper.AllowEmptyEnv(true)

	if err := viper.BindEnv("PATH"); err != nil {
		log.WithError(err).Error("could not bind REMITLY_PATH environment variable, make sure it is set")
		return err
	}
	if err := viper.BindEnv("CONTEXT"); err != nil {
		log.WithError(err).Error("could not bind REMITLY_CONTEXT environment variable, make sure it is set")
		return err
	}

	path := viper.GetString("PATH")
	if path == "" {
		log.Info("make sure REMITLY_PATH is set")
		return ErrPathVariableNotFound
	}

	var fileName string
	err := filepath.Walk(path, func(path string, f fs.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "an error occured while scanning directory: '%s'", path)
		}
		if strings.HasSuffix(f.Name(), ".yml") {
			fileName = strings.ReplaceAll(f.Name(), ".yml", "")
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("make sure $REMITLY_PATH/*.yml files exist")
			return err
		}
		return err
	}
	log.Infof("config successfully loaded, using file: '%s'", viper.ConfigFileUsed())
	return nil
}
