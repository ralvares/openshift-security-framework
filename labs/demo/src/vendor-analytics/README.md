# vendor-analytics (3rd-party)

Placeholder directory representing an externally supplied container image that requires anyuid SCC.

## Deployment Manifest Example
Apply only if risk accepted and documented.

```
apiVersion: v1
kind: Namespace
metadata:
  name: ns-vendor-analytics
  labels:
    stage: demo
    layer: vendor
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vendor-analytics
  namespace: ns-vendor-analytics
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vendor-analytics
  namespace: ns-vendor-analytics
  labels:
    app: vendor-analytics
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vendor-analytics
  template:
    metadata:
      labels:
        app: vendor-analytics
        security-exception: anyuid
    spec:
      serviceAccountName: vendor-analytics
      containers:
        - name: vendor
          image: VENDOR_REGISTRY/analytics:latest
          securityContext:
            runAsUser: 0 # Required by vendor
            allowPrivilegeEscalation: false
          ports:
            - containerPort: 8080
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 15
---
apiVersion: v1
kind: Service
metadata:
  name: vendor-analytics
  namespace: ns-vendor-analytics
spec:
  selector:
    app: vendor-analytics
  ports:
    - port: 8080
      targetPort: 8080
```

## Grant anyuid SCC (cluster-admin action)
```
oc adm policy add-scc-to-user anyuid -z vendor-analytics -n ns-vendor-analytics
```

Document business justification, owner, review date, and monitoring in your risk register.
