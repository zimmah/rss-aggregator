# rss-aggregator

`sqlc generate to generate the db scripts`

Requires goose and postgresql to be installed

---

## Ideas for extending the project

- Support pagination of the endpoints that can return many items
- Support different options for sorting and filtering posts using query parameters
- Classify different types of feeds and posts (e.g. blog, podcast, video, etc.)
- Add a CLI client that uses the API to fetch and display posts, maybe it even allows you to read them in your terminal
- Scrape lists of feeds themselves from a third-party site that aggregates feed URLs
- Add support for other types of feeds (e.g. Atom, JSON, etc.)
- Add integration tests that use the API to create, read, update, and delete feeds and posts
- Add bookmarking or "liking" to posts
- Create a simple web UI that uses your backend API
