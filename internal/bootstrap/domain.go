package bootstrap

import (
	"devteambot/internal/adapter/repository/gorm"
	"devteambot/internal/application/service"
)

func RegisterDomain() {
	RegisterSetting()
	RegisterReview()
	RegisterPoint()
	RegisterSholat()
}

func RegisterSetting() {
	appContainer.RegisterService("settingRepository", new(gorm.SettingRepository))
	appContainer.RegisterService("settingKey", gorm.NewSettingKey())
}

func RegisterReview() {
	appContainer.RegisterService("reviewRepository", new(gorm.ReviewRepository))
}

func RegisterPoint() {
	appContainer.RegisterService("pointRepository", new(gorm.PointRepository))
}

func RegisterSholat() {
	appContainer.RegisterService("sholatService", new(service.SholatService))
}
