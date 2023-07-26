package bootstrap

import (
	"devteambot/config"
	"devteambot/internal/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func BuildGorm(conf *config.Database) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port, conf.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Panic("Cannot initiate gorm ", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Hour)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)

	logger.Info("Gorm initalized")

	return db
}

func BuildRedis(conf *config.Redis) *redis.Client {
	option := &redis.Options{
		Addr:     fmt.Sprintf("%s", conf.Address),
		Password: conf.Password,
	}

	return redis.NewClient(option)
}

func BuildDiscord(conf *config.Discord) *discordgo.Session {
	bot, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		logger.Panic("Cannot authorizing bot ", err)
	}
	logger.Info("Bot initialized")

	return bot
}

// func BuildTwitter(conf *config.Twitter) *resty.Client {
// 	client := resty.New()
// 	client.BaseURL = "https://api.twitter.com"
// 	client.Token = conf.Token

// 	logger.Info("Twitter initialized")

// 	return client
// }
