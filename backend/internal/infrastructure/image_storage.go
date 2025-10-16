package infrastructure

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// ImageStorage interface for storing images
type ImageStorage interface {
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
	Delete(ctx context.Context, filename string) error
	GetURL(filename string) string
}

// LocalImageStorage stores images locally (for development)
type LocalImageStorage struct {
	basePath string
	baseURL  string
}

func NewLocalImageStorage(basePath, baseURL string) ImageStorage {
	return &LocalImageStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

func (s *LocalImageStorage) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	// In production, save to disk or upload to GCS
	// For now, return a mock URL
	url := fmt.Sprintf("%s/uploads/%s", s.baseURL, uniqueFilename)

	return url, nil
}

func (s *LocalImageStorage) Delete(ctx context.Context, filename string) error {
	// In production, delete from disk or GCS
	return nil
}

func (s *LocalImageStorage) GetURL(filename string) string {
	return fmt.Sprintf("%s/uploads/%s", s.baseURL, filename)
}

// GCSImageStorage stores images in Google Cloud Storage
type GCSImageStorage struct {
	bucketName string
	cdnURL     string
}

func NewGCSImageStorage(bucketName, cdnURL string) ImageStorage {
	return &GCSImageStorage{
		bucketName: bucketName,
		cdnURL:     cdnURL,
	}
}

func (s *GCSImageStorage) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	// In production, upload to GCS
	// For now, return mock CDN URL
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	cdnURL := fmt.Sprintf("%s/%s", s.cdnURL, uniqueFilename)

	return cdnURL, nil
}

func (s *GCSImageStorage) Delete(ctx context.Context, filename string) error {
	// In production, delete from GCS
	return nil
}

func (s *GCSImageStorage) GetURL(filename string) string {
	return fmt.Sprintf("%s/%s", s.cdnURL, filename)
}

// ImageProcessor handles image optimization
type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// OptimizeImage placeholder for image optimization
// In production, this would use proper image processing libraries
func (p *ImageProcessor) OptimizeImage(data []byte, maxWidth, maxHeight uint) ([]byte, error) {
	// In production: resize, compress, and optimize image
	// For now, return original data
	return data, nil
}

// GenerateThumbnail placeholder for thumbnail generation
func (p *ImageProcessor) GenerateThumbnail(data []byte, size uint) ([]byte, error) {
	// In production: create thumbnail
	// For now, return original data
	return data, nil
}
