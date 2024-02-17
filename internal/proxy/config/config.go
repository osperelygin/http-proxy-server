package config

import (
	"http-proxy-server/internal/pkg/config"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func GetConfig(cfgPath string) config.SrvConfig {
	v := viper.GetViper()
	v.SetConfigFile(cfgPath)
	v.SetConfigType(strings.TrimPrefix(filepath.Ext(cfgPath), "."))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	return config.SrvConfig{
		Port: v.GetString("webapi.port"),
		Host: v.GetString("webapi.host"),
	}
}
