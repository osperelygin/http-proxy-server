package webapi

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	ProxyURL string
}

func getConfig(cfgPath string) (Config, error) {
	v := viper.GetViper()
	v.SetConfigFile(cfgPath)
	v.SetConfigType(strings.TrimPrefix(filepath.Ext(cfgPath), "."))

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	return Config{
		Port:     v.GetString("webapi.port"),
		ProxyURL: v.GetString("webapi.proxy_url"),
	}, nil
}
