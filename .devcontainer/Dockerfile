ARG GO_VERSION=1.15
ARG ALPINE_VERSION=3.13

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

RUN apk update && apk add bash
RUN apk add build-base

# Install packages and Go language server
RUN apk add -q --update --progress --no-cache git sudo openssh-client zsh
# RUN go get -u -v golang.org/x/tools/gopls 2>&1

RUN go get -x -d github.com/stamblerre/gocode 2>&1 \
    && go build -o gocode-gomod github.com/stamblerre/gocode \
    && mv gocode-gomod $GOPATH/bin/ \
    # # Install other tools.
    && go get -u -v \
        github.com/mdempsky/gocode \
        github.com/uudashr/gopkgs/cmd/gopkgs \
        github.com/ramya-rao-a/go-outline \
        github.com/acroca/go-symbols \
        github.com/go-delve/delve/cmd/dlv \
        github.com/stamblerre/gocode \
        github.com/rogpeppe/godef \
        golang.org/x/lint/golint 2>&1

RUN go get -u -v github.com/slack-go/slack

CMD ["./main.go"]