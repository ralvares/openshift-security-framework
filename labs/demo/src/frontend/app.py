from flask import Flask, request, jsonify, render_template_string
import os, requests

GATEWAY = os.getenv("GATEWAY_URL", "http://api-gateway.ns-api-user.svc.cluster.local:8080")
app = Flask(__name__)

tmpl = """
<!doctype html>
<title>Weather Dashboard</title>
<h1>Weather Lookup</h1>
<form method=get action="/">
  <input name=city placeholder="City" value="{{city}}"/>
  <button>Fetch</button>
</form>
{% if data %}
<pre>{{data|tojson(indent=2)}}</pre>
{% endif %}
<hr>
<h2>Recent Forecasts</h2>
<pre>{{ forecasts|tojson(indent=2) }}</pre>
"""

@app.route("/")
def index():
    city = request.args.get("city", "Dubai")
    data = None
    try:
        r = requests.get(f"{GATEWAY}/weather", params={"city": city}, timeout=5)
        if r.ok:
            data = r.json()
        else:
            data = {"error": f"gateway status {r.status_code}"}
    except Exception as e:
        data = {"error": str(e)}
    try:
        fr = requests.get(f"{GATEWAY}/forecast", timeout=5)
        forecasts = fr.json() if fr.ok else []
    except Exception:
        forecasts = []
    return render_template_string(tmpl, city=city, data=data, forecasts=forecasts)

@app.route("/healthz")
def health():
    return jsonify({"ok": True})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
