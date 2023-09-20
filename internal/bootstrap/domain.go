package bootstrap

import "devteambot/internal/adapter/repository/sql"

func RegisterDomain() {
	RegisterSetting()
	RegisterReview()
	RegisterPoint()
}

func RegisterSetting() {
	appContainer.RegisterService("settingRepository", new(sql.SettingRepository))
}

func RegisterReview() {
	appContainer.RegisterService("reviewRepository", new(sql.ReviewRepository))
}

func RegisterPoint() {
	appContainer.RegisterService("pointRepository", new(sql.PointRepository))
}
