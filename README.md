### Bartók - the community bot

This is a tiny playful project, with the purpose of showcasing a basic integration with slack events api.

##### Tech stack:

- Golang for bot implementation using [Go slack API](https://github.com/slack-go/slack).
- Docker for developing or running the app locally inside a container.
- [GCP App Engine](https://cloud.google.com/appengine/docs/standard/go/building-app) (standard environment) for web hosting.

##### Creating the bot for your slack workspace:

A token is required in order to connect the bot to the slack API. 
For generating a token and creating a bot for your slack workspace you can find detailed information [here](https://slack.com/intl/en-ro/help/articles/115005265703-Create-a-bot-for-your-workspace).

There are multiple ways of connecting your bot to the slack workspace. One version that works well with the Google App Engine standard plan is the [Events API](https://api.slack.com/start/planning/choosing). 
App Engine instance is stopped when the app is idle and auto started when a request is made to the instance. 
By using the Events API, Slack will send an HTTP request to your dedicated endpoint each time an event of your interest is happening, starting your app only when neeeded. 
In case you're wondering about the performance aspect, it's barely noticeable at this point.

##### Contributing

Clone the repo and create a PR against the main branch and notify the other contributors for review. 
Once the PR passed the review and gets pushed, there's a [workflow available](https://github.com/HEITS-digital/bartok/actions/workflows/deploy.yml) which deploys the new version to the Google App Engine. This is also promoting the deployed version to receive all traffic automatically.

##### Local development
- make sure the `build.env` contains the needed secrets. Do not commit the changes on this file!
- `./bartok.sh build` - builds the docker image
- `./bartok run` - starts the app on `localhost:8000` with hot reload
- `./bartok hot` - starts the app on the host machine with hot reload
  - requires `brew install go@1.17.0`
  - go install github.com/cosmtrek/air@1.27.10
  - you may need to perform the following steps in case of `air` can not be found:
    - `mv $(go env GOPATH)/bin/air ~/.air`
    - or create an alias in `.zshrc` with `alias air='$(go env GOPATH)/bin/air'`

Enjoy yourself!
  
