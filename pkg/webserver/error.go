package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/errors"
)

func abort(c *gin.Context, err error) error {
	if err, ok := err.(*errors.Error); ok {
		c.String(err.Code, err.Error())
		return c.Error(err)
	}

	c.String(http.StatusNotFound, err.Error())
	return c.Error(err)
}