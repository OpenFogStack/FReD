package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

func getKeygroupReplica(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		r, err := h.HandleGetKeygroupReplica(keygroup.Keygroup{
			Name: commons.KeygroupName(kgname),
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.JSON(http.StatusOK, r)
		return
	}
}

func postKeygroupReplica(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		nodeid := context.Params.ByName("nodeid")

		err := h.HandleAddReplica(keygroup.Keygroup{
			Name: commons.KeygroupName(kgname),
		}, replication.Node{
			ID: replication.ID(nodeid),
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}

func deleteKeygroupReplica(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		nodeid := context.Params.ByName("nodeid")

		err := h.HandleRemoveReplica(keygroup.Keygroup{
			Name: commons.KeygroupName(kgname),
		}, replication.Node{
			ID: replication.ID(nodeid),
		})

		if err != nil {
			_ = abort(context, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}
