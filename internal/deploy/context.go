package deploy

import (
	"os"

	"github.com/spf13/viper"
)

func loadContexts() {
	// TODO: from env
	viper.SetConfigName("$HOME/.remitly")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("contexts")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// fmt info
		} else {
			//fmt.Println(err)
		}
		os.Exit(1)
	}
	//fmt.Println("Config loaded, path: " + viper.ConfigFileUsed())

	viper.AutomaticEnv()
	viper.SetEnvPrefix("REMITLY")

	if err := viper.BindEnv("CONTEXT"); err != nil {
		//fmt.Println(err)
	}
}
