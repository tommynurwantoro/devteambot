package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName       string               `valid:"required"`
	AppVersion    string               `valid:"required"`
	Environment   string               `valid:"required"`
	ShutdownDelay time.Duration        `valid:"required"`
	Http          HttpConfig           `valid:"required"`
	Logger        Logger               `valid:"required"`
	Database      Database             `valid:"required"`
	Redis         Redis                `valid:"required"`
	Discord       Discord              `valid:"required"`
	Schedulers    map[string]Scheduler `valid:"required"`
	GoogleAI      GoogleAIConfig       `valid:"required"`
}

type HttpConfig struct {
	Port         int `valid:"required"`
	WriteTimeout int `valid:"required"`
	ReadTimeout  int `valid:"required"`
}

type Logger struct {
	Stdout        bool
	FileLocation  string `valid:"required"`
	FileMaxSize   int    `valid:"required"`
	FileMaxBackup int    `valid:"required"`
	FileMaxAge    int    `valid:"required"`
}

type Discord struct {
	AppID            string `valid:"required"`
	Token            string `valid:"required"`
	RunInitCommand   bool
	RunDeleteCommand bool
}

type Role struct {
	ID   string `valid:"required"`
	Name string `valid:"required"`
}

type Database struct {
	Host            string `valid:"required"`
	User            string `valid:"required"`
	Password        string `valid:"required"`
	DBName          string `valid:"required"`
	Port            string `valid:"required"`
	SSLMode         string `valid:"required"`
	MaxIdleConn     int    `valid:"required"`
	ConnMaxLifetime int    `valid:"required"`
	MaxOpenConn     int    `valid:"required"`
}

type Redis struct {
	Address  string `valid:"required"`
	Port     int    `valid:"required"`
	Password string
}

type GoogleAIConfig struct {
	Token string `valid:"required"`
}

type Scheduler struct {
	Enable bool
	Time   struct {
		Hour   uint
		Minute uint
		Second uint
	}
}

func (c *Config) Load(path string) {
	viper.AddConfigPath(".")
	viper.SetConfigName(path)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}

	err = viper.Unmarshal(c)
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
}

func getEnvOrPanic(env string) string {
	split := strings.Split(env, ":")
	res := os.Getenv(split[0])
	if len(res) == 0 {
		if len(split) > 1 {
			res = strings.Join(split[1:], ":")
		}
		if len(res) == 0 {
			panic("Mandatory env variable not found:" + env)
		}
	}
	return res
}
