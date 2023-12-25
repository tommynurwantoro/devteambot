package review

import (
	"context"
)

type Service interface {
	GetAntrian(ctx context.Context, guildID string) (Reviews, error)
	UpdateDone(ctx context.Context, r *Review, reviewerID string) error
	AddReviewer(ctx context.Context, guildID, userID, title, url string, reviewers []string) error
}
