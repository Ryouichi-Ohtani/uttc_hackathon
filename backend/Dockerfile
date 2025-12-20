# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy proto files to path expected by go.mod replace (../proto resolves to /proto)
COPY proto/ /proto/

# Copy backend files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
