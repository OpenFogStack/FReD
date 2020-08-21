package webserver

import (
	"github.com/rs/zerolog/log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

func getKeygroupTrigger(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		nodes, err := h.HandleGetKeygroupTriggers(fred.Keygroup{
			Name: fred.KeygroupName(kgname),
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		/*
			{
			  "nodes": [
			    {
					id: "triggernodeB",
					host: "172.26.0.1:3333"
				},
			    {
					id: "triggernodeC",
					host: "172.26.0.10:3333"
				},
			    {
					id: "triggernodeD",
					host: "172.22.2.1:3333"
				},
			  ]
			}
		*/

		var r = struct {
			Nodes []struct {
				ID   string `json:"id" binding:"required"`
				Host string `json:"host" binding:"required"`
			} `json:"nodes" binding:"required"`
		}{}

		for _, node := range nodes {
			r.Nodes = append(r.Nodes, struct {
				ID   string `json:"id" binding:"required"`
				Host string `json:"host" binding:"required"`
			}{ID: node.ID, Host: node.Host})
		}

		context.JSON(http.StatusOK, r)
	}
}

func postKeygroupTrigger(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		nodeid := context.Params.ByName("triggernodeid")

		/*
			{
			  "host": "172.260.0.1:3333"
			}
		*/

		var jsonstruct struct {
			Host string `json:"host" binding:"required"`
		}

		if err := context.ShouldBindJSON(&jsonstruct); err != nil {
			log.Err(err).Msg("could not bind json")
			_ = abort(context, err)
			return
		}

		if _, _, err := net.SplitHostPort(jsonstruct.Host); err != nil {
			log.Err(err).Msg("not a valid address")
			_ = abort(context, err)
			return
		}

		err := h.HandleAddTriggers(fred.Keygroup{
			Name: fred.KeygroupName(kgname),
		}, fred.Trigger{
			ID:   nodeid,
			Host: jsonstruct.Host,
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
	}
}

func deleteKeygroupTrigger(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		nodeid := context.Params.ByName("triggernodeid")

		err := h.HandleRemoveTrigger(fred.Keygroup{
			Name: fred.KeygroupName(kgname),
		}, fred.Trigger{
			ID: nodeid,
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
	}
}
