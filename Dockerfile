FROM golang:1.21 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./src/cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /src/app .

EXPOSE 3000

CMD ["./app"]
