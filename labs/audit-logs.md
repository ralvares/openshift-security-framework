# 🔥 Kubernetes Tier 1 Threat Hunting Demo with Audit Logs

This is a full step-by-step demo for hunting Kubernetes threats using audit logs on an OpenShift cluster with access to `/var/log/kube-apiserver/audit.log`. It simulates a realistic attacker scenario where a pod is compromised and its ServiceAccount is used to escalate privileges and perform malicious actions.

All actions happen in the **`default` namespace**, and detections rely **only on verbs, resources, and subresources** — no assumptions about usernames, pod names, or namespaces.

---

## ✅ Step 0: Create Pod with kubectl and curl (Attacker foothold)

```bash
oc run testpod \
  --image=bitnami/kubectl:latest \
  -n default \
  --restart=Never \
  --command -- sleep infinity
```

### 🔍 Show logs (pod creation)

```bash
grep '"resource":"pods"' /var/log/kube-apiserver/audit.log | \
grep '"verb":"create"' | jq 'select(.objectRef.namespace=="default") | {
  timestamp: .requestReceivedTimestamp,
  user: .user.username,
  verb: .verb,
  resource: .objectRef.resource,
  namespace: .objectRef.namespace,
  name: .objectRef.name,
  uri: .requestURI
} | with_entries(select(.value != null))'
```

---

## ✅ Step 1: Exec into the pod (initial access)

```bash
oc exec -it testpod -n default -- sh
```

### 🔍 Show logs (exec access)

```bash
grep '"subresource":"exec"' /var/log/kube-apiserver/audit.log | \
grep '"verb":"get"' | jq 'select(.objectRef.namespace=="default") | {
  timestamp: .requestReceivedTimestamp,
  user: .user.username,
  verb: .verb,
  subresource: .objectRef.subresource,
  resource: .objectRef.resource,
  namespace: .objectRef.namespace,
  name: .objectRef.name,
  uri: .requestURI
} | with_entries(select(.value != null))'
```

---

## ✅ Step 2: Escalate privileges (bind SA to cluster-admin)

```bash
oc adm policy add-cluster-role-to-user cluster-admin system:serviceaccount:default:default
```

### 🔍 Show logs (role binding escalation)

```bash
grep '"resource":"clusterrolebindings"' /var/log/kube-apiserver/audit.log | \
grep -E '"verb":"(create|patch|update)"' | jq '{
  timestamp: .requestReceivedTimestamp,
  user: .user.username,
  verb: .verb,
  resource: .objectRef.resource,
  uri: .requestURI,
  response_code: .responseStatus.code,
  decision: .annotations["authorization.k8s.io/decision"]
} | with_entries(select(.value != null))'
```

---

## ✅ Step 3: From inside the pod, list nodes using kubectl

Inside the pod:

```sh
kubectl get nodes --token=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token) \
  --server=https://kubernetes.default.svc --insecure-skip-tls-verify
```

### 🔍 Show logs (API access with SA)

```bash
grep '"user":{"username":"system:serviceaccount:' /var/log/kube-apiserver/audit.log | \
grep -v '"user":{"username":"system:serviceaccount:openshift-' | \
grep -E '"verb":"(get|list|create|patch|update)"' | \
jq '{
  timestamp: .requestReceivedTimestamp,
  user: .user.username,
  verb: .verb,
  resource: .objectRef.resource,
  subresource: .objectRef.subresource,
  uri: .requestURI,
  response_code: .responseStatus.code
} | with_entries(select(.value != null))'
```

---

## ✅ Step 4: Create a CronJob from inside the pod

```bash
kubectl create cronjob eviljob --image=busybox --schedule="*/1 * * * *" -- echo pwned
```

### 🔍 Show logs (CronJob creation)

```bash
grep '"resource":"cronjobs"' /var/log/kube-apiserver/audit.log | \
grep -E '"verb":"(create|patch|update)"' | \
grep -v '"user":{"username":"system:serviceaccount:openshift-' | \
grep -v '"user":{"username":"system:serviceaccount:kube-system:' | \
jq '{
  timestamp: .requestReceivedTimestamp,
  user: .user.username,
  verb: .verb,
  resource: .objectRef.resource,
  namespace: .objectRef.namespace,
  name: .objectRef.name,
  uri: .requestURI,
  response_code: .responseStatus.code
} | with_entries(select(.value != null))'
```

---

## ✅ Step 5: Port-forward to internal service

From within the pod:

```bash
kubectl port-forward testpod 9000:80 --token=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token) \
  --server=https://kubernetes.default.svc --insecure-skip-tls-verify
```

### 🔍 Show logs (portforward)

```bash
grep '"subresource":"portforward"' /var/log/kube-apiserver/audit.log | \
grep '"verb":"get"' | \
grep -v '"user":{"username":"system:serviceaccount:openshift-' | \
grep -v '"user":{"username":"system:serviceaccount:kube-system:' | \
jq '{
  timestamp: .requestReceivedTimestamp,
  user: .user.username,
  verb: .verb,
  subresource: .objectRef.subresource,
  resource: .objectRef.resource,
  namespace: .objectRef.namespace,
  name: .objectRef.name,
  uri: .requestURI,
  response_code: .responseStatus.code
} | with_entries(select(.value != null))'
```

---

This sequence simulates:

1. Attacker entry into a pod
2. Privilege escalation
3. API abuse from inside the pod
4. Persistence via CronJob
5. Lateral movement or tunneling with port-forward

All tracked and confirmed via audit logs — using **only verbs and resource types**, and scoped to the **default namespace** for clarity and reproducibility.
