package bootstrap

import (
	"devteambot/internal/adapter/repository/gorm"
	"devteambot/internal/application/service"
)

func RegisterDomain() {
	RegisterSetting()
	RegisterReview()
	RegisterPresensi()
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

func RegisterPresensi() {
	appContainer.RegisterService("presensiService", new(service.PresensiService))
}

func RegisterPoint() {
	appContainer.RegisterService("pointRepository", new(gorm.PointRepository))
	appContainer.RegisterService("pointService", new(service.PointService))
}

func RegisterSholat() {
	appContainer.RegisterService("sholatService", new(service.SholatService))
}
