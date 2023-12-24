package service

import (
	"context"
	"devteambot/internal/domain/review"
	"devteambot/internal/domain/sharedkernel/entity"
)

type ReviewService struct {
	ReviewRepository review.Repository `inject:"reviewRepository"`
}

func (s *ReviewService) Startup() error { return nil }

func (s *ReviewService) Shutdown() error { return nil }

func (s *ReviewService) GetAntrian(ctx context.Context, guildID string) (review.Reviews, error) {
	listReview, err := s.ReviewRepository.GetAllPendingByGuildID(ctx, guildID)
	if err != nil {
		return nil, err
	}

	return listReview, nil
}

func (s *ReviewService) UpdateDone(ctx context.Context, r *review.Review, reviewerID string) error {
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

func (s *ReviewService) AddReviewer(ctx context.Context, guildID, userID, title, url string, reviewers []string) error {
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
