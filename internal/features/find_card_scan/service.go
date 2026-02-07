package find_card_scan

import (
	"context"
	"strings"
	"time"

	"neurofreyja/internal/entities/card"
	"neurofreyja/internal/shared/llm"
	"neurofreyja/internal/shared/storage"
)

type Service struct {
	Cards   card.Repository
	Storage *storage.Client
	LLM     *llm.Client
	Model   string
}

func (s *Service) FindCards(ctx context.Context, message string) ([]card.Card, error) {
	prompt := BuildPrompt(message)

	llmCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	response, err := s.LLM.ChatCompletion(llmCtx, s.Model, prompt)
	if err != nil {
		return nil, err
	}

	parsed, err := llm.ParseTitlesJSON(response)
	if err != nil {
		return nil, err
	}

	titles := normalizeTitles(parsed.Titles)
	if len(titles) == 0 {
		return nil, nil
	}

	return s.Cards.FindByTitles(ctx, titles)
}

func (s *Service) FetchImage(ctx context.Context, imageID string) ([]byte, error) {
	imgCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.Storage.GetObjectBytes(imgCtx, imageID)
}

func normalizeTitles(titles []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, t := range titles {
		normalized := strings.TrimSpace(strings.ToLower(t))
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	return out
}
