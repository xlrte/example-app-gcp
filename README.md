# How to run this
* Follow the first 3 steps, [starting here](https://xlrte.dev/docs/getting-started/setup-gcp) of the `xlrte` setup guide. Initialize the xlrte environment `dev`.
* Run `docker build . -t gcr.io/[gcp project selected in the steps above]/hello-app:v1`
* Run `docker push gcr.io/[gcp project selected in the steps above]/hello-app:v1`
* Run `xlrte apply -e dev`

In about 15 minutes or so (depending on CloudSQL development time), you should have an environment up on Cloud Run with the following endpoints:
* `/database` pings the Postgres database to see it can be connected to.
* `POST /bucket` creates a simple text file in the configured bucket
* `GET /bucket` shows the text of the above.
* `/publish` publishes a Pub/Sub message that then gets `POST`:ed to `/` (can be seen in logs of Cloud Run.)