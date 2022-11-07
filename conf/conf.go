package conf

import "github.com/spf13/viper"

func ReadConf() {
	viper.SetConfigType("yaml")
}