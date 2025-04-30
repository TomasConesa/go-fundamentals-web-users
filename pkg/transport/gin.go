package transport

import (
	"github.com/gin-gonic/gin"
)

func GinServer(
	endpoint Endpoint,
	decode func(c *gin.Context) (any, error),
	encode func(c *gin.Context, response any),
	encodeError func(c *gin.Context, err error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := decode(c)
		if err != nil {
			encodeError(c, err)
			return
		}

		response, err := endpoint(c.Request.Context(), data)
		if err != nil {
			encodeError(c, err)
			return
		}

		encode(c, response)
	}
}
