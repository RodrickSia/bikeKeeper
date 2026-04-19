package OCR

import "context"

type PlateProcessor struct{}

func NewPlateProcessor() *PlateProcessor {
	return &PlateProcessor{}
}

func (p *PlateProcessor) ExtractPlate(ctx context.Context, imageData []byte) (string, error) {
	// Implement the logic to extract the plate number from the image data
	// TODO implement actual OCR request to external service
	return "MOCK-PLATE", nil
}