package ginutil

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPathUUID(c *gin.Context, param string) (uuid.UUID, error) {
	if s := c.Param(param); s == "" {
		return uuid.Nil, BadRequest(fmt.Sprintf("missing path parameter %q", param))
	} else if u, err := uuid.Parse(s); err != nil {
		return uuid.Nil, BadRequest(fmt.Sprintf("invalid path parameter %q: %v", param, err))
	} else {
		return u, nil
	}
}
