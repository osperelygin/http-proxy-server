package proxy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	Script   string
	CertsDir string
	KeyFile  string
	CertFile string
}

func getConfig(cfgPath string) (Config, error) {
	v := viper.GetViper()
	v.SetConfigFile(cfgPath)
	v.SetConfigType(strings.TrimPrefix(filepath.Ext(cfgPath), "."))

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	currDir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	certsDirRelPath := v.GetString("proxy.certs_dir")

	return Config{
		Port:     v.GetString("proxy.port"),
		Script:   filepath.Join(currDir, v.GetString("proxy.certs_gen_script")),
		CertsDir: filepath.Join(currDir, certsDirRelPath),
		KeyFile:  filepath.Join(currDir, certsDirRelPath, v.GetString("proxy.key_file")),
		CertFile: filepath.Join(currDir, certsDirRelPath, v.GetString("proxy.cert_file")),
	}, nil
}
