package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/app"
)

// Fred is the main app logic.
type Fred interface {
	CreateKeygroup(kg string) error
	DeleteKeygroup(kg string) error
	Create(kg string, data string) (uint64, error)
	Read(kg string, id uint64) (string, error)
	Update(kg string, id uint64, data string) error
	Delete(kg string, id uint64) error
}

var addr = flag.String("addr", ":9001", "http service address")

func postKeygroup(context *gin.Context) {

	kgname := context.Params.ByName("kgname")

	err := app.CreateKeygroup(kgname)

	if err != nil {
		context.Status(http.StatusConflict)
		return
	}

	context.Status(http.StatusOK)
	return
}

func deleteKeygroup(context *gin.Context) {
	kgname := context.Params.ByName("kgname")

	err := app.DeleteKeygroup(kgname)

	if err != nil {
		context.Status(http.StatusConflict)
		return
	}

	context.Status(http.StatusOK)
	return
}

func postItem(context *gin.Context) {
	kgname := context.Params.ByName("kgname")
	var json struct {
		Data string `json:"data" binding:"required"`
	}

	if err := context.ShouldBindJSON(&json); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	data := json.Data

	id, err := app.Create(kgname, data)

	if err != nil {
		context.Status(http.StatusConflict)
		return
	}

	context.String(http.StatusOK, string(id))
	return

}

func getItem(context *gin.Context) {
	kgname := context.Params.ByName("kgname")

	id, err := strconv.Atoi(context.Params.ByName("id"))

	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	data, err := app.Read(kgname, uint64(id))

	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.String(http.StatusOK, data)
	return
}

func putItem(context *gin.Context) {
	kgname := context.Params.ByName("kgname")

	id, err := strconv.Atoi(context.Params.ByName("id"))

	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	var json struct {
		Data string `json:"data" binding:"required"`
	}

	if err := context.ShouldBindJSON(&json); err != nil {
		log.Print(err)
		context.Status(http.StatusBadRequest)
		return
	}

	data := json.Data
	err = app.Update(kgname, uint64(id), data)

	if err != nil {
		context.Status(http.StatusConflict)
		return
	}

	context.Status(http.StatusOK)
	return
}

func deleteItem(context *gin.Context) {
	kgname := context.Params.ByName("kgname")

	id, err := strconv.Atoi(context.Params.ByName("id"))

	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	err = app.Delete(kgname, uint64(id))

	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.Status(http.StatusOK)
	return
}

func setupRouter(app Fred) (r *gin.Engine) {
	r = gin.Default()

	r.POST("/keygroup/:kgname", postKeygroup)
	r.DELETE("/keygroup/:kgname", deleteKeygroup)

	r.POST("/keygroup/:kgname/items", postItem)
	r.GET("/keygroup/:kgname/items/:id", getItem)
	r.PUT("/keygroup/:kgname/items/:id", putItem)
	r.DELETE("/keygroup/:kgname/items/:id", deleteItem)

	return
}

func main() {
	var app Fred = app.New()

	r := setupRouter(app)

	log.Fatal(r.Run(*addr))
}
