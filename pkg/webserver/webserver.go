package webserver

import (
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
)

const apiversion string = "/v0"

func postKeygroup(h handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := h.HandleCreateKeygroup(keygroup.Keygroup{
			Name: kgname,
		})

		if err != nil {
			context.AbortWithError(http.StatusConflict, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}

func deleteKeygroup(h handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := h.HandleDeleteKeygroup(keygroup.Keygroup{
			Name: kgname,
		})

		if err != nil {
			context.AbortWithError(http.StatusNotFound, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}

func getItem(h handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		id := context.Params.ByName("id")

		data, err := h.HandleRead(data.Item{
			Keygroup: kgname,
			ID:       id,
		})

		if err != nil {
			context.AbortWithError(http.StatusNotFound, err)
			return
		}

		context.JSON(http.StatusOK, data)
		return
	}
}

func putItem(h handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		id := context.Params.ByName("id")

		var jsonstruct struct {
			Data string `json:"data" binding:"required"`
		}

		if err := context.ShouldBindJSON(&jsonstruct); err != nil {
			log.Print(err)
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		arg := jsonstruct.Data
		err := h.HandleUpdate(data.Item{
			Keygroup: kgname,
			ID:       id,
			Data:     arg,
		})

		if err != nil {
			context.AbortWithError(http.StatusConflict, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}

func deleteItem(h handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		id := context.Params.ByName("id")

		err := h.HandleDelete(data.Item{
			Keygroup: kgname,
			ID:       id,
		})

		if err != nil {
			context.AbortWithError(http.StatusNotFound, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}

}

// Setup sets up a web server client interface for the Fred node.
func Setup(addr string, h handler) error {
	gin.SetMode("release")
	r := gin.New()

	r.Use(logger.SetLogger(logger.Config{
		Logger: &log.Logger,
		UTC:    true,
	}))

	r.POST(apiversion+"/keygroup/:kgname", postKeygroup(h))
	r.DELETE(apiversion+"/keygroup/:kgname", deleteKeygroup(h))

	r.GET(apiversion+"/keygroup/:kgname/data/:id", getItem(h))
	r.PUT(apiversion+"/keygroup/:kgname/data/:id", putItem(h))
	r.DELETE(apiversion+"/keygroup/:kgname/data/:id", deleteItem(h))

	return r.Run(addr)
}
