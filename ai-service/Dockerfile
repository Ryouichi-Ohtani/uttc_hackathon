FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Rename proto directory to avoid conflict with Google's proto package
RUN mv proto protos

# Generate proto files in protos directory (skip if fails - will use mock mode)
RUN python -m grpc_tools.protoc \
    -I. \
    --python_out=. \
    --grpc_python_out=. \
    protos/product_analysis.proto || echo "Proto generation skipped"

EXPOSE 50051

CMD ["python", "-m", "app.grpc_server.server"]
