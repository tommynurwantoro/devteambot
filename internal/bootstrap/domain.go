package bootstrap

import "devteambot/internal/adapter/repository/sql"

func RegisterDomain() {
	RegisterSetting()
	RegisterReview()
}

func RegisterSetting() {
	appContainer.RegisterService("settingRepository", new(sql.SettingRepository))
}

func RegisterReview() {
	appContainer.RegisterService("reviewRepository", new(sql.ReviewRepository))
}
