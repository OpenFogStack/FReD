package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

func postSeed(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		var jsonstruct struct {
			ID   string `json:"id" binding:"required"`
			Addr string `json:"addr" binding:"required"`
		}

		if err := context.ShouldBindJSON(&jsonstruct); err != nil {
			log.Err(err).Msg("could not bind json")
			_ = abort(context, err)
			return
		}

		addr, err := replication.ParseAddress(jsonstruct.Addr)

		if err != nil {
			_ = abort(context, err)
			return
		}

		err = h.HandleSeed(replication.Node{
			ID:   replication.ID(jsonstruct.ID),
			Addr: addr,
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}
