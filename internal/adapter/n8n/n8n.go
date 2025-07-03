package n8n

import (
	"context"
	"devteambot/config"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type N8NAdapter interface {
	GenerateResponse(ctx context.Context, userId, guildId, prompt string) (*N8NResponse, error)
}

type N8N struct {
	Client *resty.Client
	Conf   *config.Config `inject:"config"`
}

func (n *N8N) Startup() error {
	client := resty.New()

	client.SetBaseURL(n.Conf.N8N.BaseURL)
	client.SetBasicAuth(n.Conf.N8N.Username, n.Conf.N8N.Password)

	n.Client = client

	return nil
}

func (n *N8N) Shutdown() error {
	return nil
}

func (n *N8N) GenerateResponse(ctx context.Context, userId, guildId, prompt string) (*N8NResponse, error) {
	request := n.Client.R().
		SetBody(map[string]any{
			"userId":  userId,
			"guildId": guildId,
			"input":   prompt,
		})

	response, err := request.Post(fmt.Sprintf("/webhook/%s", n.Conf.N8N.WebhookID))
	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("failed to generate response: %s", response.String())
	}

	var result N8NResponse
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
