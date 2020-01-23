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
			_ = abort(context, err)
			return
		}

		/*
			{
			  "id": "hello",
			  "value": "Hello World!",
			  "keygroup": "Test-Keygroup"
			}
		*/

		var r = struct {
			ID       string `json:"id" binding:"required"`
			Value    string `json:"value" binding:"required"`
			Keygroup string `json:"keygroup" binding:"required"`
		}{
			d.ID,
			d.Data,
			string(d.Keygroup),
		}

		context.JSON(http.StatusOK, r)
		return
	}
}

func putItem(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := commons.KeygroupName(context.Params.ByName("kgname"))

		id := context.Params.ByName("id")

		/*
			{
			  "id": "hello",
			  "value": "Hello World!",
			  "keygroup": "Test-Keygroup"
			}
		*/
		var jsonstruct struct {
			ID       string `json:"id" binding:"required"`
			Value    string `json:"value" binding:"required"`
			Keygroup string `json:"keygroup" binding:"required"`
		}

		if err := context.ShouldBindJSON(&jsonstruct); err != nil {
			log.Err(err).Msg("could not bind json")
			_ = abort(context, err)
			return
		}

		arg := jsonstruct.Value
		err := h.HandleUpdate(data.Item{
			Keygroup: kgname,
			ID:       id,
			Data:     arg,
		})

		if err != nil {
			_ = abort(context, err)
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
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}

}
