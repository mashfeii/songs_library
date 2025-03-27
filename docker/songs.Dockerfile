FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/songs_library ./cmd/songs_library/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/songs_library .
COPY sql/migrations/ ./sql/migrations/

EXPOSE 8080
CMD ["./songs_library"]
