package bootstrap

import (
	"devteambot/internal/adapter/repository/gorm"
)

func RegisterSetting() {
	appContainer.RegisterService("settingRepository", new(gorm.SettingRepository))
}

func RegisterReview() {
	appContainer.RegisterService("reviewRepository", new(gorm.ReviewRepository))
}

func RegisterPoint() {
	appContainer.RegisterService("pointRepository", new(gorm.PointRepository))
}

func RegisterMarketplace() {
	appContainer.RegisterService("marketplaceRepository", new(gorm.MarketplaceRepository))
}

func RegisterThanks() {
	appContainer.RegisterService("thanksRepository", new(gorm.ThanksRepository))
}
