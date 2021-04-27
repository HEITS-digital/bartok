ARG GO_VERSION=1.15
ARG ALPINE_VERSION=3.13

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

RUN apk add --update git

WORKDIR /app

COPY  . .

RUN go get -u -v github.com/slack-go/slack

RUN go build -o main .

CMD ["./main"]