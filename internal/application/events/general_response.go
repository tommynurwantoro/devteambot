package events

import (
	"context"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/n8n"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type GeneralResponseEvent struct {
	Discord        *discord.App           `inject:"discord"`
	N8N            n8n.N8NAdapter         `inject:"n8n"`
	MessageService service.MessageService `inject:"messageService"`
}

func (e *GeneralResponseEvent) Startup() error {
	e.Discord.Bot.AddHandler(e.HandleEvent)
	return nil
}

func (e *GeneralResponseEvent) Shutdown() error { return nil }

func (e *GeneralResponseEvent) HandleEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	go e.Do(s, m)
}

func (e *GeneralResponseEvent) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	mentions := m.Mentions
	mentionBullster := false
	if len(mentions) > 0 {
		for _, mention := range mentions {
			if mention.ID == s.State.User.ID {
				mentionBullster = true
			}
		}
	}

	if m.ReferencedMessage == nil && mentionBullster {
		chat, err := e.N8N.GenerateResponse(context.Background(), m.Author.ID, m.GuildID, m.Content)
		if err != nil {
			logger.Error("GeneralResponseEvent.Do", err)
			return
		}
		if _, err := s.ChannelMessageSendReply(m.ChannelID, chat.Output, m.Reference()); err != nil {
			logger.Error("GeneralResponseEvent.Do", err)
		}

		return
	}

	if m.ReferencedMessage != nil {
		refMessage := m.ReferencedMessage
		if refMessage.Author.ID == s.State.User.ID || mentionBullster {
			content := fmt.Sprintf("Referenced message from you (agent): '%s' --- User reply: '%s'", refMessage.Content, m.Content)
			chat, err := e.N8N.GenerateResponse(context.Background(), m.Author.ID, m.GuildID, content)
			if err != nil {
				logger.Error("GeneralResponseEvent.Do", err)
				return
			}
			if _, err := s.ChannelMessageSendReply(m.ChannelID, chat.Output, m.Reference()); err != nil {
				logger.Error("GeneralResponseEvent.Do", err)
			}
		}
	}
}
