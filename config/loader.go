package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Load(path string) (config *Config, err error) {
	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
		return
	}()

	var (
		dir      = filepath.Dir(path)
		filename = filepath.Base(path)
		ext      = filepath.Ext(path)[1:]
	)

	viper.AddConfigPath(dir)
	viper.SetConfigName(filename)
	viper.SetConfigType(ext)

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimPrefix(strings.TrimSuffix(value, "}"), "${")))
		}
	}

	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
