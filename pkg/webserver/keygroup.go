package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

func postKeygroup(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := h.HandleCreateKeygroup(fred.Keygroup{
			Name: fred.KeygroupName(kgname),
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
	}
}

func deleteKeygroup(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := h.HandleDeleteKeygroup(fred.Keygroup{
			Name: fred.KeygroupName(kgname),
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
	}
}
