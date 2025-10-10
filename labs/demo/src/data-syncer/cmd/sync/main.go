package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Weather struct {
	City         string   `json:"city"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	TemperatureC *float64 `json:"temperature_C"`
	WindspeedKph *float64 `json:"windspeed_kph"`
	Timestamp    string   `json:"timestamp"`
}

func logEvent(ev map[string]any) {
	if ev["timestamp"] == nil {
		ev["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}
	if ev["component"] == nil {
		ev["component"] = "data-syncer"
	}
	if ev["namespace"] == nil {
		ev["namespace"] = os.Getenv("POD_NAMESPACE")
	}
	b, _ := json.Marshal(ev)
	fmt.Println(string(b))
}

func main() {
	legacyURL := env("WEATHER_PROXY_URL", "http://httpd-cgi.ns-frontend-user.svc.cluster.local")
	forecastAPI := env("FORECAST_BASE", "http://api-forecast.ns-api-user.svc.cluster.local:8000")
	cityList := env("SYNC_CITIES", "London,Dubai,New York")
	jwtSecret := env("JWT_SECRET", "dev-change-me")
	jwtSub := env("JWT_SUBJECT", "data-syncer")
	token, err := buildJWT(jwtSecret, jwtSub, 15*time.Minute)
	if err != nil {
		log.Fatalf("cannot build jwt: %v", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	processed := 0
	errs := 0
	for _, city := range strings.Split(cityList, ",") {
		city = strings.TrimSpace(city)
		if city == "" {
			continue
		}
		w, err := fetchWeather(client, legacyURL, city)
		if err != nil {
			logEvent(map[string]any{"event_type": "egress", "method": "GET", "path": "/weather", "dst_ip": legacyURL, "status": 0, "message": "fetch_error", "error": err.Error(), "city": city})
			continue
		}
		if err := pushForecast(client, forecastAPI, w, token); err != nil {
			logEvent(map[string]any{"event_type": "egress", "method": "POST", "path": "/forecast", "dst_ip": forecastAPI, "status": 0, "message": "push_error", "error": err.Error(), "city": city})
			errs++
		} else {
			logEvent(map[string]any{"event_type": "egress", "method": "POST", "path": "/forecast", "dst_ip": forecastAPI, "status": 201, "message": "push_ok", "city": city})
			processed++
		}
	}
	logEvent(map[string]any{"event_type": "misc", "message": "sync_complete", "processed": processed, "errors": errs})
}

func fetchWeather(client *http.Client, base, city string) (*Weather, error) {
	url := fmt.Sprintf("%s/?city=%s", strings.TrimSuffix(base, "/"), city)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var w Weather
	if err := json.Unmarshal(body, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func pushForecast(client *http.Client, base string, w *Weather, token string) error {
	obsTime := w.Timestamp
	payload := map[string]any{
		"city":          w.City,
		"latitude":      w.Latitude,
		"longitude":     w.Longitude,
		"temperature_c": w.TemperatureC,
		"windspeed_kph": w.WindspeedKph,
		"observed_at":   obsTime,
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, strings.TrimSuffix(base, "/")+"/forecast", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("push failed status %d", resp.StatusCode)
	}
	return nil
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func buildJWT(secret, sub string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now.Add(-1 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
