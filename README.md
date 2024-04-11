# RSS Aggregator

This project aggregates RSS feeds provided by a user. Once a user registers, they can use their API key to:

- [x] Add feeds to the aggregator
- [x] View which feeds they're subscribed to
- [x] View a list of feeds already available to subscribe to
- [x] Subscribe/Unsubscribe from feeds
- [x] Retrieve `n` recent posts from feeds they've subscribed to, where `n` is an optional provided query limit

# Goal

I wanted to create something with a bit more actual use potential while getting more practice with using SQL in Go, and some of the tooling available, like [sqlc](https://docs.sqlc.dev/en/latest/index.html).

# Installation

Clone this repo:

```bash
git clone git@github.com:ellielle/rss-aggregator.git
```

This repo makes use of [godotenv](https://github.com/joho/godotenv) to provide the env variables `PORT` and `DB_CONNECTION`, which is making use of a postgresql database.

[Goose](https://github.com/pressly/goose) is used to handle migrations, and [sqlc](https://docs.sqlc.dev/en/latest/index.html) is used to generate Go code.

Navigate to the `sql/schema` directory and run:

```bash
goose postgres postgres://username:password@localhost:5432/dbname up
```

After migrating, run `sqlc generate` to generate Go code to interact with the database.

Then build and run the server using the following command:

```bash
go build -o out && ./out
```

# Usage

## Post /v1/users - Create User

> It currently only needs a non-unique name to get an API key

Request Body:

```json
{
  "name": "NAME"
}
```

Response Body:

```json
{
  "name": "NAME",
  "api_key": "3b354ee128803150bc541741278134ba2b3d47fecd0767e135d19239b631f385"
}
```

## GET /v1/users - Get User information

Request Headers: Authorization ApiKey APIKEY

Request Response:

```json
{
  "id": "e101afa6-8bdc-4899-bc40-5a50684056b9",
  "name": "NAME",
  "created_at": "2024-04-11T11:41:06.436996Z",
  "updated_at": "2024-04-11T11:41:06.436996Z",
  "api_key": "3b354ee128803150bc541741278134ba2b3d47fecd0767e135d19239b631f385"
}
```

## POST /v1/feeds - Add a feed to the aggregator

Returns a `feed id` and a `feed_follow id`, where the `feed_follow_id` is used for feeds the user is actively following.

> Currently only XML feeds are supported

Request Headers: Authorization ApiKey APIKEY

Request Body:

```json
{
  "name": "FEEDNAME",
  "url": "FEEDURL"
}
```

Request Response:

```json
{
  "feed": "0078a18f-f0d6-4294-abd1-5adf34409483",
  "feed_follow": "435f7482-3add-499e-bc60-1ecfef031187"
}
```

## GET /v1/feeds - Get all available feeds

Request Response:

```json
{
  "feeds": [
    {
      "id": "0078a18f-f0d6-4294-abd1-5adf34409483",
      "created_at": "2024-04-11T13:00:49.953444Z",
      "updated_at": "2024-04-11T13:00:49.953444Z",
      "name": "FEED1",
      "url": "URL1",
      "user_id": "e101afa6-8bdc-4899-bc40-5a50684056b9",
      "last_fetched_at": "0001-01-01T00:00:00Z"
    },
    {
      "id": "06e7a9cb-66fd-48df-878b-438e9442961f",
      "created_at": "2024-04-10T12:20:11.269447Z",
      "updated_at": "2024-04-11T13:01:36.508076Z",
      "name": "FEED2",
      "url": "URL2",
      "user_id": "8939ad8a-2ba3-45a3-96ed-c247db28a9e9",
      "last_fetched_at": "2024-04-11T13:01:36.508076Z"
    }
  ]
}
```

## POST /v1/feed_follows - Subscribes a user to a feed

Request Headers: Authorization ApiKey APIKEY

Request Body:

```json
{
  "feed_id": "06e7a9cb-66fd-48df-878b-438e9442961f"
}
```

Response Body:

```json
{
  "id": "f61da9f3-882d-4ddf-86b8-4955fc8bde62",
  "feed_id": "0078a18f-f0d6-4294-abd1-5adf34409483",
  "user_id": "97cb91d3-feca-4bdd-b9d8-1c9c8a17b579",
  "created_at": "2024-04-11T13:21:59.08347Z",
  "updated_at": "2024-04-11T13:21:59.08347Z"
}
```

## GET /v1/feed_follow - Gets a user's list of followed feeds

Request Headers: Authorization ApiKey APIKEY

Request Response:

```json
{
  "feeds": [
    {
      "id": "f61da9f3-882d-4ddf-86b8-4955fc8bde62",
      "created_at": "2024-04-11T13:21:59.08347Z",
      "updated_at": "2024-04-11T13:21:59.08347Z",
      "user_id": "97cb91d3-feca-4bdd-b9d8-1c9c8a17b579",
      "feed_id": "0078a18f-f0d6-4294-abd1-5adf34409483"
    }
  ]
}
```

## DELETE /v1/feed_follows/{feed_follow_id} - Unsubscribe a user from a feed with specified ID

Request Headers: Authorization ApiKey APIKEY

Request Body:

```json
{
  "message": "OK"
}
```

## GET /v1/posts/{limit} - Get a list of posts of LIMIT length from the user's subscribed feeds

Ordered from newest -> oldest

Request Headers: Authorization ApiKey APIKEY

Request Body:

```json
{
  "posts": [
    {
      "id": "f20e854d-d6c7-43a4-a872-306021db0e71",
      "created_at": "2023-01-08T00:00:00Z",
      "updated_at": "2023-01-08T00:00:00Z",
      "title": "TITLE",
      "url": "URL",
      "description": "TEXT",
      "feed_id": "06e7a9cb-66fd-48df-878b-438e9442961f"
    }
  ]
}
```

## GET /v1/posts - Same as above but with a default limit

## GET /v1/readiness

Endpoints to test that the server is ready and that errors are thrown correctly.

Response Body:

```json
{
  "status": "ok"
}
```

## GET /v1/err

Response Body:

```json
{
  "error": "Internal Server Error"
}
```
