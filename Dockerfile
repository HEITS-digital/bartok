ARG GO_VERSION=1.20
ARG ALPINE_VERSION=3.13

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

RUN apk add --update git

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.27.10
RUN go get -u -v github.com/slack-go/slack

COPY go.mod go.sum .air.toml ./

RUN go mod download

EXPOSE 8000

HEALTHCHECK  --interval=5m --timeout=3s --start-period=2m\
  CMD wget --no-verbose --tries=3 --spider http://localhost:8000/status/ || exit 1

CMD ["air"]
