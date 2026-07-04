# Build stage
FROM golang:1.26-alpine AS builder
# CACHEBUST: 2026-07-05-v2
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server/

# Run stage
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/public ./public
RUN mkdir -p uploads/audio uploads/covers data logs
EXPOSE 3001
CMD ["./server"]
