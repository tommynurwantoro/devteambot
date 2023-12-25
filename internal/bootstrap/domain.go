package bootstrap

import (
	"devteambot/internal/adapter/repository/gorm"
	"devteambot/internal/application/service"
)

func RegisterDomain() {
	RegisterMessage()
	RegisterSetting()
	RegisterReview()
	RegisterPresensi()
	RegisterPoint()
	RegisterSholat()
}

func RegisterMessage() {
	appContainer.RegisterService("messageService", new(service.MessageService))
}

func RegisterSetting() {
	appContainer.RegisterService("settingRepository", new(gorm.SettingRepository))
	appContainer.RegisterService("settingService", new(service.SettingService))
}

func RegisterReview() {
	appContainer.RegisterService("reviewRepository", new(gorm.ReviewRepository))
	appContainer.RegisterService("reviewService", new(service.ReviewService))
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
