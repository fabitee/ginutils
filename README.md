# Gin Utils

These are utilities for the gin web library.

Some features:

- Improved and unified error handling
  - Write `func (c *gin.Context) error` handlers, returning an error
  - Use `ginutils.Recovery` middleware to respond with `ErrorResponse` objects
  - Respond with a common JSON object on any errors
- Complex parameters
  - Parse UUIDs from URL path and return error if invalid

## Error Handling

One central idea is to use JSON objects for every error response, even on panics.

Clients should expect to always receive an `ErrorResponse` object, which looks like this:

```json
{
  "status": 400,
  "message": "Missing parameter 'id' in URL"
}
```
