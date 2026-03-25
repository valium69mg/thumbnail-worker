package health

import (
	"log"
	"net/http"
)

func RegisterHealthCheck() {
	port := "8080"
	log.Printf("Starting healthcheck on [::]:%s/healthz", port)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	http.ListenAndServe(":"+port, nil)
}
