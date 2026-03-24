package thumbnail

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

type Job struct {
	ImageURL  string `json:"image_url"`
	Sizes     []int  `json:"sizes"`
	Timestamp int64  `json:"timestamp"`
}

func ProcessJob(job Job, fileDir string) error {

	originalPath := filepath.Join(fileDir, job.ImageURL)

	img, err := imaging.Open(originalPath)
	if err != nil {
		return fmt.Errorf("failed to open original image: %w", err)
	}

	dir := filepath.Dir(job.ImageURL)
	base := filepath.Base(job.ImageURL)
	ext := filepath.Ext(base)
	imageId := strings.TrimSuffix(base, ext)

	var wg sync.WaitGroup
	var errOnce sync.Once
	var processErr error

	for _, size := range job.Sizes {
		wg.Add(1)
		go func(s int) {
			defer wg.Done()
			thumb := imaging.Resize(img, s, s, imaging.Lanczos)

			targetPath := filepath.Join(fileDir, dir, fmt.Sprintf("%s_%d%s", imageId, s, ext))

			if err := imaging.Save(thumb, targetPath); err != nil {
				errOnce.Do(func() { processErr = err })
			}
		}(size)
	}
	wg.Wait()
	return processErr
}
