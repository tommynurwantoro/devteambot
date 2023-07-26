package review

import (
	"devteambot/internal/domain/sharedkernel/entity"

	"github.com/lib/pq"
)

type Review struct {
	entity.Entity
	GuildID      string         `gorm:"not null"`
	Reporter     string         `gorm:"not null"`
	Title        string         `gorm:"not null"`
	Url          string         `gorm:"not null"`
	Reviewer     pq.StringArray `gorm:"type:text[]"`
	TotalPending int
}

func NewReview(guildID, reporter, title, url string, reviewer []string) *Review {
	return &Review{
		Entity:       entity.NewEntity(),
		GuildID:      guildID,
		Reporter:     reporter,
		Title:        title,
		Url:          url,
		Reviewer:     reviewer,
		TotalPending: len(reviewer),
	}
}

type Reviews []*Review
