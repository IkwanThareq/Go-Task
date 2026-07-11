# Stage 1 — build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gotask-api .

# Stage 2 — tiny final image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/gotask-api .

EXPOSE 8080

CMD ["./gotask-api"]
