#!/usr/bin/env python3
import socket, time, os, json
from urllib.parse import parse_qs

CONFIG_PATH = "/etc/probe/config.json"
CHECK_TIMEOUT = 0.5  # Reduced timeout for quick network policy checks
CACHE_TTL = 4  # seconds
CACHE_FILE = "/tmp/probe_status.json"

def build_host(service, namespace=None):
    if namespace:
        return f"{service}.{namespace}.svc.cluster.local"
    return service  # Assume external host if no namespace is provided

def check_connection(host, port):
    try:
        with socket.create_connection((host, port), timeout=CHECK_TIMEOUT):
            return True
    except:
        return False

def load_config():
    try:
        with open(CONFIG_PATH) as f:
            return json.load(f).get("destinations", [])
    except:
        return []

def load_cache():
    if os.path.exists(CACHE_FILE):
        age = time.time() - os.path.getmtime(CACHE_FILE)
        if age < CACHE_TTL:
            with open(CACHE_FILE) as f:
                return json.load(f)
    return {}

def save_cache(data):
    with open(CACHE_FILE, "w") as f:
        json.dump(data, f)

def build_html(status):
    html = """Content-type: text/html\n
<html><head><title>Connection Status</title></head><body>
<h1>Live Connection Status</h1>
<table border="1" cellpadding="6">
<tr><th>Name</th><th>Host</th><th>Port</th><th>Status</th><th>Checked At</th></tr>
"""
    for name, info in status.items():
        color = "green" if info["reachable"] else "red"
        symbol = "✅" if info["reachable"] else "❌"
        html += f"<tr><td>{name}</td><td>{info['host']}</td><td>{info['port']}</td><td style='color:{color}'>{symbol}</td><td>{info['timestamp']}</td></tr>\n"
    html += "</table></body></html>\n"
    return html

def build_ascii_table(status):
    print("Content-type: text/plain\n")
    headers = ["Name", "Host", "Port", "Status", "Checked At"]
    rows = []
    for name, info in status.items():
        symbol = "✅" if info["reachable"] else "❌"
        rows.append([name, info["host"], str(info["port"]), symbol, info["timestamp"]])

    # Calculate max column widths
    col_widths = [max(len(str(row[i])) for row in rows + [headers]) for i in range(len(headers))]

    # Format row helper
    def format_row(row, sep=" | "):
        return sep.join(cell.ljust(col_widths[i]) for i, cell in enumerate(row))

    # Output table
    print(format_row(headers))
    print("-" * (sum(col_widths) + 3 * (len(headers) - 1)))
    for row in rows:
        print(format_row(row))

def main():
    # Detect CGI query
    query = os.environ.get("QUERY_STRING", "")
    params = parse_qs(query)
    format_type = params.get("format", ["text"])[0]  # Default to "text" if no format is provided

    config = load_config()
    now = time.strftime('%H:%M:%S')
    updated = {}

    for dest in config:
        name = dest.get("name")
        namespace = dest.get("namespace")  # Namespace can be None for external hosts
        host = build_host(dest["service"], namespace)
        port = dest.get("port", 8080 if namespace else 443)  # Default to 8080 if namespace is specified, otherwise 443
        reachable = check_connection(host, port)
        updated[name] = {
            "host": host,
            "port": port,
            "reachable": reachable,
            "timestamp": now
        }

    save_cache(updated)

    if format_type == "text":
        build_ascii_table(updated)
    else:
        print(build_html(updated))

if __name__ == "__main__":
    main()