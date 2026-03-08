package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

var ctx = context.Background()

type Job struct {
	ImageURL  string `json:"image_url"`
	Sizes     []int  `json:"sizes"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	fileDir := os.Getenv("FILE_DIRECTORY")

	if redisHost == "" || redisPort == "" || fileDir == "" {
		panic("Missing required environment variables: REDIS_HOST, REDIS_PORT, FILE_DIRECTORY")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	log.Printf("Worker started. Redis at %s:%s, saving files to %s\n", redisHost, redisPort, fileDir)

	for {
		jobStr, err := rdb.RPop(ctx, "thumbnail_jobs").Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			log.Println("Redis error:", err)
			continue
		}

		var job Job
		if err := json.Unmarshal([]byte(jobStr), &job); err != nil {
			log.Println("Invalid job JSON:", err)
			continue
		}

		log.Printf("Processing job: %+v\n", job)
		processJob(job, fileDir)
	}
}

func processJob(job Job, fileDir string) {
	// Full path to original image
	originalPath := filepath.Join(fileDir, job.ImageURL)

	img, err := imaging.Open(originalPath)
	if err != nil {
		log.Println("Failed to open original image:", err)
		return
	}

	// Derive imageId and extension from ImageURL
	base := filepath.Base(job.ImageURL) // e.g. "4c8bad88-ece6-4333-b7bc-cb7f48f71d3d.png"
	ext := filepath.Ext(base)           // ".png"
	imageId := strings.TrimSuffix(base, ext)

	var wg sync.WaitGroup
	for _, size := range job.Sizes {
		wg.Add(1)
		go func(s int) {
			defer wg.Done()
			thumb := imaging.Resize(img, s, s, imaging.Lanczos)

			// Save as categories/{imageId}_{size}.ext
			filename := fmt.Sprintf("%s_%d%s", imageId, s, ext)
			targetPath := filepath.Join(fileDir, "categories", filename)

			if err := imaging.Save(thumb, targetPath); err != nil {
				log.Println("Save failed:", err)
			} else {
				log.Println("Thumbnail saved:", targetPath)
			}
		}(size)
	}
	wg.Wait()
}
