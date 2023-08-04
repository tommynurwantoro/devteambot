package scheduler

import (
	"context"
	"devteambot/internal/adapter/resty"
	"devteambot/internal/pkg/logger"
	"fmt"
	"time"
)

func (s *Scheduler) GetSholatSchedule(ctx context.Context) {
	response := new(resty.GetSholatResponse)
	req := s.MyQuranAPI.Client.R().SetContext(ctx).
		ForceContentType("application/json").
		SetResult(response)

	_, err := req.Get("/sholat/jadwal/1505/2023/8/4")
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %s", err.Error()), err)
		return
	}

	s.Cache.Put(ctx, s.RedisKey.DailySholatSchedule(), response, 24*time.Hour)
}
