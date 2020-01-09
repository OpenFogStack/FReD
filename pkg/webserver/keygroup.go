package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
)

func postKeygroup(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := h.HandleCreateKeygroup(keygroup.Keygroup{
			Name: commons.KeygroupName(kgname),
		})

		if err != nil {
			_ = context.AbortWithError(http.StatusConflict, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}

func deleteKeygroup(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := h.HandleDeleteKeygroup(keygroup.Keygroup{
			Name: commons.KeygroupName(kgname),
		})

		if err != nil {
			_ = context.AbortWithError(http.StatusNotFound, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}
