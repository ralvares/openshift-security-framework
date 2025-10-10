FROM registry.access.redhat.com/ubi9/python-311
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY app.py ./
EXPOSE 5000
USER 1001
ENV PYTHONUNBUFFERED=1
CMD ["python", "app.py"]
