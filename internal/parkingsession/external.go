package parkingsession

import "context"

// ParkingSessionResponse represents the response from the external API
type ParkingSessionResponse struct {
	Image  string         `json:"image"`
	Events []ParkingEvent `json:"events"`
}

// ParkingEvent represents a parking event
type ParkingEvent struct {
	Plate     string `json:"plate"`
	Event     string `json:"event"`
	Timestamp string `json:"timestamp"`
}

type ImageStore interface {
	SaveImage(ctx context.Context, imageData []byte) (string, error)
	DeleteImage(ctx context.Context, path string) error
}

type OCRService interface {
	ExtractPlate(ctx context.Context, imageData []byte) (string, error)
}


