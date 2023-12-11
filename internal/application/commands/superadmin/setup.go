package superadmin

import (
	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) Setup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c.Command.Discord.Bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modal_setup",
			Title:    "Insert Admin Role ID",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "setup_superadmin",
							Label:       "Please insert super admin role id",
							Style:       discordgo.TextInputShort,
							Placeholder: "982659357953118208,757119647102271588",
							Required:    true,
							MaxLength:   300,
							MinLength:   10,
						},
					},
				},
			},
		},
	})
}
