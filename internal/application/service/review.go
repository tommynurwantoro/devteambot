package service

import (
	"context"
	"devteambot/internal/domain/review"
	"devteambot/internal/domain/sharedkernel/entity"
	"devteambot/internal/pkg/constant"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ReviewService interface {
	GetAntrian(ctx context.Context, guildID string) (review.Reviews, error)
	UpdateDone(ctx context.Context, r *review.Review, reviewerID string) error
	AddReviewer(ctx context.Context, guildID, userID, title, url string, reviewers []string) error
	PrettyAntrian(reviews review.Reviews) (string, *discordgo.MessageEmbed)
}

type Review struct {
	ReviewRepository review.Repository `inject:"reviewRepository"`
}

func (s *Review) Startup() error { return nil }

func (s *Review) Shutdown() error { return nil }

func (s *Review) GetAntrian(ctx context.Context, guildID string) (review.Reviews, error) {
	listReview, err := s.ReviewRepository.GetAllPendingByGuildID(ctx, guildID)
	if err != nil {
		return nil, err
	}

	return listReview, nil
}

func (s *Review) UpdateDone(ctx context.Context, r *review.Review, reviewerID string) error {
	newReviewer := make([]string, 0)
	found := false
	for _, s := range r.Reviewer {
		if s == reviewerID {
			found = true
			continue
		}

		newReviewer = append(newReviewer, s)
	}

	if !found {
		return review.ErrReviewerNotFound
	}

	r.Reviewer = newReviewer
	r.TotalPending = len(newReviewer)

	if err := s.ReviewRepository.Update(ctx, r); err != nil {
		return err
	}

	return nil
}

func (s *Review) AddReviewer(ctx context.Context, guildID, userID, title, url string, reviewers []string) error {
	r := &review.Review{
		Entity:       entity.NewEntity(),
		GuildID:      guildID,
		Reporter:     userID,
		Title:        title,
		Url:          url,
		Reviewer:     reviewers,
		TotalPending: len(reviewers),
	}

	if err := s.ReviewRepository.Update(ctx, r); err != nil {
		return err
	}

	return nil
}

func (s *Review) PrettyAntrian(reviews review.Reviews) (string, *discordgo.MessageEmbed) {
	if len(reviews) == 0 {
		return "", &discordgo.MessageEmbed{
			Title:       "Antrian Review",
			Description: "Gak ada antrian review nih 👍🏻",
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Antrian Review",
		Description: "Reviewer dapat menggunakan command `/sudah_direview` untuk menyelesaikan review. Berikut adalah antrian review tim kamu:",
		Color:       constant.GREEN,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	content := ""

	for i, r := range reviews {
		reviewer := ""
		for _, user := range r.Reviewer {
			if !strings.Contains(content, user) {
				content = fmt.Sprintf("%s <@%s>", content, user)
			}
			reviewer = fmt.Sprintf("%s, <@%s>", reviewer, user)
		}

		reviewer = reviewer[2:]

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%d. %s", i+1, r.Title),
			Value: fmt.Sprintf("%s \nReviewer: %s", r.Url, reviewer),
		})
	}

	return content, embed
}
