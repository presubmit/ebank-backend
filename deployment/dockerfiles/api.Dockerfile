# builder image
FROM golang:alpine as builder

WORKDIR /ebank

COPY go.sum .
COPY go.mod .
COPY vendor vendor
COPY services/api services/api/
COPY shared shared
COPY pb pb

RUN CGO_ENABLED=0 GOOS=linux go build -o api services/api/cmd/main.go

# final image
FROM alpine
COPY --from=builder /ebank/api /api
CMD ["/api"] 