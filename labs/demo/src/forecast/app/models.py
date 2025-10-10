from pydantic import BaseModel, Field
from datetime import datetime
from typing import Optional

class ForecastIn(BaseModel):
    city: str = Field(max_length=64)
    latitude: float
    longitude: float
    temperature_c: Optional[float]
    windspeed_kph: Optional[float]
    observed_at: datetime

class Forecast(ForecastIn):
    id: int
    created_at: datetime

class ForecastUpdate(BaseModel):
    city: Optional[str] = Field(default=None, max_length=64)
    latitude: Optional[float] = None
    longitude: Optional[float] = None
    temperature_c: Optional[float] = None
    windspeed_kph: Optional[float] = None
    observed_at: Optional[datetime] = None
