package ginutils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetPathUUID retrieves a UUID from the specified path parameter in the Gin context.
//
// This returns an error in these cases:
//   - The path parameter is missing. This is considered a BadRequest error.
//   - The path parameter is not a valid UUID. This is also considered a BadRequest error.
func GetPathUUID(c *gin.Context, param string) (uuid.UUID, error) {
	if s := c.Param(param); s == "" {
		return uuid.Nil, BadRequest(fmt.Sprintf("missing path parameter %q", param))
	} else if u, err := uuid.Parse(s); err != nil {
		return uuid.Nil, BadRequest(fmt.Sprintf("invalid path parameter %q: %v", param, err))
	} else {
		return u, nil
	}
}
