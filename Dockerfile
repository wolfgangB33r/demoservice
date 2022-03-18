FROM golang:1.17.5-alpine AS builder
RUN apk add --no-cache ca-certificates git
RUN apk add build-base

RUN export GO111MODULE=on
WORKDIR /src
COPY favicon.ico .

COPY ./service.go .
# restore dependencies
COPY go.mod go.sum ./
RUN go mod download \
    && go mod tidy


RUN go build -ldflags="-linkmode external" -o /demoservice .

FROM alpine AS release
RUN apk add --no-cache ca-certificates
# RUN GRPC_HEALTH_PROBE_VERSION=v0.4.6 && \
#    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
#    chmod +x /bin/grpc_health_probe
WORKDIR /src
COPY --from=builder /demoservice ./demoservice

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

EXPOSE 8080
ENTRYPOINT ["/src/demoservice"]