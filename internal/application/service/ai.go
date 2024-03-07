package service

import (
	"context"
	"devteambot/internal/adapter/google"
	"devteambot/internal/pkg/logger"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

type AIService interface {
	// GetResponse return response from AI
	GetResponse(ctx context.Context, prompt string) (string, error)
}

type AI struct {
	GoogleAI *google.AI `inject:"googleai"`
}

func (a *AI) Startup() error { return nil }

func (a *AI) Shutdown() error { return nil }

func (a *AI) GetResponse(ctx context.Context, prompt string) (string, error) {
	resp, err := a.GoogleAI.Model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		logger.Error("Error to get response from AI", err)
		return "", err
	}

	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			result := fmt.Sprintf("%s", part)
			if result != "" {
				return result, nil
			}
		}
	}

	return "Ask another question", nil
}
