package superadmin

import (
	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) Setup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c.Command.Discord.Bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "setup_superadmin",
			Title:    "Insert Admin Role ID",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "insert_role_id",
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
