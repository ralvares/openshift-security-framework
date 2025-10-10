#!/usr/bin/python3
import cgi, json, urllib.request, urllib.parse, re, sys, time, os

print("Content-type: application/json\n")
form = cgi.FieldStorage()
city = form.getfirst("city", "").strip()
if not city:
    print(json.dumps({"error": "missing ?city= parameter"})); sys.exit()
if len(city) > 64 or not re.fullmatch(r"[A-Za-z .'-]+", city):
    print(json.dumps({"error": "invalid city format"})); sys.exit()
try:
    geo_url = ("https://geocoding-api.open-meteo.com/v1/search?" + urllib.parse.urlencode({"name": city, "count": 1}))
    with urllib.request.urlopen(geo_url, timeout=5) as resp:
        geo = json.load(resp)
    if not geo.get("results"):
        print(json.dumps({"error": f"city '{city}' not found"})); sys.exit()
    first = geo["results"][0]
    lat, lon = first.get("latitude"), first.get("longitude")
    if lat is None or lon is None:
        print(json.dumps({"error": "missing coordinates"})); sys.exit()
    base_url = "https://api.open-meteo.com/v1/forecast"
    params = {
        "latitude": lat,
        "longitude": lon,
        "current": ",".join([
            "temperature_2m","relative_humidity_2m","apparent_temperature","wind_speed_10m","wind_direction_10m","pressure_msl","cloud_cover","weathercode"
        ]),
        "hourly": ",".join(["temperature_2m","relative_humidity_2m","wind_speed_10m","pressure_msl","precipitation","uv_index"]),
        "daily": ",".join([
            "temperature_2m_max","temperature_2m_min","apparent_temperature_max","apparent_temperature_min","precipitation_sum","sunrise","sunset","shortwave_radiation_sum","wind_gusts_10m_max"
        ]),
        "forecast_days": 7,
        "timezone": "auto"
    }
    wx_url = base_url + "?" + urllib.parse.urlencode(params)
    with urllib.request.urlopen(wx_url, timeout=10) as resp:
        weather = json.load(resp)
    current = weather.get("current", {})
    daily = weather.get("daily", {})
    daily_summary = []
    for i, day in enumerate(daily.get("time", [])):
        daily_summary.append({
            "date": day,
            "temp_max": daily.get("temperature_2m_max", [None]*7)[i],
            "temp_min": daily.get("temperature_2m_min", [None]*7)[i],
            "apparent_max": daily.get("apparent_temperature_max", [None]*7)[i],
            "apparent_min": daily.get("apparent_temperature_min", [None]*7)[i],
            "sunrise": daily.get("sunrise", [None]*7)[i],
            "sunset": daily.get("sunset", [None]*7)[i],
            "wind_gusts_max": daily.get("wind_gusts_10m_max", [None]*7)[i],
            "radiation": daily.get("shortwave_radiation_sum", [None]*7)[i],
            "precipitation": daily.get("precipitation_sum", [None]*7)[i]
        })
    out = {
        "city": city,
        "latitude": lat,
        "longitude": lon,
        "timezone": weather.get("timezone"),
        "current": current,
        "daily_summary": daily_summary
    }
    print(json.dumps(out))
except Exception:
    print(json.dumps({"error": "upstream request failed"}))
# structured log to stderr
log = {
  "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
  "component": "weather-proxy",
  "namespace": os.getenv("POD_NAMESPACE",""),
  "event_type": "ingress",
  "method": "GET",
  "path": "/",
  "status": 200,
  "message": "cgi_complete"
}
print(json.dumps(log), file=sys.stderr)
