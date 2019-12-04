# TO BUILD 
# go build -ldflags '-linkmode external -w -extldflags "-static"'

FROM golang:latest

WORKDIR /go/src/

COPY ./vault-secret /go/src/

EXPOSE 8080

ENTRYPOINT [ "./vault-secret" ]
