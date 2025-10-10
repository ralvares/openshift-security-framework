# frontend-app (Flask)

Simple UI that calls api-gateway for weather lookup and recent forecasts.

## Build
```
IMAGE=YOUR_REGISTRY/frontend-app:dev
podman build -t $IMAGE .
podman push $IMAGE
```

## Deploy
```
oc apply -f manifest.yaml
```

## Port Forward (local test)
```
oc -n ns-frontend-user port-forward svc/frontend-app 5000:5000
open http://localhost:5000
```

## Env Vars
- GATEWAY_URL (default cluster svc URL)

## Future Enhancements
- Add static assets & styling
- Add JWT-based analytics call when gateway supports it
- Add simple caching of last city
