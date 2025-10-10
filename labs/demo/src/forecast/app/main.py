from fastapi import FastAPI, HTTPException, Depends, Header
from typing import List, Optional
from . import db
from .models import ForecastIn, Forecast, ForecastUpdate
import MySQLdb.cursors
import os, jwt

JWT_SECRET = os.getenv("JWT_SECRET", "dev-change-me")
app = FastAPI(title="api-forecast")

def require_auth(authorization: Optional[str] = Header(None)):
    if not authorization or not authorization.startswith("Bearer "):
        raise HTTPException(status_code=401, detail="missing token")
    token = authorization.split(" ",1)[1]
    try:
        jwt.decode(token, JWT_SECRET, algorithms=["HS256"])
    except Exception:
        raise HTTPException(status_code=401, detail="invalid token")
    return True

@app.get("/healthz")
def health():
    return {"ok": True}

@app.post("/forecast", response_model=Forecast, status_code=201, dependencies=[Depends(require_auth)])
def create_forecast(f: ForecastIn):
    conn = db.get_conn()
    cur = conn.cursor()
    cur.execute(
        """
        INSERT INTO forecast (city, latitude, longitude, temperature_c, windspeed_kph, observed_at)
        VALUES (%s,%s,%s,%s,%s,%s)
        """,
        (f.city, f.latitude, f.longitude, f.temperature_c, f.windspeed_kph, f.observed_at)
    )
    conn.commit()
    _id = cur.lastrowid
    cur.execute("SELECT id, city, latitude, longitude, temperature_c, windspeed_kph, observed_at, created_at FROM forecast WHERE id=%s", (_id,))
    row = cur.fetchone()
    if not row:
        raise HTTPException(500, "insert failed")
    return Forecast(
        id=row[0], city=row[1], latitude=float(row[2]), longitude=float(row[3]),
        temperature_c=row[4], windspeed_kph=row[5], observed_at=row[6], created_at=row[7]
    )

@app.get("/forecast", response_model=List[Forecast])
def list_recent(limit: int = 20):
    if limit > 100:
        limit = 100
    conn = db.get_conn()
    cur = conn.cursor()
    cur.execute(
        """
        SELECT id, city, latitude, longitude, temperature_c, windspeed_kph, observed_at, created_at
        FROM forecast ORDER BY observed_at DESC LIMIT %s
        """,
        (limit,)
    )
    out = []
    for row in cur.fetchall():
        out.append(
            Forecast(
                id=row[0], city=row[1], latitude=float(row[2]), longitude=float(row[3]),
                temperature_c=row[4], windspeed_kph=row[5], observed_at=row[6], created_at=row[7]
            )
        )
    return out

@app.get("/forecast/{forecast_id}", response_model=Forecast)
def get_one(forecast_id: int):
    conn = db.get_conn()
    cur = conn.cursor()
    cur.execute("""
        SELECT id, city, latitude, longitude, temperature_c, windspeed_kph, observed_at, created_at
        FROM forecast WHERE id=%s
    """, (forecast_id,))
    row = cur.fetchone()
    if not row:
        raise HTTPException(404, "not found")
    return Forecast(
        id=row[0], city=row[1], latitude=float(row[2]), longitude=float(row[3]),
        temperature_c=row[4], windspeed_kph=row[5], observed_at=row[6], created_at=row[7]
    )

@app.put("/forecast/{forecast_id}", response_model=Forecast, dependencies=[Depends(require_auth)])
def update_forecast(forecast_id: int, patch: ForecastUpdate):
    # Build dynamic update
    fields = []
    values = []
    mapping = {
        'city': 'city',
        'latitude': 'latitude',
        'longitude': 'longitude',
        'temperature_c': 'temperature_c',
        'windspeed_kph': 'windspeed_kph',
        'observed_at': 'observed_at'
    }
    for attr, col in mapping.items():
        val = getattr(patch, attr)
        if val is not None:
            fields.append(f"{col}=%s")
            values.append(val)
    if not fields:
        return get_one(forecast_id)
    values.append(forecast_id)
    conn = db.get_conn()
    cur = conn.cursor()
    cur.execute(f"UPDATE forecast SET {', '.join(fields)} WHERE id=%s", tuple(values))
    conn.commit()
    return get_one(forecast_id)

@app.delete("/forecast/{forecast_id}", status_code=204, dependencies=[Depends(require_auth)])
def delete_forecast(forecast_id: int):
    conn = db.get_conn()
    cur = conn.cursor()
    cur.execute("DELETE FROM forecast WHERE id=%s", (forecast_id,))
    conn.commit()
    if cur.rowcount == 0:
        raise HTTPException(404, "not found")
    return None
