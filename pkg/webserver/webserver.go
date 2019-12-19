package webserver

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/app"
)

// Fred is the main app logic.
type Fred interface {
	ExtHandleCreateKeygroup(kg string) error
	ExtHandleDeleteKeygroup(kg string) error
	ExtHandleRead(kg string, id string) (string, error)
	ExtHandleUpdate(kg string, id string, data string) error
	ExtHandleDelete(kg string, id string) error
}

const apiversion string = "/v0"
var a Fred

func postKeygroup(context *gin.Context) {

	kgname := context.Params.ByName("kgname")

	err := a.ExtHandleCreateKeygroup(kgname)

	if err != nil {
		context.Status(http.StatusConflict)
		return
	}

	context.Status(http.StatusOK)
	return
}

func deleteKeygroup(context *gin.Context) {
	kgname := context.Params.ByName("kgname")

	err := a.ExtHandleDeleteKeygroup(kgname)

	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.Status(http.StatusOK)
	return
}

func getItem(context *gin.Context) {
	kgname := context.Params.ByName("kgname")

	id := context.Params.ByName("id")

	data, err := a.ExtHandleRead(kgname, id)

	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.String(http.StatusOK, data)
	return
}

func putItem(context *gin.Context) {
	kgname := context.Params.ByName("kgname")

	id := context.Params.ByName("id")

	var json struct {
		Data string `json:"data" binding:"required"`
	}

	if err := context.ShouldBindJSON(&json); err != nil {
		log.Print(err)
		context.Status(http.StatusBadRequest)
		return
	}

	data := json.Data
	err := a.ExtHandleUpdate(kgname, id, data)

	if err != nil {
		context.Status(http.StatusConflict)
		return
	}

	context.Status(http.StatusOK)
	return
}

func deleteItem(context *gin.Context) {
	kgname := context.Params.ByName("kgname")

	id := context.Params.ByName("id")

	err := a.ExtHandleDelete(kgname, id)

	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.Status(http.StatusOK)
	return
}

// SetupRouter sets up a webserver for Fred
func SetupRouter(addr string, fred *app.App) error {
	a = fred

	r := gin.Default()

	r.POST(apiversion + "/keygroup/:kgname", postKeygroup)
	r.DELETE(apiversion + "/keygroup/:kgname", deleteKeygroup)

	r.GET(apiversion + "/keygroup/:kgname/items/:id", getItem)
	r.PUT(apiversion + "/keygroup/:kgname/items/:id", putItem)
	r.DELETE(apiversion + "/keygroup/:kgname/items/:id", deleteItem)

	return r.Run(addr)
}