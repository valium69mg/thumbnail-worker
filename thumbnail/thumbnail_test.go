package thumbnail

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// helper: crea una imagen de prueba en path dado
func createTestImage(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	// Rellenar con blanco
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.White)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func TestProcessJobValidImage(t *testing.T) {
	tmpDir := t.TempDir()

	// Crear subcarpeta como en ImageURL
	imagePath := filepath.Join(tmpDir, "categories", "uuid.png")
	if err := createTestImage(imagePath); err != nil {
		t.Fatalf("failed to create test image: %v", err)
	}

	job := Job{
		ImageURL:  "categories/uuid.png",
		Sizes:     []int{5, 8},
		Timestamp: 123,
	}

	err := ProcessJob(job, tmpDir)
	if err != nil {
		t.Fatalf("ProcessJob failed: %v", err)
	}

	// Verificar que se crearon los thumbnails
	for _, size := range job.Sizes {
		expected := filepath.Join(tmpDir, "categories", "uuid_"+strconv.Itoa(size)+".png")
		if _, err := os.Stat(expected); os.IsNotExist(err) {
			t.Errorf("expected thumbnail not found: %s", expected)
		}
	}
}

func TestProcessJobFileNotExist(t *testing.T) {
	tmpDir := t.TempDir()

	job := Job{
		ImageURL:  "categories/nonexistent.png",
		Sizes:     []int{10},
		Timestamp: 123,
	}

	// no debe panic, debe retornar error
	err := ProcessJob(job, tmpDir)
	if err == nil {
		t.Errorf("expected error for nonexistent file, got nil")
	}
}
