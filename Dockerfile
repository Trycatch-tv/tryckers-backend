# Stage 1: Build Go binary
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./src/cmd

# Stage 2: Minimal runtime image
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main ./main
COPY --from=builder /app/docs ./docs
RUN mkdir -p ./uploads

EXPOSE 8080

CMD ["./main"]
