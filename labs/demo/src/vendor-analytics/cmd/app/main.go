package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	At    string  `json:"at"`
}

func logEvent(ev map[string]any) {
	if ev["timestamp"] == nil {
		ev["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}
	if ev["component"] == nil {
		ev["component"] = "vendor-analytics"
	}
	if ev["namespace"] == nil {
		ev["namespace"] = os.Getenv("POD_NAMESPACE")
	}
	b, _ := json.Marshal(ev)
	fmt.Println(string(b))
}

func main() {
	// Bind to privileged port (requires root & anyuid SCC in OpenShift)
	addr := ":443"
	if v := os.Getenv("VENDOR_PORT"); v != "" {
		addr = ":" + v
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logEvent(map[string]any{"event_type": "ingress", "method": r.Method, "path": "/health", "src_ip": r.RemoteAddr, "message": "health"})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	mux.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		logEvent(map[string]any{"event_type": "ingress", "method": r.Method, "path": "/analytics", "src_ip": r.RemoteAddr, "message": "analytics_req"})
		w.Header().Set("Content-Type", "application/json")
		out := []Metric{{Name: "requests_total", Value: 1234, At: time.Now().UTC().Format(time.RFC3339)}}
		json.NewEncoder(w).Encode(out)
	})
	logEvent(map[string]any{"event_type": "listen", "message": "listening_ports", "listening": []map[string]any{{"port": 443, "proto": "tcp"}}})
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
