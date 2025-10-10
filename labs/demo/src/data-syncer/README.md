# data-syncer (Go)

One-shot binary that:
1. Pulls current weather for a list of cities from the legacy weather-proxy CGI
2. Pushes each observation into the api-forecast service as a new forecast row

## Env Vars
- WEATHER_PROXY_URL (default http://httpd-cgi.ns-frontend-user.svc.cluster.local)
- FORECAST_BASE (default http://api-forecast.ns-api-user.svc.cluster.local:8000)
- SYNC_CITIES (comma list, default: London,Dubai,New York)

## Build
```
IMAGE=YOUR_REGISTRY/data-syncer:dev
podman build -t $IMAGE .
podman push $IMAGE
```

## Run locally (requires access to cluster network or port-forwards)
```
go run ./cmd/sync
```

## CronJob Example
```
apiVersion: batch/v1
kind: CronJob
metadata:
  name: data-syncer
  namespace: ns-api-user
spec:
  schedule: "*/30 * * * *" # every 30 minutes
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: sync
              image: YOUR_REGISTRY/data-syncer:dev
              env:
                - name: SYNC_CITIES
                  value: "London,Dubai,New York,Singapore"
```
