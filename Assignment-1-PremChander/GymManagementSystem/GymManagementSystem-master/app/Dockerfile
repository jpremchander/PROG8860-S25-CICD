FROM python:3.9-slim

WORKDIR /app

# Install system dependencies required for mysqlclient and building wheels
RUN apt-get update && apt-get install -y \
    build-essential \
    default-libmysqlclient-dev \
    pkg-config \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*

COPY . .

RUN pip install --no-cache-dir -r requirements.txt

EXPOSE 5000

CMD ["python", "app.py", "runserver", "--host=0.0.0.0"]
