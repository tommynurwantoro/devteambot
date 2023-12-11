package repository

import (
	"devteambot/config"
	"devteambot/internal/pkg/logger"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Gorm struct {
	*gorm.DB
	Conf *config.Config `inject:"config"`
}

func (g *Gorm) Startup() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		g.Conf.Database.Host, g.Conf.Database.User, g.Conf.Database.Password, g.Conf.Database.DBName, g.Conf.Database.Port, g.Conf.Database.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(g.Conf.Database.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(g.Conf.Database.ConnMaxLifetime) * time.Hour)
	sqlDB.SetMaxOpenConns(g.Conf.Database.MaxOpenConn)

	logger.Info("Gorm initalized")
	g.DB = db

	return nil
}

func (g *Gorm) Shutdown() error {
	sqlDB, _ := g.DB.DB()
	return sqlDB.Close()
}
