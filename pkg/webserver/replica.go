package webserver

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

func getReplica(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {

		r, err := h.HandleGetReplica()

		if err != nil {
			_ = context.AbortWithError(http.StatusConflict, err)
			return
		}

		context.JSON(http.StatusOK, r)
		return
	}
}

func postReplica(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {

		var jsonstruct struct {
			Nodes []struct{
				ID string `json:"id" binding:"required"`
				IP string `json:"ip" binding:"required"`
				Port int `json:"port" binding:"required"`
			} `json:"nodes" binding:"required"`
		}

		if err := context.ShouldBindJSON(&jsonstruct); err != nil {
			log.Err(err).Msg("could not bind json")
			_ = context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		n := make([]replication.Node, len(jsonstruct.Nodes))

		for i, node := range jsonstruct.Nodes {
			ip := net.ParseIP(node.IP)

			if ip == nil {
				_ = context.AbortWithError(http.StatusConflict, fmt.Errorf("not a valid IP addr: %s", node.IP))
				return
			}

			n[i] = replication.Node{
				ID:   replication.ID(node.ID),
				IP:   ip,
				Port: node.Port,
			}
		}

		err := h.HandleAddReplica(n)

		if err != nil {
			_ = context.AbortWithError(http.StatusConflict, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}

func deleteReplica(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {

		nodeid := context.Params.ByName("nodeid")

		err := h.HandleRemoveReplica(replication.Node{
			ID: replication.ID(nodeid),
		})

		if err != nil {
			_ = context.AbortWithError(http.StatusNotFound, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}