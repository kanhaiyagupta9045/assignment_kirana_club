package image

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type ImageData struct {
	ID     string
	Image  image.Image
	Width  int
	Height int
	Format string
}

func DownloadImage(uri string) (*ImageData, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download image: received status code %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data %w", err)
	}

	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()

	height := bounds.Max.Y
	width := bounds.Max.X

	return &ImageData{
		ID:     uuid.New().String(),
		Image:  img,
		Width:  width,
		Height: height,
		Format: format,
	}, nil
}

func (id *ImageData) SaveImage(documentID, storeID string) error {
	if documentID == "" || storeID == "" {
		return fmt.Errorf("documentID and storeID must be provided")
	}

	dirPath := filepath.Join("./image", documentID, storeID)
	log.Printf("Creating directory: %s", dirPath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Printf("Failed to create directory: %v", err)
		return fmt.Errorf("failed to create directory: %w", err)
	}
	log.Println("Directory created successfully.")

	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.%s", id.ID, id.Format))
	log.Printf("Constructed file path: %s", filePath)

	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}
	defer outFile.Close()
	log.Println("File created successfully.")

	if id.Image == nil {
		return fmt.Errorf("image data is nil")
	}

	switch id.Format {
	case "png":
		log.Printf("Saving PNG file: %s", filePath)
		if err := png.Encode(outFile, id.Image); err != nil {
			log.Printf("Error encoding PNG: %v", err)
			return err
		}
	case "jpeg":
		log.Printf("Saving JPEG file: %s", filePath)
		if err := jpeg.Encode(outFile, id.Image, nil); err != nil {
			log.Printf("Error encoding JPEG: %v", err)
			return err
		}
	default:
		return fmt.Errorf("unsupported image format: %s", id.Format)
	}

	log.Printf("Image saved successfully: %s", filePath)
	return nil
}
