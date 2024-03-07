package google

import (
	"context"
	"devteambot/config"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AI struct {
	Client *genai.Client
	Model  *genai.GenerativeModel
	Conf   *config.Config `inject:"config"`
}

func (a *AI) Startup() error {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(a.Conf.GoogleAI.Token))
	if err != nil {
		return err
	}
	a.Client = client
	a.Model = client.GenerativeModel("gemini-pro")

	return nil
}

func (a *AI) Shutdown() error {
	return a.Client.Close()
}
