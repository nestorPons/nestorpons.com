# ── Stage 1: compilar ────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY *.go ./
COPY components/ ./components/

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ── Stage 2: imagen final (mínima) ───────────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 80

CMD ["./server"]