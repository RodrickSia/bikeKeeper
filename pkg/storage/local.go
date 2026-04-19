package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type LocalStorage struct {
	baseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{baseDir: baseDir}
}

func (s *LocalStorage) SaveImage(ctx context.Context, imageData []byte) (string, error) {
	err := os.MkdirAll(s.baseDir, 0755)
	if err != nil {
		return "", fmt.Errorf("creating base directory: %w", err)
	} 
	
	imageName := fmt.Sprintf("image_%s.jpg", uuid.New().String())
	imagePath := filepath.Join(s.baseDir, imageName)
	err = os.WriteFile(imagePath, imageData, 0644)
	if err != nil {
		return "", fmt.Errorf("saving image: %w", err)
	}

	return imagePath, nil
}

func (s *LocalStorage) DeleteImage(ctx context.Context, imagePath string) error {
	absPath, _ := filepath.Abs(imagePath)
	absBase, _ := filepath.Abs(s.baseDir)
	// Ensure the path is within the base directory to prevent directory traversal attacks
	if !strings.HasPrefix(absPath, absBase+string(filepath.Separator)) {
		return fmt.Errorf("path is outside allowed directory")
	}
	err := os.Remove(absPath)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	return nil
}

