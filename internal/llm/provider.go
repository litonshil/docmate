package llm

import (
	"context"
	"docmate/types"
)

type Provider interface {
	GenerateSuggestions(ctx context.Context, apiKey string, complaints []string) (*types.AISuggestionResp, error)
	GetName() string
}
