package bootstrap

import (
	"devteambot/config"
	"devteambot/internal/adapter/cache"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/resty"
)

func RegisterDatabase(conf *config.Database) {
	db := BuildGorm(conf)

	appContainer.RegisterService("database", db)
}

func RegisterCache(conf *config.Redis) {
	redisClient := BuildRedis(conf)

	appContainer.RegisterService("cache", &cache.Redis{Client: redisClient})
}

func RegisterRest(c *config.Config) {
	// googleClient := resty.New().EnableTrace()
	// googleClient.SetRetryCount(3)
	// googleClient.SetHostURL(c.Sheets.Url)
	// appContainer.RegisterService("googleSheet", googleClient)

	// discord := resty.New().EnableTrace()
	// discord.SetRetryCount(3)
	// discord.SetHostURL("https://discord.com/api/v10")
	// appContainer.RegisterService("discordAPI", discord)
}

func RegisterDiscord(conf *config.Discord) {
	bot := BuildDiscord(conf)

	appContainer.RegisterService("botSession", bot)
	appContainer.RegisterService("discordApp", new(discord.App))
}

func RegisterAPI(conf *config.Config) {
	appContainer.RegisterService("myQuranAPI", new(resty.MyQuran))
}
