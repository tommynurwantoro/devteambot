package config

import (
	"time"
)

type Config struct {
	AppName       string        `valid:"required"`
	AppVersion    string        `valid:"required"`
	Environment   string        `valid:"required"`
	ShutdownDelay time.Duration `valid:"required"`
	Logger        Logger        `valid:"required"`
	Database      Database      `valid:"required"`
	Redis         Redis         `valid:"required"`
	Discord       Discord       `valid:"required"`
}

type Logger struct {
	Stdout       bool
	FileLocation string `valid:"required"`
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
	Password string
}
