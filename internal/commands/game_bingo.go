package commands

import (
	"context"
	"devteambot/internal/domain/bingo"
	"devteambot/internal/pkg/logger"
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) Bingo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	ctx := context.Background()

	keyStart := c.RedisKey.BingoStarted(i.GuildID)
	isPlaying, err := c.Cache.Exists(ctx, keyStart)
	if err != nil {
		response = "Something went wrong, please try again later"
		c.SendStandardResponse(i.Interaction, response, true, false)
		logger.Error(err.Error(), err)
		return
	}

	if isPlaying {
		response = "Another game is in progress"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	reward := int64(0)
	min := int64(0)
	max := int64(100)

	if opt, ok := optionMap["reward"]; ok {
		reward = opt.IntValue()
	}
	if opt, ok := optionMap["min"]; ok {
		min = opt.IntValue()
	}
	if opt, ok := optionMap["max"]; ok {
		max = opt.IntValue()
	}

	b := bingo.Bingo{
		ChannelID:    i.ChannelID,
		RandomNumber: rand.Int63n(max-min) + min,
		Reward:       reward,
	}

	duration := time.Duration(5)

	err = c.Cache.Put(ctx, keyStart, b, 20*time.Minute)
	if err != nil {
		response = "Failed to start the game"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	go func(key string, duration time.Duration, randomNumber int64) {
		logger.Info(fmt.Sprintf("Waiting %d minutes", duration))
		time.Sleep(duration * time.Minute)

		isPlaying, _ := c.Cache.Exists(ctx, key)
		if !isPlaying {
			return
		}

		s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Bingo has ended. The answer is %d", randomNumber))

		c.Cache.Delete(ctx, key)
	}(keyStart, duration, b.RandomNumber)

	logger.Info("Game started")
	response = fmt.Sprintf("🚨 **Bingo Start** 🚨\nGuess the number from %d to %d. Bingo end in %d minutes.\nReward: %d", min, max, duration, reward)
	s.ChannelMessageSend(i.ChannelID, response)
	logger.Info(fmt.Sprintf("Bingo: %d", b.RandomNumber))

	response = "Success to start Bingo"
	c.SendStandardResponse(i.Interaction, response, true, false)
}
