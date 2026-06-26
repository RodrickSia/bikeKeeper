package OCR

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type PlateProcessor struct {
	serviceURL string
	client     *http.Client
}

func NewPlateProcessor() (*PlateProcessor, error) {
	serviceURL := os.Getenv("OCR_SERVICE_URL")
	if serviceURL == "" {
		return nil, fmt.Errorf("OCR_SERVICE_URL environment variable is required")
	}
	return &PlateProcessor{
		serviceURL: serviceURL,
		client:     &http.Client{Timeout: 30 * time.Second},
	}, nil
}

type OCRServiceResponse struct {
	Plates []string `json:"plates"`
}

func (p *PlateProcessor) ExtractPlate(ctx context.Context, imageData []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.serviceURL, bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("creating OCR request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling OCR service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OCR service returned status %d: %s", resp.StatusCode, string(body))
	}

	var result OCRServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding OCR response: %w", err)
	}

	if len(result.Plates) > 0 {
		return result.Plates[0], nil
	}
	return "", nil
}