# Gin Utils

These are utilities for the gin web library.

Some features:

- Improved error handling
  - Write `func (c *gin.Context) error` handlers, returning an error
  - Respond with a common JSON object on any errors
- Complex parameters
  - Parse UUIDs from URL path and return error if invalid
