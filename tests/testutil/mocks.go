package testutil

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// MockImageStore implements parkingsession.ImageStore for tests.
type MockImageStore struct{}

func (m *MockImageStore) SaveImage(ctx context.Context, imageData []byte) (string, error) {
	return fmt.Sprintf("/tmp/test_%s.jpg", uuid.New().String()), nil
}

func (m *MockImageStore) DeleteImage(ctx context.Context, path string) error {
	return nil
}

// MockOCRService implements parkingsession.OCRService for tests.
type MockOCRService struct{}

func (m *MockOCRService) ExtractPlate(ctx context.Context, imageData []byte) (string, error) {
	return "MOCK-PLATE", nil
}
