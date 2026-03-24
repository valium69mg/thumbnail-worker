package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"thumbnail-worker/thumbnail"
)

var ctx = context.Background()

func StartWorker() {
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

		var job thumbnail.Job
		if err := json.Unmarshal([]byte(jobStr), &job); err != nil {
			log.Println("Invalid job JSON:", err)
			continue
		}

		log.Printf("Processing job: %+v\n", job)
		err = thumbnail.ProcessJob(job, fileDir)
		if err != nil {
			log.Println("Failed to process job:", err)
		} else {
			log.Println("Job processed successfully")
		}
	}
}
