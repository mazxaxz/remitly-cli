package deploy

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func loadContexts() {
	viper.SetConfigName("contexts")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.remitly/")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config file not found: " + viper.ConfigFileUsed())
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	// log
	fmt.Println("Config loaded, path: " + viper.ConfigFileUsed())

	viper.AutomaticEnv()
	viper.SetEnvPrefix("REMITLY")

	viper.SetDefault("CONTEXT", "default")
	if err := viper.BindEnv("CONTEXT"); err != nil {
		fmt.Println(err)
		// log
	}
}
