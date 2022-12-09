package health

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var HEALTHY = true

func healthHandler(w http.ResponseWriter, r *http.Request) {
	body := json.RawMessage("OK")
	if !HEALTHY {
		body = json.RawMessage("NOT OK")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err := w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error writing response body: %s", err)
	}
}

func InitHealth() {
	healthServer := http.NewServeMux()
	healthServer.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r)
	})

	healthServer.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r)
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", 8081), healthServer); err != nil {
			log.Fatal(err)
		}
	}()
}
