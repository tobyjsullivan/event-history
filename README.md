# Event History

This service returns the complete event histories for requested event
IDs.

This service is meant for use in conjunction with the `event-log-reader`
service which will provide the current head event ID for a given log.

## Running Locally with Docker Compose

TK

## API

### GET /histories/:headEventId

Redirects the request to a file containing the complete history for the
 specified eventId.

#### History File Format

Histories are served as a single large file of newline delimited JSON
objects. Each line in the file will be an event object ordered from
oldest to newest.

The file as a whole is *not* a valid JSON object.
This decision was made to enable efficient line-by-line processing of
the file as it is transmitted. This allows histories of arbitrary size
to be processed without needing to store the full history in memory.

#### Event JSON Object

Each event JSON object has the following properties:

|Property| Type| Description|
|---|---|---|
|id| Hex-encoded 64-byte array | A 64-byte unique identifier|
|type| string| A type identifier. Values are specific to your application but are required for all events. |
|data| Base64-encoded byte array| Variable-length data of the event.|

Example:

```json
{"id":"3683dc3c2e1e22d3068fe0dd779aa21ec9d5ea57b1db946a5bbc311689bbafb7","type":"AccountOpened","data":"eyJhY2NvdW50SWQiOiJhODg0Mzc2NS1mZmI4LTQ0MTctYjVjNi0wYTI3N2FjY2ZkY2QifQ=="}
```
