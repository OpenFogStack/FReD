package webserver

import (


	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
)

const apiversion string = "/v0"

// Setup sets up a web server client interface for the Fred node.
func Setup(addr string, h exthandler.Handler) error {
	gin.SetMode("release")
	r := gin.New()

	r.Use(logger.SetLogger(logger.Config{
		Logger: &log.Logger,
		UTC:    true,
	}))

	r.GET(apiversion+"/replica", getReplica(h))
	r.POST(apiversion+"/replica", postReplica(h))
	r.DELETE(apiversion+"/replica/:nodeid", deleteReplica(h))

	r.POST(apiversion+"/keygroup/:kgname", postKeygroup(h))
	r.DELETE(apiversion+"/keygroup/:kgname", deleteKeygroup(h))

	r.GET(apiversion+"/keygroup/:kgname/replica", getKeygroupReplica(h))
	r.POST(apiversion+"/keygroup/:kgname/replica/:nodeid", postKeygroupReplica(h))
	r.DELETE(apiversion+"/keygroup/:kgname/replica/:nodeid", deleteKeygroupReplica(h))

	r.GET(apiversion+"/keygroup/:kgname/data/:id", getItem(h))
	r.PUT(apiversion+"/keygroup/:kgname/data/:id", putItem(h))
	r.DELETE(apiversion+"/keygroup/:kgname/data/:id", deleteItem(h))

	return r.Run(addr)
}