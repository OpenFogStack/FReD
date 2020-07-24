package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func abort(c *gin.Context, err error) error {
	d := struct {
		Error string `json:"error"`
	}{
		err.Error(),
	}

	c.JSON(http.StatusNotFound, d)
	return c.Error(err)
}
