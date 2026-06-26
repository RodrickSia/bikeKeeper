# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bikekeeper ./cmd/bikeKeeper

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
RUN mkdir -p /app/images && chown appuser:appgroup /app/images
COPY --from=builder /app/bikekeeper .

USER appuser

EXPOSE 8080

CMD ["./bikekeeper"]