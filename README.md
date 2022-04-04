# Bitrise Step Analytics

## Endpoints

### /track

This endpoint handles `POST` requests with `json` body with the following structure.

```json
{
  "id": "67688665-c71a-4dc4-abef-c577adb1fb83",
  "timestamp": 1648078534000000,
  "event_name": "step_started",
  "properties": {
    "first_property": "first_value",
    "bool_property": false,
    ...
  }
}
```
The service flattens the json and sends it to a PubSub topic defined by the `PUBSUB_PROJECT` and `PUBSUB_TOPIC`
environment variables when the service starts. Google Cloud service account's json (which has right to publish to PubSub)
has to be provided via the `PUBSUB_CREDENTIALS` environment variable.