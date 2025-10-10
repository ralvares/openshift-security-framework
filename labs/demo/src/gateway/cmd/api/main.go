package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Forecast struct {
	ID           int       `json:"id"`
	City         string    `json:"city"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	TemperatureC *float64  `json:"temperature_C"`
	WindspeedKph *float64  `json:"windspeed_kph"`
	ObservedAt   time.Time `json:"observed_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Claims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

// ===== Simple structured logging (minimal, stdout) =====
type ctxKey string

const userKey ctxKey = "userID"

func logEvent(ev map[string]any) {
	// required baseline fields
	if ev["timestamp"] == nil {
		ev["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}
	if ev["component"] == nil {
		ev["component"] = "api-gateway"
	}
	if ev["namespace"] == nil {
		ev["namespace"] = os.Getenv("POD_NAMESPACE")
	}
	if ev["user_id"] == nil {
		ev["user_id"] = "anonymous"
	}
	if ev["event_type"] == nil {
		ev["event_type"] = "misc"
	}
	b, _ := json.Marshal(ev)
	fmt.Println(string(b))
}

func withUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, userKey, user)
}
func userFrom(ctx context.Context) string {
	if v, ok := ctx.Value(userKey).(string); ok && v != "" {
		return v
	}
	return "anonymous"
}

// ingress logging middleware
func ingress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &logWriter{ResponseWriter: w, status: 200}
		logEvent(map[string]any{
			"event_type": "ingress", "method": r.Method, "path": r.URL.Path,
			"src_ip": r.RemoteAddr, "message": "request_start", "status": 0,
			"user_id": userFrom(r.Context()),
		})
		next.ServeHTTP(lw, r)
		logEvent(map[string]any{
			"event_type": "ingress", "method": r.Method, "path": r.URL.Path,
			"src_ip": r.RemoteAddr, "message": "request_complete", "status": lw.status,
			"latency_ms": time.Since(start).Milliseconds(),
			"user_id":    userFrom(r.Context()),
		})
	})
}

type logWriter struct {
	http.ResponseWriter
	status int
}

func (lw *logWriter) WriteHeader(c int) { lw.status = c; lw.ResponseWriter.WriteHeader(c) }

func main() {
	addr := getEnv("GATEWAY_ADDR", ":8080")
	jwtSecret := getEnv("JWT_SECRET", "dev-change-me")
	forecastURL := getEnv("FORECAST_BASE", "http://api-forecast.ns-api-user.svc.cluster.local:8000")
	weatherProxyURL := getEnv("WEATHER_PROXY_URL", "http://httpd-cgi.ns-frontend-user.svc.cluster.local")
	vendorURL := getEnv("VENDOR_ANALYTICS_URL", "http://vendor-analytics.ns-vendor-analytics.svc.cluster.local:443")

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})

	mux.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method"})
			return
		}
		city := r.URL.Query().Get("city")
		if city == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing city"})
			return
		}
		// proxy to legacy CGI
		url := fmt.Sprintf("%s/?city=%s", strings.TrimSuffix(weatherProxyURL, "/"), city)
		resp, err := httpClient.Get(url)
		logEvent(map[string]any{"event_type": "egress", "method": "GET", "path": "/weather-proxy", "dst_ip": weatherProxyURL, "status": func() {
			if resp != nil {
				return resp.StatusCode
			}
			return 0
		}(), "message": "proxy_weather"})
		if err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": "upstream"})
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	})

	// JWT-protected analytics proxy -> vendor-analytics
	mux.Handle("/analytics", authMiddleware(jwtSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method"})
			return
		}
		url := strings.TrimSuffix(vendorURL, "/") + "/analytics"
		resp, err := httpClient.Get(url)
		logEvent(map[string]any{"event_type": "egress", "method": "GET", "path": "/vendor-analytics", "dst_ip": vendorURL, "status": func() {
			if resp != nil {
				return resp.StatusCode
			}
			return 0
		}(), "message": "proxy_analytics"})
		if err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": "upstream"})
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	})))

	// Public list forecasts
	mux.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		base := strings.TrimSuffix(forecastURL, "/")
		switch r.Method {
		case http.MethodGet:
			url := fmt.Sprintf("%s/forecast?limit=20", base)
			resp, err := httpClient.Get(url)
			logEvent(map[string]any{"event_type": "egress", "method": "GET", "path": "/forecast list", "dst_ip": forecastURL, "status": func() {
				if resp != nil {
					return resp.StatusCode
				}
				return 0
			}(), "message": "list_forecast"})
			if err != nil {
				writeJSON(w, http.StatusBadGateway, map[string]string{"error": "upstream"})
				return
			}
			defer resp.Body.Close()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			_, _ = io.Copy(w, resp.Body)
		case http.MethodPost:
			// require auth
			authMiddleware(jwtSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				forwardJSON(w, r, base+"/forecast")
			})).ServeHTTP(w, r)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method"})
		}
	})

	mux.Handle("/forecast/", authMiddleware(jwtSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Expect /forecast/{id}
		base := strings.TrimSuffix(forecastURL, "/")
		idPart := strings.TrimPrefix(r.URL.Path, "/forecast/")
		if idPart == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing id"})
			return
		}
		backend := base + "/forecast/" + idPart
		switch r.Method {
		case http.MethodGet:
			resp, err := httpClient.Get(backend)
			logEvent(map[string]any{"event_type": "egress", "method": "GET", "path": "/forecast/{id}", "dst_ip": forecastURL, "status": func() {
				if resp != nil {
					return resp.StatusCode
				}
				return 0
			}(), "message": "get_forecast"})
			if err != nil {
				writeJSON(w, http.StatusBadGateway, map[string]string{"error": "upstream"})
				return
			}
			defer resp.Body.Close()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			_, _ = io.Copy(w, resp.Body)
		case http.MethodPut, http.MethodDelete:
			forwardJSON(w, r, backend)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method"})
		}
	})))

	logEvent(map[string]any{"event_type": "listen", "message": "listening_ports", "listening": []map[string]any{{"port": 8080, "proto": "tcp"}}})
	if err := http.ListenAndServe(addr, ingress(mux)); err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims := &Claims{}
		tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("alg")
			}
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}
		r = r.WithContext(withUser(r.Context(), claims.Sub))
		logEvent(map[string]any{"event_type": "auth", "user_id": claims.Sub, "message": "jwt_valid", "status": 200})
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// legacy logRequest removed (superseded by structured ingress)

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// forwardJSON forwards the current request (method + body) as JSON to target URL
func forwardJSON(w http.ResponseWriter, r *http.Request, target string) {
	body := r.Body
	defer body.Close()
	req, err := http.NewRequest(r.Method, target, body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "build"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// propagate auth if present (already validated)
	if auth := r.Header.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "upstream"})
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, _ = http.Copy(w, resp.Body)
}
