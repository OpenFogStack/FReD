package webserver

import (
	"fmt"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// Setup sets up a web server client interface for the Fred node.
func Setup(host string, port int, h fred.ExtHandler, apiversion string, useTLS bool, loglevel string) error {
	gin.SetMode(loglevel)
	r := gin.New()

	r.Use(logger.SetLogger(logger.Config{
		Logger: &log.Logger,
		UTC:    true,
	}))

	r.GET(apiversion+"/replica", getAllReplica(h))
	r.POST(apiversion+"/replica", postReplica(h))
	r.GET(apiversion+"/replica/:nodeid", getReplica(h))
	r.DELETE(apiversion+"/replica/:nodeid", deleteReplica(h))

	r.POST(apiversion+"/keygroup/:kgname", postKeygroup(h))
	r.DELETE(apiversion+"/keygroup/:kgname", deleteKeygroup(h))

	r.GET(apiversion+"/keygroup/:kgname/replica", getKeygroupReplica(h))
	r.POST(apiversion+"/keygroup/:kgname/replica/:nodeid", postKeygroupReplica(h))
	r.DELETE(apiversion+"/keygroup/:kgname/replica/:nodeid", deleteKeygroupReplica(h))

	r.GET(apiversion+"/keygroup/:kgname/data/:id", getItem(h))
	r.PUT(apiversion+"/keygroup/:kgname/data/:id", putItem(h))
	r.DELETE(apiversion+"/keygroup/:kgname/data/:id", deleteItem(h))

	if useTLS {
		if port != 443 {
			log.Warn().Msgf("HTTPS server needs to run on port 443 but port %d was given. Port 443 will be used anyway. To request a certificate, port 80 also needs to be available.", port)
		}

		return autotls.Run(r, host)
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Debug().Msgf("Starting web server on %s", addr)
	return r.Run(addr)
}
