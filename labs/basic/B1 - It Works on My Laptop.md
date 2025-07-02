# Lab B1: It Works on My Laptop

## Skill

Understand OpenShift’s secure-by-default behavior

## Objective

- Deploy a container using the `httpd:2.4` image and observe failure
- Investigate the failure using logs and container metadata
- Replace the image with `registry.access.redhat.com/ubi9/httpd-24`
- Explain why one image fails and the other succeeds
- Deploy and access the working application

## Background

OpenShift enforces security through defaults like running containers as non-root and disallowing privileged ports (<1024). Containers that expect to run as root or bind to port 80 will fail at runtime, not during deployment.

## Instructions

### 1. Create a namespace

```sh
kubectl create namespace httpd-demo
kubectl config set-context --current --namespace=httpd-demo
```

### 2. Deploy the privileged image

```sh
kubectl create deployment httpd-privileged --image=httpd:2.4
```

### 3. Check pod status and logs

```sh
kubectl get pods
kubectl describe pod -l app=httpd-privileged
kubectl logs -l app=httpd-privileged
```

**Expected output in logs:**

```
Permission denied: AH00072: make_sock: could not bind to address [::]:80
no listening sockets available, shutting down
```

### 4. Remove the failing deployment

```sh
kubectl delete deployment httpd-privileged
```

### 5. Deploy the secure UBI image

```sh
kubectl create deployment httpd-unprivileged --image=registry.access.redhat.com/ubi9/httpd-24
kubectl expose deployment httpd-unprivileged --port=8080 --target-port=8080
```

### 6. Expose the application using a Route

```sh
kubectl expose deployment httpd-unprivileged --port=8080 --target-port=8080 --type=ClusterIP
oc create route edge --service=httpd-unprivileged --port=8080
```

Get the route URL:

```sh
oc get route httpd-unprivileged
```

Test access (replace <ROUTE_URL> with the actual route hostname):

```sh
curl http://<ROUTE_URL>
```

### 7. Confirm UID again

```sh
kubectl exec -it deploy/httpd-unprivileged -- id
```

You should see a UID like 1001, which is non-root.

## Summary

| Image                | UID  | Port | Works on OpenShift | Reason                        |
|----------------------|------|------|--------------------|-------------------------------|
| httpd:2.4 (httpd-privileged) | 0    | 80   | No                 | Requires root and port 80      |
| ubi9/httpd-24 (httpd-unprivileged)   | 1001 | 8080 | Yes                | Runs as non-root               |

## Key Points

- OpenShift enforces non-root execution and blocks privileged ports by default
- Containers that assume root will fail silently at runtime
- Use images designed for non-root operation and high ports