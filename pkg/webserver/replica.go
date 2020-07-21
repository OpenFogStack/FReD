package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

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

func postReplica(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {

		/*
			{
			  "nodes": [
			    {
			      "id": "nodeB",
			      "addr": "172.12.0.3",
			      "zmqPort": 5555
			    },
			    {
			      "id": "nodeC",
			      "addr": "nodeC.nodes.mcc-f.red",
			      "zmqPort": 5554
			    },
			    {
			      "id": "nodeD",
			      "addr": "localhost",
			      "zmqPort": 5553
			    }
			  ]
			}
		*/

		type node struct {
			ID   string `json:"id" binding:"required"`
			Addr string `json:"addr" binding:"required"`
			Port int    `json:"zmqPort" binding:"required"`
		}

		var jsonstruct struct {
			Nodes []node `json:"nodes" binding:"required"`
		}

		if err := context.ShouldBindJSON(&jsonstruct); err != nil {
			log.Err(err).Msg("could not bind json")
			_ = abort(context, err)
			return
		}

		n := make([]fred.Node, len(jsonstruct.Nodes))

		for i, node := range jsonstruct.Nodes {
			addr, err := fred.ParseAddress(node.Addr)

			if err != nil {
				_ = abort(context, err)
				return
			}

			n[i] = fred.Node{
				ID:   fred.NodeID(node.ID),
				Addr: addr,
				Port: node.Port,
			}
		}

		err := h.HandleAddNode(n)

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
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

func deleteReplica(h fred.ExtHandler) func(context *gin.Context) {
	return func(context *gin.Context) {

		nodeid := context.Params.ByName("nodeid")

		err := h.HandleRemoveNode(fred.Node{
			ID: fred.NodeID(nodeid),
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
	}
}
