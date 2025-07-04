package bootstrap

import (
	"devteambot/config"
	"devteambot/internal/application/api"
	"devteambot/internal/application/commands/member"
	"devteambot/internal/application/commands/superadmin"
	"devteambot/internal/application/events"
	"devteambot/internal/application/scheduler"
	"devteambot/internal/application/service"
)

func RegisterScheduler(conf *config.Config) {
	appContainer.RegisterService("scheduler", new(scheduler.Scheduler))

	appContainer.RegisterService("presensiScheduler", new(scheduler.PresensiScheduler))
	appContainer.RegisterService("poinScheduler", new(scheduler.PoinScheduler))
	appContainer.RegisterService("sholatScheduler", new(scheduler.SholatScheduler))
}

func RegisterService() {
	appContainer.RegisterService("messageService", new(service.Message))
	appContainer.RegisterService("settingService", new(service.Setting))
	appContainer.RegisterService("reviewService", new(service.Review))
	appContainer.RegisterService("presensiService", new(service.Presensi))
	appContainer.RegisterService("pointService", new(service.Point))
	appContainer.RegisterService("sholatService", new(service.Sholat))
	appContainer.RegisterService("aiService", new(service.AI))
}

func RegisterAPI() {
	appContainer.RegisterService("api", new(api.Api))
}

func RegisterCommand() {
	// Superadmin
	appContainer.RegisterService("panelSettingCommand", new(superadmin.PanelSettingCommand))
	appContainer.RegisterService("activatePointCommand", new(superadmin.ActivatePointCommand))
	appContainer.RegisterService("activateReminderPresensiCommand", new(superadmin.ActivateReminderPresensiCommand))
	appContainer.RegisterService("ActivateReminderSholatCommand", new(superadmin.ActivateReminderSholatCommand))
	appContainer.RegisterService("addButtonFeatureCommand", new(superadmin.AddButtonFeatureCommand))
	appContainer.RegisterService("deleteButtonFeatureCommand", new(superadmin.DeleteButtonFeatureCommand))
	appContainer.RegisterService("editEmbedCommand", new(superadmin.EditEmbedCommand))
	appContainer.RegisterService("sendEmbedCommand", new(superadmin.SendEmbedCommand))

	// Member
	appContainer.RegisterService("antrianReviewCommand", new(member.AntrianReviewCommand))
	appContainer.RegisterService("askCommand", new(member.AskCommand))
	appContainer.RegisterService("claimRoleCommand", new(member.ClaimRoleCommand))
	appContainer.RegisterService("pingCpmmand", new(member.PingCommand))
	appContainer.RegisterService("sudahDireviewCommand", new(member.SudahDireviewCommand))
	appContainer.RegisterService("thanksCommand", new(member.ThanksCommand))
	appContainer.RegisterService("thanksLeaderboardCommand", new(member.ThanksLeaderboardCommand))
	appContainer.RegisterService("titipReviewCommand", new(member.TitipReviewCommand))

	appContainer.RegisterService("commandSuperAdmin", new(superadmin.Command))
	appContainer.RegisterService("commandMember", new(member.Command))

	// events
	appContainer.RegisterService("generalResponseEvent", new(events.GeneralResponseEvent))
}
