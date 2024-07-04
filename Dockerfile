# syntax=docker/dockerfile:1
FROM golang:1.22

WORKDIR /app

COPY go.mod *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /demoservice

EXPOSE 8080

CMD ["/demoservice"]