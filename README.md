# xlrte example application for GCP

This app is the simplest possible app that exercises all the services of GCP that [xlrte](https://xlrte.dev) currently supports.
It is written in Go, is not intended to be "pretty", but simply demonstrate ability to connect and interact with the different services. Anyone familiar with any language should hopefully be able to roughly understand what the code in `cmd/main.go` does.

## xlrte config

The `xlrte` config lives in `.xlrte/config` and contains a single service in the `services` sub-folder that users `CloudRun` as a runtime and also uses `CloudSQL`, `Pub/Sub` and `Cloud Storage`. There is a single environment called `dev`.

## Setup

#### Prepare pre-requisites

- [Setup GCP for xlrte](https://xlrte.dev/docs/getting-started/setup-gcp)
- [Install xlrte](https://xlrte.dev/docs/getting-started/install)
- Decide on a new GCP project name for this exercise, let's say `europe-west6` for the region and `my-cool-dev-env` for the GCP project for the purposes of demonstration (the GCP project needs to be unique across all GCP projects, so please change it).

#### Build the example service and push it

- Run `gcloud config set project my-cool-dev-env` (change the project as necessaru)
- Run `gcloud auth configure-docker` so that you can push images.
- Run `docker build . -t gcr.io/my-cool-dev-env/hello-app:v1` (change the project as necessary)
- Run `docker push gcr.io/my-cool-dev-env/hello-app:v1` (change the project as necessary)

#### Deploy with xlrte

- Initialize the environment, let's call it dev, with `xlrte init -e dev -p gcp -c my-cool-dev-env -r europe-west6`
- Run `xlrte apply -e dev` to create and deploy the environment.

**Warning! Initial CloudSQL setup takes a long time!**

In about 10-15 minutes or so (depending on CloudSQL development time), you should have an environment up on Cloud Run with the following endpoints:

- `/database` pings the Postgres database to see it can be connected to.
- `POST /bucket` creates a simple text file in the configured bucket
- `GET /bucket` shows the text of the above.
- `/publish` publishes a Pub/Sub message that then gets `POST`:ed to `/` (can be seen in logs of Cloud Run).
