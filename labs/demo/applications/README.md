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
oc exec deploy/connectivity-probe -n demo-app -- curl -s http://localhost:8080/cgi-bin/index.py?format=text
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

## Configuration File: `config.json`

The `config.json` file is used to define the destinations that the HTTPD server interacts with. Below is a breakdown of its structure:

### Structure
The file contains a JSON object with the following key:
- `destinations`: An array of destination objects. Each object specifies the details of a service.

### Destination Object
Each destination object has the following properties:
- `name`: The name of the destination.
- `service`: The name of the service.
- `namespace`: The namespace where the service resides.
- `port`: The port the service listens on.

### Example
```json
{
  "destinations": [
    {
      "name": "team-b1",
      "service": "httpd",
      "namespace": "team-b1",
      "port": 8080
    },
    {
      "name": "weather-proxy",
      "service": "weather-proxy",
      "namespace": "weather-proxy",
      "port": 80
    }
  ]
}
```

### Explanation
- **team-b1**: Represents a service named `httpd` in the `team-b1` namespace, listening on port `8080`.
- **weather-proxy**: Represents a service named `weather-proxy` in the `weather-proxy` namespace, listening on port `80`.

This configuration allows the HTTPD server to route traffic to the specified destinations.

## Notes
- Ensure that the `connectivity-probe` deployment is running in the `demo-app` namespace.
- The application is designed to be used in OpenShift environments.