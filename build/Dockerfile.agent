FROM python:3.11-slim

WORKDIR /app

# Copy requirements and install dependencies
COPY ./agent/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY ./agent/app.py .
COPY ./agent/system.prompt .

# Expose Streamlit port
EXPOSE 8501

ENTRYPOINT ["streamlit", "run", "app.py", "--server.port=8501", "--server.address=0.0.0.0"]
