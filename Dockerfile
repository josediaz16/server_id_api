FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /go/src/server_id_api

RUN apk add whois

ENV GO111MODULE on
ENV CGO_ENABLED 0

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go mod tidy

EXPOSE 8080

CMD ["server_id_api"]
