FROM golang:alpine AS build

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /go/src/server_id_api

ENV GO111MODULE on
ENV CGO_ENABLED 0

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go mod tidy

RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN apk add whois bash

WORKDIR /usr/bin
COPY --from=build /go/src/server_id_api/bin /go/bin
EXPOSE 3000
CMD ["/go/bin/api"]
