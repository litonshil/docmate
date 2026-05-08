package llm

import (
	"bytes"
	"context"
	"docmate/types"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type GeminiProvider struct {
	ModelName string
}

func NewGeminiProvider(modelName string) *GeminiProvider {
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}

	return &GeminiProvider{ModelName: modelName}
}

func (p *GeminiProvider) GetName() string {
	return "gemini"
}

type geminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (p *GeminiProvider) GenerateSuggestions(ctx context.Context, apiKey string, complaints []string) (*types.AISuggestionResp, error) {
	prompt := fmt.Sprintf(`As a medical assistant, analyze these chief complaints: "%s".
Suggest potential diagnoses and medical investigations.
Return the result strictly as a JSON object with two arrays: "diagnoses" and "investigations".
Example: {"diagnoses": ["Viral Fever", "Common Cold"], "investigations": ["CBC", "Chest X-Ray"]}
Do not include any other text or formatting.`, strings.Join(complaints, ", "))

	reqBody := geminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/%s:generateContent?key=%s", p.ModelName, apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	slog.Info("sending request to gemini", "url", url)
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("gemini request failed", "error", err.Error())

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		slog.Error("gemini api error", "status", resp.StatusCode, "body", errResp)

		return nil, fmt.Errorf("gemini api returned status: %d", resp.StatusCode)
	}

	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from gemini")
	}

	rawResult := geminiResp.Candidates[0].Content.Parts[0].Text

	// Robustly extract JSON if it's wrapped in markdown code blocks
	if start := strings.Index(rawResult, "{"); start != -1 {
		if end := strings.LastIndex(rawResult, "}"); end != -1 && end > start {
			rawResult = rawResult[start : end+1]
		}
	}
	rawResult = strings.TrimSpace(rawResult)

	var result types.AISuggestionResp
	if err := json.Unmarshal([]byte(rawResult), &result); err != nil {
		return nil, fmt.Errorf("failed to parse suggestions: %w", err)
	}

	result.Disclaimer = "AI-generated suggestions. Please review and confirm."

	return &result, nil
}
