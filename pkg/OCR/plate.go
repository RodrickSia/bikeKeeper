package OCR

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type PlateProcessor struct {
	serviceURL string
	client     *http.Client
}

func NewPlateProcessor() *PlateProcessor {
	return &PlateProcessor{
		serviceURL: os.Getenv("OCR_SERVICE_URL"),
		client:     &http.Client{},
	}
}

type OCRServiceResponse struct {
	PlateNumber string `json:"plateNumber"`
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
		return "", fmt.Errorf("OCR service returned status %d", resp.StatusCode)
	}

	var result OCRServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding OCR response: %w", err)
	}

	return result.PlateNumber, nil
}