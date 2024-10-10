# builder image
FROM golang:alpine as builder

WORKDIR /ebank

COPY go.sum .
COPY go.mod .
COPY vendor vendor
COPY services/card services/card/
COPY shared shared
COPY pb pb

RUN CGO_ENABLED=0 GOOS=linux go build -o card services/card/cmd/main.go

# final image
FROM alpine

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
  wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe

COPY --from=builder /ebank/card /card
CMD ["/card"] 