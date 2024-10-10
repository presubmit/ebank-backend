FROM golang:alpine
RUN apk --no-cache add git
RUN go get "github.com/cespare/reflex"

WORKDIR /ebank

COPY go.sum .
COPY go.mod .
COPY vendor vendor

COPY services services
COPY shared shared
COPY pb pb

ENTRYPOINT [ "/bin/sh", "-c" ]
CMD [ "reflex -R vendor -R pb -s -d none -- sh -c \"go run services/*.go\"" ]