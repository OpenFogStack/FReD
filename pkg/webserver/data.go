package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

func getItem(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := fred.KeygroupName(context.Params.ByName("kgname"))

		id := context.Params.ByName("id")

		d, err := h.HandleRead(fred.Item{
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
			d.Val,
			string(d.Keygroup),
		}

		context.JSON(http.StatusOK, r)
	}
}

func putItem(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := fred.KeygroupName(context.Params.ByName("kgname"))

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
		err := h.HandleUpdate(fred.Item{
			Keygroup: kgname,
			ID:       id,
			Val:      arg,
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
	}
}

func deleteItem(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := fred.KeygroupName(context.Params.ByName("kgname"))

		id := context.Params.ByName("id")

		err := h.HandleDelete(fred.Item{
			Keygroup: kgname,
			ID:       id,
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
	}

}
