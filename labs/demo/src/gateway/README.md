# api-gateway (Go)

Minimal gateway providing:
- /healthz (public)
- /weather (proxies legacy weather-proxy CGI)
- /forecast (lists forecasts from api-forecast)
- /analytics (JWT protected placeholder)

## Env Vars
- GATEWAY_ADDR (:8080 default)
- JWT_SECRET (HMAC secret)
- FORECAST_BASE (default http://api-forecast.ns-api-user.svc.cluster.local:8000)
- WEATHER_PROXY_URL (default http://httpd-cgi.ns-frontend-user.svc.cluster.local)

## Build & Run (local)
```
go build -o api-gateway ./cmd/api
JWT_SECRET=dev token example:
python - <<'PY'
import jwt, time
print(jwt.encode({"sub":"tester","exp":int(time.time())+3600},"dev-change-me",algorithm="HS256"))
PY
./api-gateway
```

## Container
```
IMAGE=YOUR_REGISTRY/api-gateway:dev
podman build -t $IMAGE .
podman push $IMAGE
```

## Kubernetes (example Deployment)
Add to your manifests:
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: ns-api-user
spec:
  replicas: 1
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
        - name: api-gateway
          image: YOUR_REGISTRY/api-gateway:dev
          env:
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: jwt-secret
                  key: secret
          ports:
            - containerPort: 8080
          readinessProbe:
            httpGet:
              path: /healthz
              port: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: api-gateway
  namespace: ns-api-user
spec:
  selector:
    app: api-gateway
  ports:
    - port: 8080
      targetPort: 8080
```
