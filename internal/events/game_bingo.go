package events

import (
	"context"
	"devteambot/internal/domain/bingo"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func (e *Event) Bingo(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.Background()
	isPlaying := false

	startedKey := fmt.Sprintf("bingo-%s-started", m.GuildID)
	isPlaying, _ = e.Cache.Exists(ctx, startedKey)
	if !isPlaying {
		return
	}

	var b bingo.Bingo
	err := e.Cache.Get(ctx, startedKey, &b)
	if err != nil {
		logger.Error(fmt.Sprintf("Error get random number %q", err.Error()), err)
		return
	}

	if m.ChannelID != b.ChannelID {
		return
	}

	stringNumber, err := strconv.Atoi(m.Content)
	if err != nil {
		return
	}

	if stringNumber == int(b.RandomNumber) {
		s.MessageReactionAdd(m.ChannelID, m.Message.ID, "✅")
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("🚨 **Bingo** 🚨\nCongratulations <@%s> guess the number %d", m.Author.ID, b.RandomNumber),
		)

		e.Cache.Delete(ctx, startedKey)
	}
}
