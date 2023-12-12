package bootstrap

import "devteambot/internal/adapter/repository/gorm"

func RegisterDomain() {
	RegisterSetting()
	RegisterReview()
	RegisterPoint()
}

func RegisterSetting() {
	appContainer.RegisterService("settingRepository", new(gorm.SettingRepository))
}

func RegisterReview() {
	appContainer.RegisterService("reviewRepository", new(gorm.ReviewRepository))
}

func RegisterPoint() {
	appContainer.RegisterService("pointRepository", new(gorm.PointRepository))
}
