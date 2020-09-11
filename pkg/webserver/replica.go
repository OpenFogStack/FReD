package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

func getAllReplica(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {

		d, err := h.HandleGetAllReplica()

		if err != nil {
			_ = abort(context, err)
			return
		}

		/*
			{
			  "nodes": [
			    "nodeB",
			    "nodeC",
			    "nodeD"
			  ]
			}
		*/

		nodes := make([]string, len(d))

		for i, n := range d {
			nodes[i] = string(n.ID)
		}

		var r = struct {
			Nodes []string `json:"nodes" binding:"required"`
		}{
			Nodes: nodes,
		}

		context.JSON(http.StatusOK, r)
	}
}

func getReplica(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {

		nodeid := context.Params.ByName("nodeid")

		d, err := h.HandleGetReplica(fred.Node{
			ID: fred.NodeID(nodeid),
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		/*
			{
			  "id": "nodeA",
			  "addr": "172.12.0.3",
			  "zmqPort": 5555
			}
		*/

		var r = struct {
			ID      string `json:"id" binding:"required"`
			Addr    string `json:"addr" binding:"required"`
			ZMQPort int    `json:"zmqPort" binding:"required"`
		}{
			string(d.ID),
			d.Addr.Addr,
			d.Port,
		}

		context.JSON(http.StatusOK, r)
	}
}