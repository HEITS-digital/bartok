#!/usr/bin/env bash

set -e

BASEDIR=$(dirname "$0")
PROJECT_NAME="$(basename $PWD)"


function build() {
    echo "Building: $PROJECT_NAME"
    docker build \
    -f $PWD/Dockerfile \
    -t $PROJECT_NAME \
    $BASEDIR
}

function run() {
    # docker image rm $PROJECT_NAME
    docker stop $PROJECT_NAME | true
    build
    echo "Running the app..."
    docker run --name $PROJECT_NAME -e SLACK_TOKEN=$SLACK_TOKEN $PROJECT_NAME
}

"$@"

if [[ $# -eq 0 ]] ; then
echo "

USAGE: ./bartok.sh [option]

where [option] is one of:

- run
- build

"
fi