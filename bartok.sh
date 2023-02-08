#!/usr/local/bin/zsh

source ~/.zshrc

set -e

export $(cat build.env | grep -v "^#" | xargs)

PROJECT_NAME="$(basename $PWD)"

function validateEnv() {
  if [[ -z "$SLACK_TOKEN" ]]; then
      echo "Must provide SLACK_TOKEN in build.env file." 1>&2
      exit 1
  fi

  if [[ -z "$FIRESTORE_PROJECT" ]]; then
      echo "Must provide FIRESTORE_PROJECT in build.env file." 1>&2
      exit 1
  fi

  if [[ -z "SLACK_SIGNING_SECRET" ]]; then
      echo "Must provide SLACK_SIGNING_SECRET in build.env file." 1>&2
      exit 1
  fi

  if [[ -z "SLACK_VERIFICATION_TOKEN" ]]; then
      echo "Must provide SLACK_VERIFICATION_TOKEN in build.env file." 1>&2
      exit 1
  fi

  if [[ -z "GOOGLE_CLIENT_ID" ]]; then
      echo "Must provide GOOGLE_CLIENT_ID in build.env file." 1>&2
      exit 1
  fi

  if [[ -z "GOOGLE_CLIENT_SECRET" ]]; then
      echo "Must provide GOOGLE_CLIENT_SECRET in build.env file." 1>&2
      exit 1
  fi

  if [[ -z "GOOGLE_REFRESH_TOKEN" ]]; then
      echo "Must provide GOOGLE_REFRESH_TOKEN in build.env file." 1>&2
      exit 1
  fi
}

function hot() {
    validateEnv
    go mod tidy
    air
}


function build() {
    echo "Building: $PROJECT_NAME"
    validateEnv
    docker build \
    -f $PWD/Dockerfile \
    -t $PROJECT_NAME \
    .
}

function run() {
    # docker image rm $PROJECT_NAME
    docker stop $PROJECT_NAME | true
    build
    echo "Running the app..."

    go mod tidy && \
    docker run \
    -e SLACK_TOKEN=$SLACK_TOKEN \
    -e SLACK_VERIFICATION_TOKEN=$SLACK_VERIFICATION_TOKEN \
    -e FIRESTORE_PROJECT=$FIRESTORE_PROJECT \
    -e SLACK_SIGNING_SECRET=$SLACK_SIGNING_SECRET \
    -e BIRTHDAY_CHANNEL_ID=$BIRTHDAY_CHANNEL_ID \
    -e GOOGLE_CALENDAR_EMPLOYEES=$GOOGLE_CALENDAR_EMPLOYEES \
    -e GOOGLE_CLIENT_SECRET=$GOOGLE_CLIENT_SECRET \
    -e GOOGLE_CLIENT_ID=$GOOGLE_CLIENT_ID \
    -e GOOGLE_REFRESH_TOKEN=$GOOGLE_REFRESH_TOKEN \
    -e PORT=$PORT \
    -p 8000:8000 \
    --volume=$(PWD)/internal:/app/internal \
    --volume=$(PWD)/cmd:/app/cmd  \
    -t $PROJECT_NAME
}

"$@"

if [[ $# -eq 0 ]] ; then
echo "

USAGE: ./bartok.sh [option]

where [option] is one of:

- run
- build
- hot

"
fi
