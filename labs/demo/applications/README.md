# Demo App

## Overview
This is an HTTPD server used to simulate network traffic. It listens on port `8080` and runs on the `UBI 9` base image with `HTTPD 2.4`.

## Features
- Simulates network traffic for testing and demonstration purposes.
- Listens on port `8080`.
- Runs on a Universal Base Image (UBI 9) with HTTPD version 2.4.

## Usage
The application is automatically patched and configured. Below is an example of how to check the connectivity status using the `connectivity-probe` deployment:

```bash
oc exec deploy/team-a1-connectivity-probe -n team-a1 -- curl -s http://localhost:8080
```

### Example Output
```
Name          | Host                                          | Port | Status | Checked At
------------------------------------------------------------------------------------------
team-b1       | httpd.team-b1.svc.cluster.local               | 8080 | ✅      | 07:10:04  
weather-proxy | weather-proxy.weather-proxy.svc.cluster.local | 80   | ❌      | 07:10:04  
```

The output provides the following details:
- **Name**: The name of the service.
- **Host**: The service's hostname.
- **Port**: The port the service is listening on.
- **Status**: The connectivity status (✅ for success, ❌ for failure).
- **Checked At**: The timestamp of the connectivity check.

## Using the Connectivity Probe for Team A1

To check the connectivity status for `team-a1`, you can use the following command:

```bash
oc exec deploy/team-a1-connectivity-probe -n team-a1 -- curl -s http://localhost:8080
```

This will return an ASCII table with the connectivity status, similar to the example below:

```
Name       | Host                                                 | Port | Status | Checked At
----------------------------------------------------------------------------------------------
team-a1    | team-a1-connectivity-probe.team-a1.svc.cluster.local | 80   | ✅      | 10:20:17  
team-a2    | team-a2-connectivity-probe.team-a2.svc.cluster.local | 80   | ❌      | 10:20:17  
team-b1    | team-b1-connectivity-probe.team-b1.svc.cluster.local | 80   | ❌      | 10:20:17  
team-b2    | team-b2-connectivity-probe.team-b2.svc.cluster.local | 80   | ❌      | 10:20:17  
open-meteo | api.open-meteo.com                                   | 443  | ✅      | 10:20:17  
[//]: # (Demo README — friendly, step-by-step guide)
# Connectivity Demo — applications

This README explains how to run the connectivity demo for the applications in this repo. It includes exact commands to apply each network overlay (OVS, UDN, CUDN), how to validate with the connectivity probes, expected output examples, and cleanup/troubleshooting tips.

Prerequisites
- You have access to an OpenShift cluster and `oc` is configured and authenticated.
- You have cluster privileges required to apply/delete the overlay kustomize directories used by this demo.

Repository overlays
- labs/demo/applications/kustomize/overlays/ovs — flat OVS network (everyone on the same flat network)
- labs/demo/applications/kustomize/overlays/udn — a unique UDN applied for team-a (isolates team-a)
- labs/demo/applications/kustomize/overlays/cudn — a common UDN for teams b1 & b2 (groups B1/B2 together)

Quick summary of the demo scenarios
- OVS: everyone on the same flat network. Expect all teams to reach each other (subject to service ports).
- UDN: team-a runs on its own UDN — it can reach external services but may be isolated from other teams depending on policy.
- CUDN: teams b1 and b2 share a common UDN (CUDN) and can reach each other.

Core commands
- Apply an overlay (example: OVS):

```bash
oc apply -k labs/demo/applications/kustomize/overlays/ovs
```

- Delete the same overlay when finished:

```bash
oc delete -k labs/demo/applications/kustomize/overlays/ovs
```

- Check that the connectivity-probe deployments are running (replace `team-a1`):

```bash
oc -n team-a1 get deploy team-a1-connectivity-probe
oc -n team-a1 get pods -l app=connectivity-probe
```

- Exec into the probe and print the ASCII table from the local HTTP probe (example for team-a1):

```bash
oc exec deploy/team-a1-connectivity-probe -n team-a1 -- curl -s http://localhost:8080
```

Recommended demo flow (copy/paste friendly)

1) Start with the flat OVS overlay

```bash
oc apply -k labs/demo/applications/kustomize/overlays/ovs
# wait for deployments to become ready (watch pods)
oc get pods -A -l app=connectivity-probe
```

Run the probes from each team (these print the ASCII table):

```bash
oc exec deploy/team-a1-connectivity-probe -n team-a1 -- curl -s http://localhost:8080
oc exec deploy/team-a2-connectivity-probe -n team-a2 -- curl -s http://localhost:8080
oc exec deploy/team-b1-connectivity-probe -n team-b1 -- curl -s http://localhost:8080
oc exec deploy/team-b2-connectivity-probe -n team-b2 -- curl -s http://localhost:8080
```

Expected example (OVS / flat network):

```
Name       | Host                                                 | Port | Status | Checked At
----------------------------------------------------------------------------------------------
team-a1    | team-a1-connectivity-probe.team-a1.svc.cluster.local | 80   | ✅      | 10:20:17  
team-a2    | team-a2-connectivity-probe.team-a2.svc.cluster.local | 80   | ✅      | 10:20:17  
team-b1    | team-b1-connectivity-probe.team-b1.svc.cluster.local | 80   | ✅      | 10:20:17  
team-b2    | team-b2-connectivity-probe.team-b2.svc.cluster.local | 80   | ✅      | 10:20:17  
open-meteo | api.open-meteo.com                                   | 443  | ✅      | 10:20:17  
```

2) Run the Unique UDN demo for team-a

```bash
oc delete -k labs/demo/applications/kustomize/overlays/ovs
oc apply -k labs/demo/applications/kustomize/overlays/udn
```

Verify from team-a and other teams:

```bash
oc exec deploy/team-a2-connectivity-probe -n team-a2 -- curl -s http://localhost:8080
oc exec deploy/team-a1-connectivity-probe -n team-a1 -- curl -s http://localhost:8080
oc exec deploy/team-b1-connectivity-probe -n team-b1 -- curl -s http://localhost:8080
```

Expected example (team-a on its own UDN):

```
Name       | Host                                                 | Port | Status | Checked At
----------------------------------------------------------------------------------------------
team-a1    | team-a1-connectivity-probe.team-a1.svc.cluster.local | 80   | ✅      | 10:20:54  
team-a2    | team-a2-connectivity-probe.team-a2.svc.cluster.local | 80   | ❌      | 10:20:54  
team-b1    | team-b1-connectivity-probe.team-b1.svc.cluster.local | 80   | ❌      | 10:20:54  
open-meteo | api.open-meteo.com                                   | 443  | ✅      | 10:20:54  
```

3) Run the CUDN demo for teams b1 & b2

Verify from B teams:

```bash
oc exec deploy/team-b1-connectivity-probe -n team-b1 -- curl -s http://localhost:8080
oc exec deploy/team-b2-connectivity-probe -n team-b2 -- curl -s http://localhost:8080
```

Expected example (CUDN for B teams):

```
Name       | Host                                                 | Port | Status | Checked At
----------------------------------------------------------------------------------------------
team-a1    | team-a1-connectivity-probe.team-a1.svc.cluster.local | 80   | ❌      | 10:21:27  
team-a2    | team-a2-connectivity-probe.team-a2.svc.cluster.local | 80   | ❌      | 10:21:27  
team-b1    | team-b1-connectivity-probe.team-b1.svc.cluster.local | 80   | ✅      | 10:21:27  
team-b2    | team-b2-connectivity-probe.team-b2.svc.cluster.local | 80   | ✅      | 10:21:27  
open-meteo | api.open-meteo.com                                   | 443  | ✅      | 10:21:27  
```

Cleanup (remove whatever overlay is active)

```bash
oc delete -k labs/demo/applications/kustomize/overlays/udn  || true
```
