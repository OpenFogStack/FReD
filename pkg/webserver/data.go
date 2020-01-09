package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
)

func getItem(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := commons.KeygroupName(context.Params.ByName("kgname"))

		id := context.Params.ByName("id")

		d, err := h.HandleRead(data.Item{
			Keygroup: kgname,
			ID:       id,
		})

		if err != nil {
			_ = context.AbortWithError(http.StatusNotFound, err)
			return
		}

		context.JSON(http.StatusOK, d)
		return
	}
}

func putItem(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := commons.KeygroupName(context.Params.ByName("kgname"))

		id := context.Params.ByName("id")

		var jsonstruct struct {
			Data string `json:"data" binding:"required"`
		}

		if err := context.ShouldBindJSON(&jsonstruct); err != nil {
			log.Err(err).Msg("could not bind json")
			_ = context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		arg := jsonstruct.Data
		err := h.HandleUpdate(data.Item{
			Keygroup: kgname,
			ID:       id,
			Data:     arg,
		})

		if err != nil {
			_ = context.AbortWithError(http.StatusConflict, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}

func deleteItem(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := commons.KeygroupName(context.Params.ByName("kgname"))

		id := context.Params.ByName("id")

		err := h.HandleDelete(data.Item{
			Keygroup: kgname,
			ID:       id,
		})

		if err != nil {
			_ = context.AbortWithError(http.StatusNotFound, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}

}
