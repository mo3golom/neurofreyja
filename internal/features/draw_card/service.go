package draw_card

import (
	"context"
	"strings"
	"time"

	"neurofreyja/internal/entities/card"
	"neurofreyja/internal/entities/drawn"
	"neurofreyja/internal/shared/llm"
	"neurofreyja/internal/shared/storage"
)

type Service struct {
	Cards   card.Repository
	Drawn   drawn.Repository
	Storage *storage.Client
	LLM     *llm.Client
	Model   string
}

func (s *Service) Draw(ctx context.Context, chatID int64) (*card.Card, string, error) {
	cardItem, err := s.Cards.FindRandomNotDrawn(ctx, chatID)
	if err != nil || cardItem == nil {
		return cardItem, "", err
	}

	description := ""
	if cardItem.Description.Valid {
		description = strings.TrimSpace(cardItem.Description.String)
	}
	if description == "" {
		prompt := BuildPrompt(cardItem.Title)
		llmCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		response, err := s.LLM.ChatCompletion(llmCtx, s.Model, prompt)
		if err != nil {
			return nil, "", err
		}

		description = strings.TrimSpace(response)
		if description != "" {
			_ = s.Cards.UpdateDescription(ctx, cardItem.ID, description)
		}
	}

	return cardItem, description, nil
}

func (s *Service) FetchImage(ctx context.Context, imageID string) ([]byte, error) {
	imgCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.Storage.GetObjectBytes(imgCtx, imageID)
}

func (s *Service) MarkDrawn(ctx context.Context, cardID int64, chatID int64) error {
	return s.Drawn.Insert(ctx, cardID, chatID)
}
