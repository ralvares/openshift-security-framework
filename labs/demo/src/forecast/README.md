# api-forecast Demo Component

This directory contains a minimal FastAPI service and a MySQL backing store to support the broader weather platform.

## Contents
- `requirements.txt` Python deps
- `Dockerfile` UBI9-based image build
- `app/` FastAPI code (DB connection + CRUD)
- `manifest.yaml` Namespaces, MySQL, and api-forecast Deployment/Service

## Build Image
```
IMAGE=YOUR_REGISTRY/api-forecast:dev
podman build -t $IMAGE .
podman push $IMAGE
```

## Deploy (creates namespaces if absent)
```
oc apply -f manifest.yaml
```

## Initialize Schema
Port-forward the DB or run a Job. Example local forward:
```
oc -n ns-data-user port-forward deploy/mysql 3306:3306 &
mysql -h 127.0.0.1 -u root -p$MYSQL_ROOT_PASSWORD <<'SQL'
CREATE TABLE IF NOT EXISTS forecast (
  id INT AUTO_INCREMENT PRIMARY KEY,
  city VARCHAR(64) NOT NULL,
  latitude DECIMAL(8,5) NOT NULL,
  longitude DECIMAL(8,5) NOT NULL,
  temperature_c DECIMAL(5,2),
  windspeed_kph DECIMAL(5,2),
  observed_at DATETIME NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_city_time (city, observed_at)
);
SQL
```

## Test API
```
oc -n ns-api-user port-forward deploy/api-forecast 8000:8000 &
curl localhost:8000/healthz
curl -X POST localhost:8000/forecast \
  -H 'Content-Type: application/json' \
  -d '{"city":"Dubai","latitude":25.0657,"longitude":55.1713,"temperature_c":34.2,"windspeed_kph":12.3,"observed_at":"2025-10-06T12:00:00"}'
```

## Next Steps
- Add NetworkPolicies
- Add ServiceAccount + RBAC
- Add migrations Job instead of manual schema
- Introduce api-gateway with JWT validation
- Implement data-syncer CronJob
