import os
import threading
import MySQLdb

_lock = threading.Lock()
_conn = None

def get_conn():
    global _conn
    if _conn is None:
        with _lock:
            if _conn is None:
                _conn = MySQLdb.connect(
                    host=os.getenv("MYSQL_HOST","mysql.ns-data-user.svc.cluster.local"),
                    user=os.getenv("MYSQL_USER","forecast"),
                    passwd=os.getenv("MYSQL_PASSWORD","changeme"),
                    db=os.getenv("MYSQL_DATABASE","forecasts"),
                    connect_timeout=3,
                    charset="utf8mb4"
                )
    return _conn
