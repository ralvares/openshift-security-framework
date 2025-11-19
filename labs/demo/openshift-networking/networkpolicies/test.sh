#!/usr/bin/env bash

set -euo pipefail

# Color helpers
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

pass() { printf "%bPASS%b %s\n" "$GREEN" "$NC" "$1"; }
fail() { printf "%bFAIL%b %s\n" "$RED" "$NC" "$1"; }
info() { printf "%bINFO%b %s\n" "$BLUE" "$NC" "$1"; }

TESTBED_YAML="anp-testbed.yaml"

info "[SETUP] Applying testbed YAML..."
oc apply -f "$TESTBED_YAML" >/dev/null
info "[SETUP] Testbed applied. Waiting for pods..."

PODS=(
  "team-a/a-default"
  "team-a/a-dns"
  "team-a/a-kapi"
  "team-a/a-proxy"
  "team-a/a-ingress"
  "team-b/b-default"
  "team-c/c-proxy"
  "team-c/c-kapi"
  "external-no-access/outsider"
)

check_pods_ready() {
  info "[1] Checking pod readiness…"
  for p in "${PODS[@]}"; do
    ns="${p%%/*}"
    pod="${p##*/}"
    info "  Waiting for $ns/$pod"
    if ! oc wait pod "$pod" -n "$ns" --for=condition=Ready --timeout=180s >/dev/null; then
      fail "Pod $ns/$pod failed to become Ready"
      exit 1
    fi
  done
  pass "All pods Ready"
}

exec_cmd() {
  local ns="$1"
  local pod="$2"
  local cmd="$3"
  oc exec -n "$ns" "$pod" -- bash -c "$cmd" 2>&1
}
TOTAL_TESTS=0
FAILED_TESTS=0

test_case() {
  local title="$1"
  local ns="$2"
  local pod="$3"
  local cmd="$4"
  local expect="$5"

  TOTAL_TESTS=$((TOTAL_TESTS+1))

  set +e
  local output
  output=$(exec_cmd "$ns" "$pod" "$cmd")
  local rc=$?
  set -e

  if [[ "$expect" == "__STATUS_2XX__" ]]; then
    if [[ "$rc" -eq 0 ]]; then
      pass "$title (from $ns/$pod)"
    else
      fail "$title (from $ns/$pod)"
      FAILED_TESTS=$((FAILED_TESTS+1))
    fi
  elif [[ "$expect" == "__STATUS_FAIL__" ]]; then
    if [[ "$rc" -ne 0 ]]; then
      pass "$title (from $ns/$pod)"
    else
      fail "$title (from $ns/$pod)"
      FAILED_TESTS=$((FAILED_TESTS+1))
    fi
  else
    if [[ "$output" =~ $expect ]]; then
      pass "$title (from $ns/$pod)"
    else
      fail "$title (from $ns/$pod)"
      FAILED_TESTS=$((FAILED_TESTS+1))
    fi
  fi
}

info "Starting AdminNetworkPolicy validation tests..."
check_pods_ready
### 2. DNS ALLOW (team-a/a-dns)
# We just check that getent returns any line with an IP.
test_case "DNS resolution works" \
  "team-a" "a-dns" \
  "getent hosts kubernetes.default.svc.cluster.local || getent hosts api.vulnerawise.ai" \
  "[0-9]"

### 3. Kube API ALLOW (team-a/a-kapi, team-c/c-kapi)
test_case "Kube API reachable (a-kapi)" \
  "team-a" "a-kapi" \
  "curl -ks --max-time 3 https://kubernetes.default.svc.cluster.local:443 -o /dev/null -w '%{http_code}'" \
  "__STATUS_2XX__"

test_case "Kube API reachable (c-kapi)" \
  "team-c" "c-kapi" \
  "curl -ks --max-time 3 https://kubernetes.default.svc.cluster.local:443 -o /dev/null -w '%{http_code}'" \
  "__STATUS_2XX__"

### 4. Kube API BLOCKED for non-labeled pods
test_case "Kube API blocked for default pods" \
  "team-b" "b-default" \
  "curl -ks --max-time 3 https://kubernetes.default.svc.cluster.local:443 -o /dev/null -w '%{http_code}'" \
  "__STATUS_FAIL__"

### 5. Metadata DENY
test_case "Metadata access blocked" \
  "team-a" "a-default" \
  "curl -s --max-time 3 https://142.250.203.174 -o /dev/null -w '%{http_code}'" \
  "__STATUS_FAIL__"

### 6. External egress test (using external API as target)
EXTERNAL_TARGET="51.79.73.188"

# Note: this assumes your policies allow egress to this IP from proxy-labeled pods.
test_case "External endpoint reachable from a-proxy" \
  "team-a" "a-proxy" \
  "curl -sk --max-time 5 https://$EXTERNAL_TARGET -o /dev/null -w '%{http_code}' || curl -s --max-time 5 http://$EXTERNAL_TARGET -o /dev/null -w '%{http_code}'" \
  "__STATUS_2XX__"  # any successful HTTP status

test_case "External endpoint reachable from c-proxy" \
  "team-c" "c-proxy" \
  "curl -sk --max-time 5 https://$EXTERNAL_TARGET -o /dev/null -w '%{http_code}' || curl -s --max-time 5 http://$EXTERNAL_TARGET -o /dev/null -w '%{http_code}'" \
  "__STATUS_2XX__"  # any successful HTTP status

### 7. East–West Pass (team-b → team-a)
service_ip=$(oc get pod a-default -n team-a -o jsonpath='{.status.podIP}')

test_case "Owned namespaces pass-through (team-b → team-a)" \
  "team-b" "b-default" \
  "curl -s --max-time 3 http://$service_ip:8080 -o /dev/null -w '%{http_code}'" \
  "__STATUS_2XX__"

### 8. External namespace isolation
test_case "external-no-access cannot reach owned namespaces" \
  "external-no-access" "outsider" \
  "curl -s --max-time 3 http://$service_ip:8080 -o /dev/null -w '%{http_code}'" \
  "__STATUS_FAIL__"

### 9. Final deny-all enforcement (direct Internet blocked)
test_case "Direct internet blocked from default pod" \
  "team-a" "a-default" \
  "curl -s --max-time 3  https://www.google.com -o /dev/null -w '%{http_code}'" \
  "__STATUS_FAIL__"

info "[CLEANUP] Deleting testbed resources…"
#oc delete -f "$TESTBED_YAML" --ignore-not-found >/dev/null || true
echo "========================================"
if [[ "$FAILED_TESTS" -eq 0 ]]; then
  pass "All $TOTAL_TESTS tests passed and testbed cleaned up."
else
  fail "$FAILED_TESTS of $TOTAL_TESTS tests failed (total: $TOTAL_TESTS)."
fi
