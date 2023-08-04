package scheduler

import (
	"context"
	"devteambot/internal/pkg/logger"
	"fmt"
	"time"
)

type GetSholatResponse struct {
	Data struct {
		Jadwal struct {
			Tanggal string `json:"tanggal"`
			Dzuhur  string `json:"dzuhur"`
			Ashar   string `json:"ashar"`
		} `json:"jadwal"`
	} `json:"data"`
}

func (s *Scheduler) GetSholatSchedule(ctx context.Context) {
	response := new(GetSholatResponse)
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
