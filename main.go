package main

import (
	"thumbnail-worker/health"
	"thumbnail-worker/worker"
)

func main() {
	go worker.StartWorker()

	go health.RegisterHealthCheck()

	select {}
}
