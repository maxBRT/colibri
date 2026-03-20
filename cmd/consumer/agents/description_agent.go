package agents

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
	"www.github.com/maxbrt/colibri/internal/rss"
	"www.github.com/maxbrt/colibri/internal/utils"
)

const (
	modelID = "gemini-3.1-flash-lite-preview"
)

type DescriptionAgent struct {
	agent.Agent
}

func NewDescriptionAgent() (*DescriptionAgent, error) {
	ctx := context.Background()

	key, err := utils.GetSecret("google-api-key")
	if err != nil {
		return nil, err
	}

	model, err := gemini.NewModel(ctx, modelID, &genai.ClientConfig{
		APIKey: key,
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	prompt, err := os.ReadFile("/app/config/description_agent_prompt.md")
	if err != nil {
		log.Fatalf("Failed to read prompt: %v", err)
	}

	descAgent, err := llmagent.New(llmagent.Config{
		Name:        "description_agent",
		Model:       model,
		Description: "Generates a description for a given post.",
		Instruction: string(prompt),
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	return &DescriptionAgent{Agent: descAgent}, nil
}

func (a *DescriptionAgent) GenerateDescription(p *rss.Post) error {
	prompt := fmt.Sprintf("Generate a description for the post at this address: %s", p.Link)

	sessionService := session.InMemoryService()

	r, err := runner.New(runner.Config{
		AppName:        "colibri_consumer",
		Agent:          a.Agent,
		SessionService: sessionService,
	})
	if err != nil {
		return err
	}

	sess, err := sessionService.Create(context.Background(), &session.CreateRequest{
		AppName: "colibri_consumer",
		UserID:  "consummer",
	})
	if err != nil {
		return err
	}

	var result strings.Builder
	events := r.Run(
		context.Background(),
		"consummer",
		sess.Session.ID(),
		genai.NewContentFromText(prompt, genai.RoleUser),
		agent.RunConfig{StreamingMode: agent.StreamingModeNone},
	)

	for event, err := range events {
		if err != nil {
			return err
		}
		if event.IsFinalResponse() && event.Content != nil {
			for _, part := range event.Content.Parts {
				result.WriteString(part.Text)
			}
		}
	}

	p.Description = result.String()

	return nil
}
