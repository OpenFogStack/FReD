package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/errors"
)

func abort(c *gin.Context, err error) error {
	d := struct {
		Error string `json:"error"`
	}{
		err.Error(),
	}

	if err, ok := err.(*errors.Error); ok {
		c.JSON(err.Code, d)
		return c.Error(err)
	}

	c.JSON(http.StatusNotFound, d)
	return c.Error(err)
}
