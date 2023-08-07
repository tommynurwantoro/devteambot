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
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	_, err := req.Get(fmt.Sprintf("/sholat/jadwal/1505/%d/%d/%d", now.Year(), int(now.Month()), now.Day()))
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %s", err.Error()), err)
		return
	}

	s.Cache.Put(ctx, s.RedisKey.DailySholatSchedule(), response, 24*time.Hour)
}
