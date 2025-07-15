FROM golang:1.24.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main ./src/cmd

EXPOSE 8080

CMD ["sh", "-c", "./main"]
