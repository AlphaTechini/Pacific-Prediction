package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"prediction/internal/realtime"
)

func NewStreamHandler(controller realtime.Controller, heartbeatInterval time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "streaming_not_supported"})
			return
		}

		subscription, err := controller.Subscribe(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		defer subscription.Close()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		if _, err := fmt.Fprint(w, ": stream connected\n\n"); err != nil {
			return
		}
		flusher.Flush()

		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-ticker.C:
				if _, err := fmt.Fprint(w, ": heartbeat\n\n"); err != nil {
					return
				}
				flusher.Flush()
			case event, ok := <-subscription.Events:
				if !ok {
					return
				}

				if err := writeStreamEvent(w, event); err != nil {
					return
				}
				flusher.Flush()
			}
		}
	})
}

func writeStreamEvent(w http.ResponseWriter, event realtime.StreamEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal stream event: %w", err)
	}

	if _, err := fmt.Fprintf(w, "event: %s\n", event.Type); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", body); err != nil {
		return err
	}

	return nil
}
