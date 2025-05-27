# Use official Python image
FROM python:3.9-slim

# Set working directory inside container
WORKDIR /app

# Copy all app files into the container
COPY . .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Expose Flask default port
EXPOSE 5000

# Run the Flask app
CMD ["python", "app.py"]
