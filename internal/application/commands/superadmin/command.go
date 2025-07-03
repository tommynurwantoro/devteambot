package superadmin

import (
	"devteambot/config"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/pkg/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Conf    *config.Config `inject:"config"`
	Discord *discord.App   `inject:"discord"`
	cmdList []*discordgo.ApplicationCommand
}

func (c *Command) Startup() error {
	if c.Conf.Discord.RunResetCommand {
		c.runAddCommand()
	}

	return nil
}

func (c *Command) Shutdown() error { return nil }

func (c *Command) AppendCommand(cmd *discordgo.ApplicationCommand) {
	serverManager := int64(discordgo.PermissionManageServer)
	dmPermission := new(bool)
	*dmPermission = false

	cmd.DefaultMemberPermissions = &serverManager
	cmd.DMPermission = dmPermission

	c.cmdList = append(c.cmdList, cmd)
}

func (c *Command) runAddCommand() {
	logger.Info("Adding superadmin commands...")
	for _, v := range c.cmdList {
		logger.Info(fmt.Sprintf("Adding command: %v", v.Name))
		cmd, err := c.Discord.Bot.ApplicationCommandCreate(c.Conf.Discord.AppID, "", v)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Cannot create command: %v %s", v.Name, err.Error()), err)
		}

		c.cmdList = append(c.cmdList, cmd)
	}
}
