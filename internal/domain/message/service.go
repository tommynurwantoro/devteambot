package message

import "github.com/bwmarrin/discordgo"

type Service interface {
	SendStandardResponse(i *discordgo.Interaction, response string, isPrivate, isRemovePreview bool)
	EditStandardResponse(i *discordgo.Interaction, response string)
	SendEmbedResponse(i *discordgo.Interaction, content string, embed *discordgo.MessageEmbed, isPrivate bool)
	SendStandardMessage(channelID, message string)
}
